package archive

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	studentNews := New("http://www.calvin.edu/archive/student-news/")
	c := studentNews.getMonths()
	<-c
	for _, m := range studentNews.Months {
		fmt.Println(m.Date)
	}
}
