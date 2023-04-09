package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cocktails-ai/backend/gpt" // Update with the correct import path for your 'gpt' package
)

type RequestPayload struct {
	Drinks []string `json:"drinks"`
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// Convert the drinks list to a prompt for the GPT message.
	drinksList := strings.Join(payload.Drinks, ", ")
	inputMessage := fmt.Sprintf("I have the following ingredients: %s. What cocktails can I make with these ingredients?", drinksList)

	response, err := gpt.Message(inputMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling Message function: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/cocktails", messageHandler)

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
