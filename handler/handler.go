package handler

import (
	"bytes"
	"encoding/json"
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
	Sender string `json:"sender"`
	Text   string `json:"text"`
	// Timestamp is optional and can be any format - we don't use it internally
	Timestamp interface{} `json:"timestamp,omitempty"`
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
	log.Println("/health ->Health check endpoint hit")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func StartConvo(w http.ResponseWriter, r *http.Request) {
	log.Println("/api/engage ->StartConvo endpoint hit")
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
	// body, _ := io.ReadAll(r.Body)
	// log.Println("Received request body: ", string(body))

	var request Request
	err := json.NewDecoder(r.Body).Decode(&request)
	log.Println("Received request: ", request)
	// log.Println("here 1 ")
	if err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	// log.Println("here 2 ")
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

	// Update scam detection status using combination logic
	if internal.IsScam(&indicators) {
		session.Context.ScamDetected = true
	}

	// Add suspicious keywords
	for _, keyword := range indicators.Words {
		session.AddKeyword(keyword)
	}

	// Extract intelligence from message
	newIntel := internal.ExtractIntel(request.Message.Text, indicators.Score)
	session.Context.Intel = internal.MergeIntel(session.Context.Intel, newIntel)

	// Log current intel status
	log.Printf("Session %s - Turn %d - Intel: UPI=%d, Phone=%d, Link=%d, Bank=%d",
		request.SessionID, session.Context.TurnCount,
		len(session.Context.Intel.UPI), len(session.Context.Intel.Phone),
		len(session.Context.Intel.Link), len(session.Context.Intel.Bank))

	// Update state based on context
	session.Context.CurrentState = internal.GetState(session.Context)
	log.Printf("Session %s - Current State: %s", request.SessionID, session.Context.CurrentState)

	// Derive intent for response
	intent := internal.DeriveIntent(
		session.Context.CurrentState,
		session.Context.Intel,
		session.Context.TurnCount,
		session.Context.AskCount,
	)

	// Increment ask count based on intent
	switch intent {
	case internal.IntentAskUPI:
		session.Context.AskCount.UPI++
	case internal.IntentAskPhone:
		session.Context.AskCount.Phone++
	case internal.IntentAskLink:
		session.Context.AskCount.Link++
	case internal.IntentAskBank:
		session.Context.AskCount.Bank++
	}

	reply := internal.GetResponse(intent)
	log.Println("reply: ", reply)

	// Check if conversation should end
	if session.Context.CurrentState == internal.StateComplete || session.Context.TurnCount >= 15 {
		// If we have all intel, give a conclusive response
		hasAllIntel := len(session.Context.Intel.UPI) > 0 && len(session.Context.Intel.Phone) > 0 &&
			len(session.Context.Intel.Link) > 0 && len(session.Context.Intel.Bank) > 0

		if hasAllIntel {
			reply = "Thank you for the information. I will verify everything and get back to you shortly."
		}

		log.Printf("Session %s - Ending conversation. State: %s, Turns: %d",
			request.SessionID, session.Context.CurrentState, session.Context.TurnCount)

		go sendFinalCallback(session)

		response := Response{
			Status: "success",
			Reply:  reply,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		store.Delete(request.SessionID)
		return
	}

	store.Update(session)

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

	notes := buildAgentNotes(session)

	// Calculate total messages: scammer messages (TurnCount) + agent responses (TurnCount)
	totalMessages := session.Context.TurnCount * 2

	finalReport := FinalResponse{
		SessionID:       session.SessionID,
		ScamDetect:      session.Context.ScamDetected,
		TotalMessagesEx: totalMessages,
		ExtractIntel: ExtractedIntel{
			BankAccounts:       session.Context.Intel.Bank,
			UPIIds:             session.Context.Intel.UPI,
			PhishingLinks:      session.Context.Intel.Link,
			PhoneNumbers:       session.Context.Intel.Phone,
			SuspiciousKeywords: session.Keywords,
		},
		AgentNote: notes,
	}

	jsonData, err := json.Marshal(finalReport)
	if err != nil {
		log.Printf("Error marshaling final report: %v", err)
		return
	}
	log.Println("+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+_+")
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

	if session.Context.ScamDetected {
		notes = append(notes, "Scam detected with high confidence.")
	} else {
		notes = append(notes, "No scam indicators detected.")
	}

	intelCount := len(session.Context.Intel.UPI) + len(session.Context.Intel.Phone) +
		len(session.Context.Intel.Link) + len(session.Context.Intel.Bank)

	if intelCount > 0 {
		notes = append(notes, "Successfully extracted intelligence through strategic engagement.")
	}

	if len(session.Keywords) > 0 {
		notes = append(notes, "Scammer used tactics: "+strings.Join(session.Keywords, ", ")+".")
	}

	return strings.Join(notes, " ")
}
