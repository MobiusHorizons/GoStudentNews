package main
import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"time"
	"calvin.edu/studentNews/month"
	"calvin.edu/studentNews/issue"
)


var archiveURL = "http://www.calvin.edu/archive/student-news/"

func main_m(){
	doc, err := goquery.NewDocument(archiveURL);
	if (err != nil){
		log.Fatal(err);
	}
	doc.Find("body > a").Each(func(i int, s *goquery.Selection){
		date := s.Text();
		href, _ := s.Attr("href");
		href = archiveURL + href;
		fmt.Println(date, href);
		cb := func(url string, date time.Time){
			fmt.Println(url, date);
		}
		go month.Issues(href, cb );
	});
}

func main(){
	issues := issue.Parse("http://www.calvin.edu/archive/student-news/201509/0017.html", time.Now());
	for _, i := range issues {
		i.Print();
	}
}
