// search in vietnamese is "tim kiem"
// it is a simple library so just call it "tim"
// this is an upgraded version of tim2: can find specific part of the documentation

package tim2

import (
	"github.com/gocql/gocql"

	"fmt"
	"sort"
	"strconv"
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
// + Sử dụng 2 bảng:
// CREATE KEYSPACE tim2 WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};
// CREATE TABLE terms (col ASCII, acc ASCII, term ASCII, day INT, doc ASCII, part ASCII, PRIMARY KEY((col, acc, term), day, doc, part)) WITH CLUSTERING ORDER BY (day DESC, doc DESC, part DESC);
// CREATE TABLE docs2 (col ASCII, acc ASCII, terms LIST<ASCII>, doc ASCII, part ASCII, day INT, PRIMARY KEY((col, acc, doc), part));

// INSERT INTO terms(col, acc, term, date, doc, part) VALUES('test', 'acctest', 'thanh', 4000, 'cs123', 'p1');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4000, 'cs1', 'p1');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4100, 'cs2', 'p2');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4200, 'cs3', 'p3');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4300, 'cs6', 'p4');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4400, 'cs4', 'p5');
// INSERT INTO terms(col, acc, term, day, doc, part) VALUES('test', 'acctest', 'thanh', 4500, 'cs5', 'p6');
// docs (col ASCII, acc ASCII, doc ASCII, term ASCII, PRIMARY KEY((col, acc, doc), term));

func SearchPart(col, accid, query, anchor string) ([]Hit, string, error) {
	return doSearch(col, accid, query, anchor, nil, 30, false)
}

// Search all docs that match the query
// docs, terms, anchor, error
func Search(col, accid, query, anchor string) ([]Hit, string, error) {
	return doSearch(col, accid, query, anchor, nil, 30, true)
}

func SearchPartOnly(col, accid, query, anchor string, only_parts []string, limit int) ([]Hit, string, error) {
	return doSearch(col, accid, query, anchor, only_parts, 30, false)
}

type Hit struct {
	doc  string
	part string
}

func doSearch(collection, accid, query, anchor string, only_parts []string, limit int, doc_distinct bool) ([]Hit, string, error) {
	waitforstartup()
	terms := Tokenize(query)
	if len(terms) == 0 {
		return nil, anchor, nil
	}

	// remove all child terms ('xin chao', 'xin', 'cong' => 'xin' is child term, 'cong' is not)
	// only using parent terms to search, since doc match parent term will match child term
	parentTerms := []string{}
	for i := 0; i < len(terms); i++ {
		isParent := true
		for j := i + 1; j < len(terms); j++ {
			if strings.Contains(terms[j], terms[i]) {
				isParent = false
				break
			}

		}
		if isParent {
			parentTerms = append(parentTerms, terms[i])
		}
	}

	terms = parentTerms

	// long query
	sort.Slice(terms, func(i, j int) bool { return len(terms[i]) > len(terms[j]) })
	// contain all matched doc

	// anchor
	anchorsplit := strings.Split(anchor, "_")
	anchorday := int(time.Now().UnixNano() / 86400 * 10) // a very large day
	matchanchor := true
	if len(anchorsplit) > 2 {
		i, err := strconv.Atoi(anchorsplit[0])
		if err == nil {
			// anchor is valid
			matchanchor = false // need match anchor
			anchorday = i
		}
	}

	hits := []Hit{}
	var day int
	iter := db.Query("SELECT day, doc, part FROM tim2.terms WHERE col=? AND acc=? AND term=? AND day<=? ORDER BY day DESC", collection, accid, terms[0], anchorday).Iter()
	var docid, part string
	for iter.Scan(&day, &docid, &part) {

		// skip until pass the anchor mark
		if day == anchorday && !matchanchor {
			if fmt.Sprintf("%d_%s_%s", day, docid, part) == anchor {
				matchanchor = true
			}
			continue
		}

		// filter parts if only_parts is pass in
		if len(only_parts) > 0 {
			found := false
			for _, p := range only_parts {
				if p == part {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		// de-duplicate by doc id, since we index document in multiple day
		alreadyhitdoc := false
		alreadyhitdocpart := false
		for _, hit := range hits {
			if docid == hit.doc {
				alreadyhitdoc = true
				if hit.part == part {
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

		docterms := []string{}
		var docday int
		db.Query("SELECT day, terms FROM tim2.docs2 WHERE col=? AND acc=? AND doc=? AND part=? LIMIT 1", collection, accid, docid, part).Scan(&docday, &terms)
		if docday != day {
			continue // invalid doc
		}
		// the doc must match all other terms (max 5)
		matchAll := true
		for i := 1; i < len(terms) && i < 5; i++ {
			term := terms[i]
			found := false
			for _, dt := range docterms {
				if term == dt {
					found = true
					break
				}
			}

			if !found {
				matchAll = false
				break
			}
		}

		if !matchAll {
			continue
		}

		hits = append(hits, Hit{doc: docid, part: part})
		anchor = fmt.Sprintf("%d_%s_%s", day, docid, part)
		if len(hits) >= limit {
			break
		}
	}
	if err := iter.Close(); err != nil {
		return nil, anchor, err
	}
	return hits, anchor, nil
}

// Index creates reverse index for document part
// remove all old terms, write new terms
func Index(col, accid, doc, part string, day int, text string) error {
	waitforstartup()
	terms := Tokenize(text)

	var oldday int
	var oldterms []string
	db.Query("SELECT day, terms FROM tim2.docs2 WHERE col=? AND acc=? AND doc=? AND part=? LIMIT 1", col, accid, doc, part).Scan(&oldday, &oldterms)

	// findings new terms (so we only have to remove outdated terms instead of all old terms)
	var outdates, news []string
	if oldday == day {
		if len(oldterms) > 0 {
			var oldM = map[string]bool{}
			for _, o := range oldterms {
				oldM[o] = true
			}

			var newM = map[string]bool{}
			for _, n := range terms {
				newM[n] = true
			}

			for _, n := range terms {
				if !oldM[n] {
					news = append(news, n)
				}
			}

			for _, o := range oldterms {
				if !newM[o] {
					outdates = append(outdates, o)
				}
			}

		} else {
			// quick path for empty oldterms
			news = terms
		}
	} else {
		// difference day, must delete all outdated term
		news = terms
		outdates = oldterms
	}

	// delete old (outdated) terms
	batch := db.NewBatch(gocql.LoggedBatch)
	for i, term := range outdates {
		batch.Query("DELETE FROM tim2.terms WHERE col=? AND acc=? AND term=? AND doc=? AND part=? AND day=?", col, accid, term, doc, part, day)
		if batch.Size()%50 == 0 || (i == len(outdates)-1 && batch.Size() > 0) {
			if err := db.ExecuteBatch(batch); err != nil {
				return err
			}
			batch = db.NewBatch(gocql.LoggedBatch)
		}
	}

	if err := db.Query("INSERT INTO tim2.docs2(col,acc,doc,terms,part,day) VALUES(?,?,?,?,?,?)", col, accid, doc, terms, part, day).Exec(); err != nil {
		return err
	}

	// write new terms to docs
	for i, term := range news {
		batch.Query("INSERT INTO tim2.terms(col,acc,term,day,doc,part) VALUES(?,?,?,?,?,?)", col, accid, term, day, doc, part)
		if batch.Size()%50 == 0 || (i == len(news)-1 && batch.Size() > 0) {
			if err := db.ExecuteBatch(batch); err != nil {
				return err
			}
			batch = db.NewBatch(gocql.LoggedBatch)
		}
	}
	return nil
}

var startupLock sync.Mutex
var db *gocql.Session

func waitforstartup() {
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
