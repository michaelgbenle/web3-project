package main

import "fmt"

func main() {
	r := mux.newRouter()
	r.HandleFunc("/", ).Methods("GET")
	fmt.Println("connectring to server")

}