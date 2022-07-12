package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
func (b *Block) generateHash()  {
	bytes,_ := json.Marshal(b.Data)

	data := string(b.Position) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))

}

func (b *Block) ValidateHash(hash string) bool {
b.generateHash()
if b.Hash != hash {
	return false
}
return true
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
func (bc *Blockchain)AddBlock(data BookCheckout)  {
prevBlock := 	bc.blocks[len(bc.blocks)-1]
block := CreateBlock(prevBlock, data)
if ValidBlock(block, prevBlock){
	bc.blocks= append(bc.blocks, block)
}
}

func ValidBlock(block, prevBlock *Block) bool {
	if prevBlock.Hash !=  block.PrevHash {
		return false
	}
	if !block.ValidateHash(block.Hash) {
		return false
	}
	if prevBlock.Position + 1 != block.Position {
		return false
	}
return true
}



func CreateBlock(prevBlock *Block, checkoutItem BookCheckout) *Block{
block := &Block{}
block.Position = prevBlock.Position + 1
block.TimeStamp = time.Now().String()
block.PrevHash = prevBlock.Hash
block.generateHash()
return block
}

func GenesisBlock() *Block  {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockchain() *Blockchain  {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func main() {
	Blockchain = NewBlockchain()
	r := mux.NewRouter()
	r.HandleFunc("/", GetBlockchain ).Methods("GET")
	r.HandleFunc("/", WriteBlock ).Methods("POST")
	r.HandleFunc("/newbook", NewBook ).Methods("POST")

	go func() {
		for _,block := range Blockchain.Blocks {
				fmt.Printf("prev.hash :%x\n", block.prevHash)
				bytes,_:=json.MarshalIndent(block.data,""," ")
		}
	}

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
if err:= json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil{
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("could not write block"))
	return
}
Blockchain.AddBlock(checkoutItem)

}