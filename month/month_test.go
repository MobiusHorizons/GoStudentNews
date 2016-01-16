package month

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	m := New("http://www.calvin.edu/archive/student-news/201503/", "March 2015")
	c := m.GetIssues()
	<-c
	for _, i := range m.Issues {
		fmt.Println(i.Date)
	}
}
