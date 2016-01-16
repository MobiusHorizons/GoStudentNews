package month

import (
	"github.com/MobiusHorizons/GoStudentNews/issue"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

type Month struct {
	url    string
	Date   time.Time
	Issues []issue.Issue
}

func New(url string, date string) *Month {
	d, e := time.Parse("January 2006", date)
	if e != nil {
		log.Fatal(e)
	}
	m := Month{url, d, []issue.Issue{}}
	return &m
}

func (m *Month) GetIssues() chan *Month {
	c := make(chan *Month)
	go func() {
		doc, err := goquery.NewDocument(m.url)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(".messages-list ul li").Each(func(i int, s *goquery.Selection) {
			issueUrl, _ := s.Find("a[href]").Attr("href")
			date := s.Find("em").Get(1).FirstChild.Data
			issueUrl = m.url + issueUrl
			newIssue := issue.New(issueUrl, date)
			m.Issues = append(m.Issues, *newIssue)
		})
		c <- m
	}()
	return c
}
