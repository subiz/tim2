// search in vietnamese is "tim kiem"
// it is a simple library so just call it "tim"
// this is an upgraded version of tim2: can find specific part of the documentation

package tim2

import (
	"github.com/gocql/gocql"

	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

//
// Can be used to:
// + search conversation (conversation service)
// + search lead or user (user service)
// + search account (account service)
// + search for a product (product service)

// Thuật ngữ:
// + Document: thứ mà người dùng cần tìm (conversation, user, lead, account, ...)
// + Query: đoạn text mà user nhập để tìm các document
// + Tokenize: Hàm dùng để biến một đoạn text thành các term
//
//
// Cách sử dụng:
// indexer.AppendText("convo", "acc1", "convo1", "xin chào")
// indexer.AppendText("convo", "acc1", "convo1", "cộng hòa xã hội chủ nghĩa việt nam")
//
// indexer.AppendText("convo", "acc1", "convo2", "xin tạm biệt")
// indexer.AppendText("convo", "acc1", "convo2", "độc lập tự do hạnh phúc")
//
// indexer.Search("convo", "acc1", "ag1", "xin")
//   => 2 hits: [convo1, convo2]
// indexer.Search("convo", "acc1", "ag2", "xin")
//   => 1 hit: [convo1]
//
//
// Thiết kế:
// + Trong database chỉ lữu trữ docId, nội dung doc không quan tâm và không lưu
// + Sử dụng 1 bảng:
//   CREATE KEYSPACE tim2 WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};
//
//   CREATE TABLE tim2.term_doc(date BIGINT, col ASCII, acc ASCII, term ASCII, doc ASCII, PRIMARY KEY ((col, acc, term), doc, part)) WITH CLUSTERING ORDER BY (doc DESC, part DESC);
//     dùng để tìm kiếm docs theo term
//

// task (date BIGINT, acc ASCII, term TEXT, doc ASCII, PRIMARY KEY((date, acc, term), doc))
// order(date BIGINT, acc ASCII, term TEXT, doc ASCII, PRIMARY KEY((date, acc, term), doc)) // product // user email
// convo(date BIGINT, acc ASCII, term TEXT, doc ASCII, PRIMARY KEY((date, acc, term), doc))

// CREATE TABLE terms (col ASCII, acc ASCII, term TEXT, day BIGINT, doc ASCII, part ASCII, PRIMARY KEY((col, acc, term), day, doc, part)) WITH CLUSTERING ORDER BY (day DESC, doc DESC, part DESC);
// CREATE TABLE docs (col ASCII, acc ASCII, term TEXT, doc ASCII, part ASCII, PRIMARY KEY ((col, acc, term, doc), term, part))

// INSERT INTO terms(col, acc, term, date, doc, part) VALUES('test', 'acctest', 'thanh', 4000, 'cs123', 'p1');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4000, 'cs1', 'p1');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4100, 'cs2', 'p2');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4200, 'cs3', 'p3');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4300, 'cs6', 'p4');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4400, 'cs4', 'p5');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4500, 'cs5', 'p6');
// docs (col ASCII, acc ASCII, doc ASCII, term ASCII, PRIMARY KEY((col, acc, doc), term));

// # list all task that have user contains keyword 'thanh'
// SELECT * FROM terms where acc='acctest' AND col='test' AND term='thanh' ORDER BY day DESC;
//
// # when update user, just

// # list all task that contains keyword 'thanh'
// SELECT from terms where date='05-2022' AND col='task' AND acc='acsble' AND term='thanh' AND typ='task'

// text => term

func SearchPart(collection, accid, query string, limit int, validate func(doc, part string) bool) ([]string, []string, error) {
	return doSearch(collection, accid, query, limit, validate, false)
}

// Search all docs that match the query
// docs, terms, anchor, error
func Search(collection, accid, query string, limit int, validate func(doc, part string) bool) ([]string, []string, error) {
	return doSearch(collection, accid, query, limit, validate, true)
}

func doSearch(collection, accid, query string, limit int, validate func(doc, part string) bool, doc_distinct bool) ([]string, []string, error) {
	waitforstartup(collection, accid)
	interms := Tokenize(query)
	if len(interms) == 0 {
		return []string{}, []string{}, nil
	}

	// long query
	var terms []string
	if len(interms) > 5 {
		biwords := make([]string, 0)
		for _, term := range interms {
			if strings.Contains(term, " ") {
				biwords = append(biwords, term)
			}
		}
		terms = make([]string, 0)
		for i := 0; i < 2 && i < len(biwords); i++ {
			terms = append(terms, biwords[i])
		}
		if len(terms) < 2 {
			for i := 0; i < 5-len(terms); i++ {
				terms = append(terms, interms[i])
			}
		}
	} else {
		terms = interms
	}
	// order by length desc
	sort.Slice(terms, func(i, j int) bool { return len(terms[i]) > len(terms[j]) })

	// contain all matched doc
	hitDocs := []string{}
	hitParts := []string{}

	iter := db.Query("SELECT doc, part FROM tim2.terms WHERE col=? AND acc=? AND term=? ORDER BY day DESC", collection, accid, terms[0]).Iter()
	var docid, part string
	for iter.Scan(&docid, &part) {
		// de-duplicate by doc id, since we index document in multiple day
		alreadyhitdoc := false
		alreadyhitdocpart := false
		for i, hit := range hitDocs {
			if docid == hit {
				alreadyhitdoc = true
				if hitParts[i] == part {
					alreadyhitdocpart = true
				}
			}
		}

		// both doc id and part id already found
		if alreadyhitdocpart {
			continue
		}

		// skip if user want doc distinct and doc already found
		if doc_distinct && alreadyhitdoc {
			continue
		}

		if !validate(docid, part) {
			continue
		}

		// the doc must match all other terms
		matchAll := true
		for i := 1; i < len(terms); i++ {
			term := terms[i]
			dump := ""
			db.Query("SELECT doc FROM tim2.docs WHERE col=? AND acc=? AND term=? AND doc=? AND part=? LIMIT 1", collection, accid, term, docid, part).Scan(&dump)
			if dump == "" {
				matchAll = false
				break
			}
		}

		if !matchAll {
			continue
		}

		hitDocs, hitParts = append(hitDocs, docid), append(hitParts, part)
		if len(hitDocs) >= limit {
			break
		}
	}
	if err := iter.Close(); err != nil {
		return nil, nil, err
	}
	return hitDocs, hitParts, nil
}

func doIndex(collection, accid, docId, part string, day int, text string) error {
	terms := Tokenize(text)
	batch := db.NewBatch(gocql.LoggedBatch)
	for i, term := range terms {
		batch.Query("INSERT INTO tim2.terms(col,acc,term,day,doc,part) VALUES(?,?,?,?,?,?)", collection, accid, term, day, docId, part)
		batch.Query("INSERT INTO tim2.docs(col,acc,doc,term,part) VALUES(?,?,?,?,?)", collection, accid, docId, term, part)
		if i%50 == 0 {
			if err := db.ExecuteBatch(batch); err != nil {
				return err
			}
			batch = db.NewBatch(gocql.LoggedBatch)
		}
	}

	if batch.Size() > 0 {
		if err := db.ExecuteBatch(batch); err != nil {
			return err
		}
	}
	return nil

}

func IndexText(collection, accid, docId, part string, day int, text string) error {
	waitforstartup(collection, accid)

	// clear term first
	err := db.Query("DELETE FROM tim2.docs WHERE col=? AND acc=? AND doc=? VALUES(?,?,?)", collection, accid, docId).Exec()
	if err != nil {
		return err
	}

	return doIndex(collection, accid, docId, part, day, text)
}

func AppendText(collection, accid, docId, part string, day int, text string) error {
	waitforstartup(collection, accid)
	return doIndex(collection, accid, docId, part, day, text)
}

var startupLock sync.Mutex
var db *gocql.Session

func waitforstartup(collection, accid string) {
	startupLock.Lock()
	defer startupLock.Unlock()

	// connect db
	if db != nil {
		return
	}

	cluster := gocql.NewCluster("db-0")
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Keyspace = "tim2"
	var err error
	for {
		if db, err = cluster.CreateSession(); err == nil {
			break
		}
		fmt.Println("cassandra", err, ". Retring after 5sec...")
		time.Sleep(5 * time.Second)
	}
}
