package gpt

import (
	"fmt"
	"testing"
)

func TestMessage(t *testing.T) {
	inputMessage := "I have vodka, gin, tequila, rum, triple sec, and lemon juice. What cocktails can I make with these ingredients?"

	// Call the Message function.
	response, err := Message(inputMessage)

	// If there's an error, print it and mark the test as failed.
	if err != nil {
		t.Errorf("Error while calling Message function: %v", err)
		return
	}

	// Print the GPT response to the terminal.
	fmt.Printf("GPT Response: %s\n", response.Choices[0].Message.Content)
}
