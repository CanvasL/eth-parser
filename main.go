package main

import (
	"encoding/json"
	"eth_parser/parser"
	"fmt"
	"net/http"
)

func main() {
	parser := parser.NewParser()

	http.HandleFunc("/current_block", func(w http.ResponseWriter, r *http.Request) {
		block := parser.GetCurrentBlock()
		json.NewEncoder(w).Encode(block)
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		success := parser.Subscribe(address)
		json.NewEncoder(w).Encode(success)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		transactions := parser.GetTransactions(address)
		json.NewEncoder(w).Encode(transactions)
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
