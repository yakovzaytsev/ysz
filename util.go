package ysz

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HandleHomePage(mux *http.ServeMux) {
	var CLIENT_HOME string
	var ok bool // XXX
	if CLIENT_HOME, ok = os.LookupEnv("CLIENT_HOME"); !ok {
		log.Fatal("CLIENT_HOME environment variable not set. Should point to some.html or client/dist")
	}
	log.Printf("path to home %v\n", CLIENT_HOME)

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
			log.Fatal("Could not ReadFile %v", err)
		}
		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(page)
		}))
	}
}
