package issue
import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"log"
	"time"
	"regexp"
	"strings"
	"bytes"
	"golang.org/x/net/html"
)


var headSeperator = "----------------------------------------------------------------------";
var issueSeperator = "------------------------------";

var dateRegex = regexp.MustCompile("Date:\\s+(.+)\\s+From:");
var fromRegex = regexp.MustCompile("From:\\s+(.+)\\s+Subject:");
var subjectRegex = regexp.MustCompile("Subject:\\s(.+)");

type Issue struct{
	Date time.Time
	From string
	Subject string
	Body string
};

func (i *Issue) Print (){
	fmt.Println("Date: " , i.Date);
	fmt.Println("From: " , i.From);
	fmt.Println("Subject: ", i.Subject);
	fmt.Println(i.Body);
}

func outerHtml (n *html.Node) string{
	var b bytes.Buffer;
	html.Render(&b,n);
	s:= b.String();
	return s;
}

func Parse( url string, date time.Time) {
	doc, err := goquery.NewDocument(url);
	if (err != nil){
		log.Fatal(err);
	}
	head := true;
	currentIssue := new(Issue);
	var issues []Issue;

	doc.Find(".mail p").Each(func(i int, s *goquery.Selection){
		line := s.Text();
		if (head){
			if (strings.TrimSpace(line) == headSeperator){
				head = false;
			}
			return;
		}
		haveDate, haveSubject, haveFrom := false,false,false;
		if (strings.TrimSpace(line) == issueSeperator){
			currentIssue.Print();
			issues = append(issues, *currentIssue);
			currentIssue = new(Issue);
			fmt.Println("starting New issue");
			return;
		}

		if (dateRegex.MatchString(line)){
			dateMatch := dateRegex.FindStringSubmatch(line);
			if (dateMatch != nil && len(dateMatch) > 1){
				haveDate = true;
				currentIssue.Date, _  = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", dateMatch[1]);
			}
		}
		if (fromRegex.MatchString(line)){
			fromMatch := fromRegex.FindStringSubmatch(line);
			haveFrom = true;
			if (fromMatch != nil && len(fromMatch) > 1){
				currentIssue.From = fromMatch[1];
			}
		}
		if (subjectRegex.MatchString(line)){
			haveSubject = true;
			subjectMatch := subjectRegex.FindStringSubmatch(line);
			if (subjectMatch != nil && len(subjectMatch) > 1){
				currentIssue.Subject = subjectMatch[1];
			}
		}
		if ( !(haveDate && haveFrom && haveSubject)){
			pTag := goquery.NewDocumentFromNode(s.Nodes[0]);
			fmt.Println(outerHtml(pTag.Nodes[0]));
			pTag.Find("br").Remove();
			fmt.Println(outerHtml(pTag.Nodes[0]));
			
			line = outerHtml(pTag.Nodes[0]);
			//line = strings.NewReplacer("\n", " ", "\r", " ", "  ", " ").Replace(line);
			if (err == nil) {currentIssue.Body += line;}
		}
	});
	//return issues;
}
