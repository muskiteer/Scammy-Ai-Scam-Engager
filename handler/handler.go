package handler

import (
	"net/http"
	"os"
)

type Request struct {
	SessionID string `json:"sessionId"`
	Message   MessageResponse `json:"message"`
	ConvoHistory  []MessageResponse `json:"conversationHistory"`
	Metadata    Metadata `json:"metadata"`
}

type MessageResponse struct {
	Sender    string  `json:"sender"`
	Text 	string  `json:"text"`
	Timestamp string `json:"timestamp"`
}

type Metadata struct {
	Channel string `json:"channel"`
	Language string `json:"language"`
	Locale string `json:"locale"`
}

type Response struct {
	Status string `json:"status"`
	Reply 	string `json:"reply"`
}

type ExtractedIntel struct {
	UPIs   []string `json:"upiIds"`
	Phones []string `json:"phoneNumbers"`
	Links  []string `json:"phishingLinks"`
	Banks  []string `json:"bankAccounts"`
	Suswords []string `json:"suspiciousKeywords"`
}

type FinalResponse struct {
	SessionID string `json:"sessionId"`
	ScamDetect     bool `json:"scamDetected"`
	TotalMessagesEx  int `json:"totalMessagesExchanged"`
	ExtractIntel   ExtractedIntel   `json:"extractedIntelligence"` 	
	AgentNote  string `json:"agentNotes"`	
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func StartConvo(w http.ResponseWriter, r *http.Request) {
	API_KEY := os.Getenv("API_KEY")
	if API_KEY == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}
	
}
	