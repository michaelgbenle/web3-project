package main

import (
	"bytes"
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

var BlockChain *Blockchain
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
block.Data = checkoutItem
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
	BlockChain = NewBlockchain()
	r := mux.NewRouter()
	r.HandleFunc("/", GetBlockchain ).Methods("GET")
	r.HandleFunc("/", WriteBlock ).Methods("POST")
	r.HandleFunc("/newbook", NewBook ).Methods("POST")

	go func() {
		for _,block := range BlockChain.blocks {
				fmt.Printf("prev.hash :%x\n", block.PrevHash)
				bytes,_:=json.MarshalIndent(block.Data,""," ")
				fmt.Printf("Data:%v\n", string(bytes))
				fmt.Printf("Hash:%x\n", block.Hash)
				fmt.Println()
		}
	}()

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
	jbytes,err:= json.MarshalIndent(BlockChain.blocks,""," ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jbytes))
}
func  WriteBlock (w http.ResponseWriter, r *http.Request){
var checkoutItem BookCheckout
if err:= json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil{
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("could not write block"))
	return
}
BlockChain.AddBlock(checkoutItem)
resp,err := json.MarshalIndent(checkoutItem,""," ")
if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("could not write block"))
	return
	}
w.WriteHeader(http.StatusOK)

}