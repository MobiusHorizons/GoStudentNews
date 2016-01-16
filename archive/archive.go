package archive

import (
	//"fmt"
	"github.com/MobiusHorizons/GoStudentNews/month"
	"github.com/PuerkitoBio/goquery"
	"log"
	//  "time"
)

type Archive struct {
	url    string
	Months []month.Month
}

func New(url string) *Archive {
	a := Archive{url, []month.Month{}}
	return &a
}

func (a *Archive) GetMonths() chan *Archive {
	c := make(chan *Archive)
	go func() {
		doc, err := goquery.NewDocument(a.url)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find("body > a").Each(func(i int, s *goquery.Selection) {
			date := s.Text()
			href, _ := s.Attr("href")
			href = a.url + href
			//fmt.Println(date, href)
			m := month.New(href, date)
			a.Months = append(a.Months, *m)
		})
		c <- a
	}()
	return c
}
