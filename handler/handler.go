package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	
	"github.com/muskiteer/Ai-Scam/internal"
)

type Request struct {
	SessionID    string            `json:"sessionId"`
	Message      MessageResponse   `json:"message"`
	ConvoHistory []MessageResponse `json:"conversationHistory"`
	Metadata     Metadata          `json:"metadata"`
}

type MessageResponse struct {
	Sender    string `json:"sender"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

type Metadata struct {
	Channel  string `json:"channel"`
	Language string `json:"language"`
	Locale   string `json:"locale"`
}

type Response struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

type ExtractedIntel struct {
	BankAccounts       []string `json:"bankAccounts"`
	UPIIds             []string `json:"upiIds"`
	PhishingLinks      []string `json:"phishingLinks"`
	PhoneNumbers       []string `json:"phoneNumbers"`
	SuspiciousKeywords []string `json:"suspiciousKeywords"`
}

type FinalResponse struct {
	SessionID       string         `json:"sessionId"`
	ScamDetect      bool           `json:"scamDetected"`
	TotalMessagesEx int            `json:"totalMessagesExchanged"`
	ExtractIntel    ExtractedIntel `json:"extractedIntelligence"`
	AgentNote       string         `json:"agentNotes"`
}

const SCAM_THRESHOLD = 50

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func StartConvo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		requestKey := r.Header.Get("x-api-key")
		if requestKey != apiKey {
			http.Error(w, "Unauthorized: Invalid or missing API key", http.StatusUnauthorized)
			return
		}
	}
	body, _ := io.ReadAll(r.Body)
	log.Println("Received request body: ", string(body))

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.SessionID == "" || request.Message.Text == "" {
		http.Error(w, "sessionId and message.text are required", http.StatusBadRequest)
		return
	}

	// Get or create session
	store := internal.GetStore()
	session := store.Get(request.SessionID)

	// Add incoming message to history
	session.AddMessage(request.Message.Text)
	session.Context.TurnCount++

	// Run scam detection
	indicators := internal.ScamIndicators{}
	internal.ScamDetection(request.Message.Text, &indicators)

	// Update scam detection status
	if indicators.Score >= SCAM_THRESHOLD {
		session.Context.ScamDetected = true
	}

	// Add suspicious keywords
	for _, keyword := range indicators.Words {
		session.AddKeyword(keyword)
	}

	// Extract intelligence from message
	newIntel := internal.ExtractIntel(request.Message.Text, indicators.Score)
	session.Context.Intel = internal.MergeIntel(session.Context.Intel, newIntel)

	// Update state based on context
	session.Context.CurrentState = internal.GetState(session.Context)

	// Derive intent for response
	intent := internal.DeriveIntent(
		session.Context.CurrentState,
		session.Context.Intel,
		session.Context.TurnCount,
	)

	// Generate response
	reply := internal.GetResponse(intent)
	log.Println("reply: ", reply)

	// Check if we should send final callback
	if session.Context.CurrentState == internal.StateComplete {
		go sendFinalCallback(session)

		// Return completion response
		response := Response{
			Status: "success",
			Reply:  reply,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		// Clean up session
		store.Delete(request.SessionID)
		return
	}

	// Save session
	store.Update(session)

	// Return normal response
	response := Response{
		Status: "success",
		Reply:  reply,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendFinalCallback(session *internal.SessionData) {

	callbackURL := os.Getenv("CALLBACK_URL")
	if callbackURL == "" {
		callbackURL = "https://hackathon.guvi.in/api/updateHoneyPotFinalResult"
		log.Printf("Using default GUVI callback endpoint: %s", callbackURL)
	}

	// Build agent notes
	notes := buildAgentNotes(session)

	// Prepare final report matching GUVI format exactly
	finalReport := FinalResponse{
		SessionID:       session.SessionID,
		ScamDetect:      session.Context.ScamDetected,
		TotalMessagesEx: session.Context.TurnCount,
		ExtractIntel: ExtractedIntel{
			BankAccounts:       session.Context.Intel.Bank,
			UPIIds:             session.Context.Intel.UPI,
			PhishingLinks:      session.Context.Intel.Link,
			PhoneNumbers:       session.Context.Intel.Phone,
			SuspiciousKeywords: session.Keywords,
		},
		AgentNote: notes,
	}

	// Send callback
	jsonData, err := json.Marshal(finalReport)
	if err != nil {
		log.Printf("Error marshaling final report: %v", err)
		return
	}
	log.Println("Final report JSON: ", string(jsonData))

	resp, err := http.Post(callbackURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending callback: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Final callback sent successfully for session %s", session.SessionID)
	} else {
		log.Printf("Callback failed with status %d for session %s", resp.StatusCode, session.SessionID)
	}
}

func buildAgentNotes(session *internal.SessionData) string {
	var notes []string

	notes = append(notes, "Session completed after extracting sufficient intelligence.")

	if len(session.Keywords) > 0 {
		notes = append(notes, "Detected scam indicators: "+strings.Join(session.Keywords, ", "))
	}

	intelCount := len(session.Context.Intel.UPI) + len(session.Context.Intel.Phone) +
		len(session.Context.Intel.Link) + len(session.Context.Intel.Bank)

	notes = append(notes, "Total intelligence items extracted: "+string(rune(intelCount+'0')))

	return strings.Join(notes, " ")
}
