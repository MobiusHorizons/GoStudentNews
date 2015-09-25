package month
import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)



func Issues( url string, cb func(string, time.Time) ){
	doc, err := goquery.NewDocument(url);
	if (err != nil){
		log.Fatal(err);
	}
	doc.Find(".messages-list ul li").Each(func(i int, s *goquery.Selection){
		issueUrl, _ := s.Find("a[href]").Attr("href");
		date, error := time.Parse("(Mon Jan 2 2006 - 15:04:05 MST)", s.Find("em").Get(1).FirstChild.Data);
		if (error != nil){
			log.Fatal(error);
		}
		issueUrl = url + issueUrl;
		cb(issueUrl, date);
	});
}
