package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var Blockchain *Blockchain
type Block struct {
Position 	int
Data		BookCheckout
TimeStamp	string
Hash		string
PrevHash	string
}
type Book struct {
ID 			string `json: "id"`
Title 		string	`json: "title"`
Author	 	string	`json: "author"`
PublishDate string	`json: "publish_date"`
Isbn 		string	`json: "isbn"`

}

type BookCheckout struct{
	BookId 			string	`json: "book_id"`
	User			string	`json: "user"`
	CheckoutDate 	string	`json: "checkout_date"`
	IsGenesis 		bool	`json: "is_genesis"`

}
type Blockchain struct {
	blocks []*Block
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", GetBlockchain ).Methods("GET")
	r.HandleFunc("/", WriteBlock ).Methods("POST")
	r.HandleFunc("/newbook", NewBook ).Methods("POST")

	log.Println("listening on port 2000")
	log.Fatal(http.ListenAndServe(":2000", r))

}

func NewBook (w http.ResponseWriter, r *http.Request){
	var book Book
	if err:= json.NewDecoder(r.Body).Decode(&book); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not bind json"))
		return
	}
	h:=  md5.New()
	io.WriteString(h, book.Isbn + book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	response, err :=json.MarshalIndent(book,"", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not save book data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
func GetBlockchain (w http.ResponseWriter, r *http.Request){

}
func WriteBlock (w http.ResponseWriter, r *http.Request){
var checkoutItem BookCheckout
if err:= json.NewDecoder(r.Body).Decode(&checkoutItem)

}