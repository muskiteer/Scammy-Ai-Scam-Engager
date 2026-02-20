package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	EmailAddresses     []string `json:"emailAddresses"`
	CaseIDs            []string `json:"caseIDs,omitempty"`
	PolicyNumbers      []string `json:"policyNumbers,omitempty"`
	OrderNumbers       []string `json:"orderNumbers,omitempty"`
	CardNumbers        []string `json:"cardNumbers,omitempty"`
	IFSCCodes          []string `json:"ifscCodes,omitempty"`
	SuspiciousKeywords []string `json:"suspiciousKeywords"`
}

type EngagementMetrics struct {
	EngagementDurationSeconds int `json:"engagementDurationSeconds"`
	TotalMessagesExchanged    int `json:"totalMessagesExchanged"`
}

type FinalResponse struct {
	SessionID                 string            `json:"sessionId"`
	ScamDetect                bool              `json:"scamDetected"`
	TotalMessagesEx           int               `json:"totalMessagesExchanged"`
	EngagementDurationSeconds int               `json:"engagementDurationSeconds"`
	EngagementMetrics         EngagementMetrics `json:"engagementMetrics"`
	ExtractIntel              ExtractedIntel    `json:"extractedIntelligence"`
	AgentNote                 string            `json:"agentNotes"`
	ScamType                  string            `json:"scamType,omitempty"`
	ConfidenceLevel           string            `json:"confidenceLevel,omitempty"`
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

	// Track red flags identified
	for _, keyword := range indicators.Words {
		if !containsString(session.Context.RedFlagsIdentified, keyword) {
			session.Context.RedFlagsIdentified = append(session.Context.RedFlagsIdentified, keyword)
		}
	}

	// Add suspicious keywords
	for _, keyword := range indicators.Words {
		session.AddKeyword(keyword)
	}

	// Extract intelligence from current message
	newIntel := internal.ExtractIntel(request.Message.Text, indicators.Score)
	session.Context.Intel = internal.MergeIntel(session.Context.Intel, newIntel)

	// Also scan conversation history for any missed intel
	for _, msg := range request.ConvoHistory {
		if msg.Sender == "scammer" || msg.Sender == "user" {
			histIntel := internal.ExtractIntel(msg.Text, indicators.Score)
			session.Context.Intel = internal.MergeIntel(session.Context.Intel, histIntel)
			// Also run scam detection on history
			histIndicators := internal.ScamIndicators{}
			internal.ScamDetection(msg.Text, &histIndicators)
			if internal.IsScam(&histIndicators) {
				session.Context.ScamDetected = true
			}
			for _, keyword := range histIndicators.Words {
				session.AddKeyword(keyword)
				if !containsString(session.Context.RedFlagsIdentified, keyword) {
					session.Context.RedFlagsIdentified = append(session.Context.RedFlagsIdentified, keyword)
				}
			}
		}
	}

	// Log current intel status
	log.Printf("Session %s - Turn %d - Intel: UPI=%d, Phone=%d, Link=%d, Bank=%d, Email=%d",
		request.SessionID, session.Context.TurnCount,
		len(session.Context.Intel.UPI), len(session.Context.Intel.Phone),
		len(session.Context.Intel.Link), len(session.Context.Intel.Bank),
		len(session.Context.Intel.Email))

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
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskPhone:
		session.Context.AskCount.Phone++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskLink:
		session.Context.AskCount.Link++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskBank:
		session.Context.AskCount.Bank++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskEmail:
		session.Context.AskCount.Email++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskCaseID:
		session.Context.AskCount.CaseID++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskPolicyNumber:
		session.Context.AskCount.PolicyNumber++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskOrderNumber:
		session.Context.AskCount.OrderNumber++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskCardNumber:
		session.Context.AskCount.CardNumber++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskIFSCCode:
		session.Context.AskCount.IFSCCode++
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
		session.Context.InformationElicitations++
	case internal.IntentAskIdentity:
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
	case internal.IntentConfirmDetails:
		session.Context.QuestionsAsked++
		session.Context.InvestigativeQuestions++
	}

	reply := internal.GetResponse(intent)
	log.Println("reply: ", reply)

	// 15-second delay for engagement duration (well within 30s timeout)
	// 10 turns × ~18s (15s delay + network/processing) = 180+ seconds
	time.Sleep(15 * time.Second)

	// Send final callback ONLY at turn 10
	if session.Context.TurnCount >= 10 {
		log.Printf("Session %s - Turn 10 reached. Sending final callback.",
			request.SessionID)
		go sendFinalCallback(session)
		store.Delete(request.SessionID)
	} else {
		store.Update(session)
	}

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

	// Calculate engagement duration in seconds using session start time
	engagementDuration := int(time.Since(session.StartTime).Seconds())

	// Determine scam type based on keywords and indicators
	scamType := determineScamType(session)

	// Determine confidence level
	confidenceLevel := determineConfidenceLevel(session)

	finalReport := FinalResponse{
		SessionID:                 session.SessionID,
		ScamDetect:                session.Context.ScamDetected,
		TotalMessagesEx:           totalMessages,
		EngagementDurationSeconds: engagementDuration,
		EngagementMetrics: EngagementMetrics{
			EngagementDurationSeconds: engagementDuration,
			TotalMessagesExchanged:    totalMessages,
		},
		ExtractIntel: ExtractedIntel{
			BankAccounts:       session.Context.Intel.Bank,
			UPIIds:             session.Context.Intel.UPI,
			PhishingLinks:      session.Context.Intel.Link,
			PhoneNumbers:       session.Context.Intel.Phone,
			EmailAddresses:     session.Context.Intel.Email,
			CaseIDs:            session.Context.Intel.CaseIDs,
			PolicyNumbers:      session.Context.Intel.PolicyNumbers,
			OrderNumbers:       session.Context.Intel.OrderNumbers,
			CardNumbers:        session.Context.Intel.CardNumbers,
			IFSCCodes:          session.Context.Intel.IFSCCodes,
			SuspiciousKeywords: session.Keywords,
		},
		AgentNote:       notes,
		ScamType:        scamType,
		ConfidenceLevel: confidenceLevel,
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

		// Enhanced red flag identification - scan ALL conversation messages + keywords
		redFlags := []string{}
		allText := strings.ToLower(strings.Join(session.MessageHistory, " ") + " " + strings.Join(session.Keywords, " "))

		if strings.Contains(allText, "urgent") || strings.Contains(allText, "immediately") || strings.Contains(allText, "right now") || strings.Contains(allText, "final warning") {
			redFlags = append(redFlags, "urgency tactics")
		}
		if strings.Contains(allText, "otp") || strings.Contains(allText, "password") || strings.Contains(allText, "pin") || strings.Contains(allText, "cvv") {
			redFlags = append(redFlags, "credential harvesting")
		}
		if strings.Contains(allText, "blocked") || strings.Contains(allText, "suspended") || strings.Contains(allText, "closed") || strings.Contains(allText, "frozen") || strings.Contains(allText, "disabled") {
			redFlags = append(redFlags, "account threat")
		}
		if strings.Contains(allText, "verify") || strings.Contains(allText, "kyc") || strings.Contains(allText, "verification") || strings.Contains(allText, "re-activate") {
			redFlags = append(redFlags, "verification scam")
		}
		if strings.Contains(allText, "click") || strings.Contains(allText, "link") || strings.Contains(allText, "tap") || strings.Contains(allText, "visit") {
			redFlags = append(redFlags, "phishing attempt")
		}
		if strings.Contains(allText, "bank") || strings.Contains(allText, "sbi") || strings.Contains(allText, "hdfc") || strings.Contains(allText, "icici") || strings.Contains(allText, "rbi") || strings.Contains(allText, "customer care") {
			redFlags = append(redFlags, "impersonation")
		}
		if strings.Contains(allText, "payment") || strings.Contains(allText, "transfer") || strings.Contains(allText, "deposit") || strings.Contains(allText, "send money") {
			redFlags = append(redFlags, "financial request")
		}
		if strings.Contains(allText, "prize") || strings.Contains(allText, "lottery") || strings.Contains(allText, "winner") || strings.Contains(allText, "reward") {
			redFlags = append(redFlags, "lottery/prize scam")
		}
		if strings.Contains(allText, "police") || strings.Contains(allText, "arrest") || strings.Contains(allText, "legal") || strings.Contains(allText, "court") {
			redFlags = append(redFlags, "threat of legal action")
		}

		if len(redFlags) > 0 {
			uniqueFlags := deduplicateStrings(redFlags)
			notes = append(notes, "Red flags identified: "+strings.Join(uniqueFlags, ", ")+".")
		}
	} else {
		notes = append(notes, "No scam indicators detected.")
	}

	intelCount := len(session.Context.Intel.UPI) + len(session.Context.Intel.Phone) +
		len(session.Context.Intel.Link) + len(session.Context.Intel.Bank) + len(session.Context.Intel.Email) +
		len(session.Context.Intel.CaseIDs) + len(session.Context.Intel.PolicyNumbers) +
		len(session.Context.Intel.OrderNumbers) + len(session.Context.Intel.CardNumbers) + len(session.Context.Intel.IFSCCodes)

	if intelCount > 0 {
		notes = append(notes, "Successfully extracted intelligence through strategic engagement.")
	}

	if len(session.Keywords) > 0 {
		notes = append(notes, "Scammer used tactics: "+strings.Join(session.Keywords, ", ")+".")
	}

	return strings.Join(notes, " ")
}

func determineScamType(session *internal.SessionData) string {
	// Scan ALL text: keywords + full conversation history
	allText := strings.ToLower(strings.Join(session.Keywords, " ") + " " + strings.Join(session.MessageHistory, " "))

	// Check for bank fraud indicators
	if strings.Contains(allText, "bank") || strings.Contains(allText, "account") ||
		strings.Contains(allText, "blocked") || strings.Contains(allText, "suspended") ||
		strings.Contains(allText, "otp") {
		return "bank_fraud"
	}

	// Check for UPI fraud indicators
	if strings.Contains(allText, "upi") || strings.Contains(allText, "payment") ||
		strings.Contains(allText, "verify") || strings.Contains(allText, "kyc") {
		return "upi_fraud"
	}

	// Check for phishing indicators
	if len(session.Context.Intel.Link) > 0 || strings.Contains(allText, "click") ||
		strings.Contains(allText, "link") {
		return "phishing"
	}

	// Check for impersonation
	if strings.Contains(allText, "customer care") || strings.Contains(allText, "support team") ||
		strings.Contains(allText, "rbi") || strings.Contains(allText, "police") {
		return "impersonation"
	}

	// Default to generic scam if detected
	if session.Context.ScamDetected {
		return "generic_scam"
	}

	return "unknown"
}

func determineConfidenceLevel(session *internal.SessionData) string {
	if !session.Context.ScamDetected {
		return "low"
	}

	// Calculate confidence based on intel count and red flags
	intelCount := len(session.Context.Intel.UPI) + len(session.Context.Intel.Phone) +
		len(session.Context.Intel.Link) + len(session.Context.Intel.Bank) + len(session.Context.Intel.Email) +
		len(session.Context.Intel.CaseIDs) + len(session.Context.Intel.PolicyNumbers) +
		len(session.Context.Intel.OrderNumbers) + len(session.Context.Intel.CardNumbers) + len(session.Context.Intel.IFSCCodes)
	redFlagCount := len(session.Context.RedFlagsIdentified)

	// High confidence: 3+ red flags or 2+ intel items
	if redFlagCount >= 3 || intelCount >= 2 {
		return "high"
	}

	// Medium confidence: 1+ red flags or 1+ intel items
	if redFlagCount >= 1 || intelCount >= 1 {
		return "medium"
	}

	return "low"
}

func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func deduplicateStrings(items []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
