package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/a{id:[0-9]+}", foo)
	http.ListenAndServe(":80", r)
}

func foo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	fmt.Println(r.Form)
	fmt.Println(vars)
	w.Write([]byte(fmt.Sprintf("id: %v", vars["id"])))
}
