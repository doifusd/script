package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Input your port")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Host)
		fmt.Fprintf(w, "request:%s", r.RequestURI)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
