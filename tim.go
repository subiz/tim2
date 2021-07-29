// search in vietnamese is "tim kiem"
// it is a simple library so just call it "tim"
// this is an upgraded version of tim2: can find specific part of the documentation

package tim

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
//   CREATE TABLE tim2.term_doc(col ASCII, acc ASCII, term ASCII, doc ASCII, PRIMARY KEY ((col, acc, term), doc, part)) WITH CLUSTERING ORDER BY (doc DESC, part DESC);
//     dùng để tìm kiếm docs theo term
//

// text => term

// isMatchOtherTerms checks to see if the document part match all the terms (except for the first one)
func isMatchOtherTerms(collection, accid, docid, part string, terms []string) bool {
	for i, term := range terms {
		if i == 0 {
			continue
		}
		dump := ""
		db.Query("SELECT doc FROM tim2.term_doc WHERE col=? AND acc=? AND term=? AND doc=? AND part=? LIMIT 1", collection, accid, term, docid, part).Scan(&dump)
		if dump == "" {
			return false
		}
	}
	return true
}

// Search all docs that match the query
// docs, terms, anchor, error
func Search(collection, accid, query, anchor string, limit int, ownerCheck func(doc string) bool) ([]string, []string, string, error) {
	waitforstartup(collection, accid)
	interms := Tokenize(query)
	if len(interms) == 0 {
		return []string{}, nil, "zzzzz", nil
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
	startdoc := "zzzz"
	startpart := "zzzz"
	anchorSplit := strings.Split(anchor, ";;;;")
	if len(anchorSplit) == 2 {
		// search inside the doc first
		startdoc, startpart = anchorSplit[0], anchorSplit[1]
		iter := db.Query("SELECT doc, part FROM tim2.term_doc WHERE col=? AND acc=? AND term=? AND doc=? AND part<?", collection, accid, terms[0], startdoc, startpart).Iter()
		var docid, part string
		for iter.Scan(&docid, &part) {
			// the doc must match owners condition
			if !ownerCheck(docid) {
				continue
			}

			if isMatchOtherTerms(collection, accid, docid, part, terms) {
				continue
			}
			hitDocs, hitParts = append(hitDocs, docid), append(hitParts, part)
			if len(hitDocs) >= limit {
				break
			}
		}
		if err := iter.Close(); err != nil {
			return nil, nil, "", err
		}
	}

	iter := db.Query("SELECT doc, part FROM tim2.term_doc WHERE col=? AND acc=? AND term=? AND doc<?", collection, accid, terms[0], startdoc).Iter()
	var docid, part string
	for iter.Scan(&docid, &part) {
		// the doc must match owners condition
		if !ownerCheck(docid) {
			continue
		}

		// the doc must match all other terms
		if isMatchOtherTerms(collection, accid, docid, part, terms) {
			continue
		}
		hitDocs, hitParts = append(hitDocs, docid), append(hitParts, part)
		if len(hitDocs) >= limit {
			break
		}

		if err := iter.Close(); err != nil {
			return nil, nil, "", err
		}
	}

	if len(hitDocs) == 0 {
		return []string{}, []string{}, anchor, nil
	}

	anchor = hitDocs[len(hitDocs)-1] + ";;;;" + hitParts[len(hitParts)-1]
	return hitDocs, hitParts, anchor, nil
}

func AppendText(collection, accid, docId, part, text string) error {
	waitforstartup(collection, accid)
	terms := Tokenize(text)
	for _, term := range terms {
		if err := db.Query("INSERT INTO tim2.term_doc(col,acc,term,doc,part) VALUES(?,?,?,?,?)", collection, accid, term, docId, part).Exec(); err != nil {
			return err
		}
	}
	return nil
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
