package issue

import (
	"bytes"
	"github.com/MobiusHorizons/GoStudentNews/entry"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"log"
	"regexp"
	"strings"
	"time"
)

var headSeperator = "----------------------------------------------------------------------"
var entrySeperator = "------------------------------"

var dateRegex = regexp.MustCompile("Date:\\s+(.+)\\s+From:")
var fromRegex = regexp.MustCompile("From:\\s+(.+)\\s+Subject:")
var subjectRegex = regexp.MustCompile("Subject:\\s(.+)")

type Issue struct {
	Url     string
	Date    time.Time
	Entries []entry.Entry
}

func outerHtml(n *html.Node) string {
	var b bytes.Buffer
	html.Render(&b, n)
	s := b.String()
	return s
}

func New(url, date string) *Issue {
	d, e := time.Parse("(Mon Jan 02 2006 - 15:04:05 MST)", date)
	if e != nil {
		log.Fatal(e)
	}
	issue := Issue{url, d, []entry.Entry{}}
	return &issue
}

func (issue *Issue) GetEntries() chan *entry.Entry {
	c := make(chan *entry.Entry)
	go func() {
		doc, err := goquery.NewDocument(issue.Url)
		if err != nil {
			log.Fatal(err)
		}
		head := true
		currentEntry := new(entry.Entry)

		doc.Find(".mail p").Each(func(i int, s *goquery.Selection) {
			line := s.Text()
			if head {
				if strings.TrimSpace(line) == headSeperator {
					head = false
				}
				return
			}
			haveDate, haveSubject, haveFrom := false, false, false
			if strings.TrimSpace(line) == entrySeperator {
				//currentEntry.Print()
				issue.Entries = append(issue.Entries, *currentEntry)
				c <- currentEntry
				currentEntry = new(entry.Entry)
				return
			}

			if dateRegex.MatchString(line) {
				dateMatch := dateRegex.FindStringSubmatch(line)
				if dateMatch != nil && len(dateMatch) > 1 {
					haveDate = true
					currentEntry.Date, _ = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", dateMatch[1])
				}
			}
			if fromRegex.MatchString(line) {
				fromMatch := fromRegex.FindStringSubmatch(line)
				haveFrom = true
				if fromMatch != nil && len(fromMatch) > 1 {
					currentEntry.From = fromMatch[1]
				}
			}
			if subjectRegex.MatchString(line) {
				haveSubject = true
				subjectMatch := subjectRegex.FindStringSubmatch(line)
				if subjectMatch != nil && len(subjectMatch) > 1 {
					currentEntry.Subject = subjectMatch[1]
				}
			}
			if !(haveDate && haveFrom && haveSubject) {
				pTag := goquery.NewDocumentFromNode(s.Nodes[0])
				pTag.Find("br").Remove()

				line = outerHtml(pTag.Nodes[0])
				//line = strings.NewReplacer("\n", " ", "\r", " ", "  ", " ").Replace(line);
				if err == nil {
					currentEntry.Body += line
				}
			}
		})
		close(c)
	}()
	return c
}
