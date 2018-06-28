package ysz

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HandleHomePage() {
	var HOME_PAGE string
	var ok bool // XXX
	if HOME_PAGE, ok = os.LookupEnv("HOME_PAGE"); !ok {
		log.Fatal("HOME_PAGE environment variable not set. Should point to some.html")
	}
	log.Printf("path to home page %v\n", HOME_PAGE)
	page, err := ioutil.ReadFile(HOME_PAGE)
	if err != nil {
		log.Fatal("Could not ReadFile %v", err)
	}
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))
}
