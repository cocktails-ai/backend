package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/b0noi/go-utils/v2/gcp"
	"github.com/cocktails-ai/backend/barcode"
	"github.com/cocktails-ai/backend/gpt"
	"github.com/cocktails-ai/backend/model"
)

type RequestPayload struct {
	Drinks []string `json:"drinks"`
}

type CocktailResponse struct {
	Cocktails       string           `json:"cocktails"`
	CocktailsParsed []model.Cocktail `json:"cocktails_parsed"`
}

type BarcodeRequestPayload struct {
	Barcode string `json:"barcode"`
}

type BarcodeResponse struct {
	ProductName string `json:"product_name"`
}

func messageHandlerBarcode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload BarcodeRequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	productName, err := barcode.FindBarcode(payload.Barcode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling FindBarcode function: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responseJson := BarcodeResponse{
		ProductName: productName,
	}

	json.NewEncoder(w).Encode(responseJson)
}

func messageHandlerGpt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		fmt.Printf("Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	inputMessageTemplate, err := gcp.ReadFile("cocktails-ai", "gpt4-prompt.txt")
	if err != nil {
		fmt.Printf("Error calling ReadFile function: %v", err)
		http.Error(w, fmt.Sprintf("Error calling ReadFile function: %v", err), http.StatusInternalServerError)
		return
	}
	response, err := requestGpt(payload.Drinks, inputMessageTemplate)
	if err != nil {
		fmt.Printf("Error calling requestGpt function: %v", err)
		http.Error(w, fmt.Sprintf("Error calling requestGpt function: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	rawResponse := response.Choices[0].Message.Content
	parsedCocktails, err := model.ParseCocktailsMarkdown(rawResponse)
	if err != nil {
		parsedCocktails = nil
	}
	responseJson := CocktailResponse{
		Cocktails:       rawResponse,
		CocktailsParsed: parsedCocktails,
	}

	json.NewEncoder(w).Encode(responseJson)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If it's an OPTIONS request, respond with just the headers, otherwise call the next handler.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func requestGpt(drinks []string, promptTemplate string) (gpt.GptChatCompletionMessage, error) {
	drinksList := strings.Join(drinks, ", ")
	inputMessage := fmt.Sprintf(promptTemplate, drinksList)

	response, err := gpt.Message(inputMessage)
	if err != nil {
		return gpt.GptChatCompletionMessage{}, err
	}
	return response, nil
}

func main() {
	http.Handle("/cocktails", corsMiddleware(http.HandlerFunc(messageHandlerGpt)))
	http.Handle("/barcodes", corsMiddleware(http.HandlerFunc(messageHandlerBarcode)))

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// func main() {
// 	productName, err := barcode.FindBarcode("0737628064502")
// 	if err != nil {
// 		fmt.Sprintf("Error calling ReadFile function: %v", err)
// 		return
// 	}
// 	fmt.Println(productName)
// }

// func main() {
// 	drinksList := strings.Join([]string{"vodka", "lime juice", "gyn"}, ", ")
// 	// inputMessageTemplate, err := gcp.ReadFile("cocktails-ai", "gpt4-prompt.txt")

// 	inputMessageTemplateBuffer, err := ioutil.ReadFile("./prompt.txt")
// 	if err != nil {
// 		log.Fatalf("Error reading file: %v", err)
// 	}

// 	inputMessageTemplate := string(inputMessageTemplateBuffer)
// 	if err != nil {
// 		fmt.Sprintf("Error calling ReadFile function: %v", err)
// 		return
// 	}
// 	inputMessage := fmt.Sprintf(inputMessageTemplate, drinksList)

// 	response, err := gpt.Message(inputMessage)

// 	if err != nil {
// 		fmt.Sprintf("Error calling Message function: %v", err)
// 		return
// 	}

// 	rawResponse := response.Choices[0].Message.Content
// 	parsedCocktails, err := model.ParseCocktailsMarkdown(rawResponse)
// 	if err != nil {
// 		parsedCocktails = nil
// 	}
// 	responseJson := CocktailResponse{
// 		Cocktails:       rawResponse,
// 		CocktailsParsed: parsedCocktails,
// 	}

// 	jsonOutput, err := json.MarshalIndent(responseJson, "", "  ")
// 	fmt.Println(string(jsonOutput))
// }
