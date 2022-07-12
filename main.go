package main

import (
	"log"
	"net/http"
)

var Blockchain *Blockchain
type Block struct {

}
type Book struct {

}

type BookCheckout struct{

}
type Blockchain struct {
	blocks []*Block
}


func main() {
	r := mux.newRouter()
	r.HandleFunc("/", GetBlockchain ).Methods("GET")
	r.HandleFunc("/", WriteBlock ).Methods("POST")
	r.HandleFunc("/newbook", NewBook ).Methods("POST")

	log.Println("listening on port 2000")
	log.Fatal(http.ListenAndServe(":2000", r))

}

