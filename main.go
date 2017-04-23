package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
)

var jsFiles map[string][]byte
var index []byte

func main() {
	generateResources()
	prepareFiles()

	r := mux.NewRouter()
	r.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		doHTTPWork(w, r)
	})
	r.HandleFunc("/http2", func(w http.ResponseWriter, r *http.Request) {
		doHTTP2Work(w, r)
	})
	r.HandleFunc("/res/{file}", func(w http.ResponseWriter, r *http.Request) {
		ServeFile(w, r)
	})

	var srv http.Server
	srv.Addr = ":8080"
	http2.ConfigureServer(&srv, nil)
	http.Handle("/", r)

	log.Fatal(srv.ListenAndServeTLS("certs/localhost.cert", "certs/localhost.key"))
}

func prepareFiles() {
	jsFiles = make(map[string][]byte)
	for i := 0; i <= jsFilesNumber; i++ {
		filename := fmt.Sprintf("f%03d.js", i)
		f, err := os.Open("res/" + filename)
		if err != nil {
			log.Fatal(err)
			return
		}
		buff, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
			return
		}
		jsFiles[filename] = buff
	}
	f, err := os.Open("res/index.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	buff, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
		return
	}
	index = buff
}

func doHTTPWork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(index)
}

func doHTTP2Work(w http.ResponseWriter, r *http.Request) {
	p, ok := w.(http.Pusher)
	if ok {
		for i := 0; i <= jsFilesNumber; i++ {
			filename := fmt.Sprintf("f%03d.js", i)
			err := p.Push("/res/"+filename, nil)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(index)
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	w.Write(jsFiles[mux.Vars(r)["file"]])
}
