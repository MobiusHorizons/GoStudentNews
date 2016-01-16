package entry

import (
	"fmt"
	"time"
)

type Entry struct {
	Date    time.Time
	From    string
	Subject string
	Body    string
}

func (e *Entry) Valid() bool {
	return !e.Date.IsZero() && e.From != "" && e.Subject != "" && e.Body != ""
}

func (e *Entry) Print() {
	fmt.Println("<div class=\"entry\">")
	fmt.Println("<div>Date: ", e.Date, "</div>")
	fmt.Println("<div>From: ", e.From, "</div>")
	fmt.Println("<div>Subject: ", e.Subject, "</div>")
	fmt.Println(e.Body)
	fmt.Println("</div>")
}
