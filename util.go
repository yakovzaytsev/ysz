package ysz

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HandleHomePage(mux *http.ServeMux, CLIENT_HOME string) {
	f, err := os.Stat(CLIENT_HOME)
	if err != nil {
		log.Fatal(err)
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		files := http.FileServer(http.Dir(CLIENT_HOME))
		mux.Handle("/", files)
	case mode.IsRegular():
		page, err := ioutil.ReadFile(CLIENT_HOME)
		if err != nil {
			log.Fatalf("Could not ReadFile %v", err)
		}
		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(page)
		}))
	}
}
