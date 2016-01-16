package main

import (
	"fmt"
	"github.com/MobiusHorizons/GoStudentNews/archive"
	"github.com/MobiusHorizons/GoStudentNews/issue"
	"github.com/MobiusHorizons/GoStudentNews/month"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"sync"
)

var threadPool = make(chan int, 6)
var gtkList *gtk.ListBox

func progress(i issue.Issue, wg *sync.WaitGroup) {
	defer wg.Done()
	threadPool <- 1
	c := i.GetEntries()

	for e := range c {
		fmt.Println(e.Subject)
		listAppend(e.Subject)
	}
	<-threadPool
}

func waitIssue(m month.Month, wg *sync.WaitGroup) {
	defer wg.Done()
	threadPool <- 1
	c := m.GetIssues()
	<-c
	<-threadPool

	fmt.Printf("Got %d issues\n", len(m.Issues))
	/*
		wg.Add(len(m.Issues))

		for _, i := range m.Issues {
			fmt.Println("fetching Entries for ", i.Date)
			go progress(i, wg)
		}*/
	wg.Add(1)
	go progress(m.Issues[0], wg)
}

func getEntries() {
	a := archive.New("http://www.calvin.edu/archive/student-news/")
	c := a.GetMonths()
	var wg sync.WaitGroup
	<-c

	fmt.Printf("Got %d months\n", len(a.Months))
	/*
		wg.Add(len(a.Months))
		for _, m := range a.Months {
			fmt.Println("Fetching Issues from " + m.Date.Format("January 2006"))
			go waitIssue(m, &wg)
		}
	*/
	//*/
	m := a.Months[0]
	wg.Add(1)
	fmt.Println("Fetching Issues from " + m.Date.Format("January 2006"))
	go waitIssue(m, &wg)
	//*/

	wg.Wait()
}

func listAppend(labelText string) {
	l, err := gtk.LabelNew(labelText)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	gtkList.Add(l)
	gtkList.ShowAll()
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new label widget to show in the window.
	list, err := gtk.ListBoxNew()
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	gtkList = list
	// Add the label to the window.
	win.Add(gtkList)

	// Set the default window size.
	win.SetDefaultSize(640/2, 1136/2)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	go getEntries()
	gtk.Main()
}
