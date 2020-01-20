package apiserver

import (
	"fmt"
	"net/http"
)

func Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	http.ListenAndServe(":5000", nil)
}