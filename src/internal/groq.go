package internal

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
)

type GroqRequest struct {
    Model    string        `json:"model"`
    Messages []GroqMessage `json:"messages"`
}

type GroqMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type GroqResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}

func GenerateGroqResponse(intent Intent, scammerMessage string, conversationHistory []string, turnCount int) string {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        log.Println("GROQ_API_KEY not set, falling back to static response")
        return GetResponse(intent)
    }

    systemPrompt := `You are playing the role of a naive, trusting elderly person who has received a suspicious message. Your goal is to:
1. Sound genuinely concerned and cooperative (not suspicious)
2. Ask the scammer for specific information based on the INTENT provided
3. Keep responses SHORT (1-2 sentences max)
4. Never reveal you know it's a scam
5. Act slightly confused to encourage the scammer to share more details
6. Reference things the scammer said to sound engaged

IMPORTANT RULES:
- Never say "scam", "fraud", "fake", "suspicious"
- Sound like a real worried person
- Always end with a question or request for more info
- Keep it under 40 words`

    intentInstructions := map[Intent]string{
        IntentConfirmDetails:  "Ask clarifying questions about what they said. Sound worried. Ask them to explain the problem again.",
        IntentAskPhone:        "Say you want to call them back for safety. Ask for their phone number or direct line.",
        IntentAskUPI:          "Say you're ready to cooperate. Ask which UPI ID you should use or ask them to confirm theirs.",
        IntentAskBank:         "Say you want to verify. Ask them to confirm the account number they're referring to.",
        IntentAskLink:         "Say you're not sure which link to use. Ask them to share the correct link or website.",
        IntentAskEmail:        "Say you want to send documents. Ask for their official email address.",
        IntentAskIdentity:     "Ask for their employee ID, department name, supervisor name, or office address to verify their identity.",
        IntentAskCaseID:       "Ask what is the reference number or case ID for this matter so you can track it.",
        IntentAskPolicyNumber: "Ask them to confirm the policy number or insurance details they are referring to.",
        IntentAskOrderNumber:  "Ask them to share the order number or booking reference so you can check.",
        IntentAskCardNumber:   "Say you have multiple cards. Ask them which card number they are referring to.",
        IntentAskIFSCCode:     "Ask them to confirm the IFSC code or branch details for verification.",
        IntentStall:           "Say you're looking for the information they asked for. Buy time. Sound cooperative but slow.",
        IntentNeutral:         "Respond naturally to what they said. Sound concerned and ask a follow-up question.",
    }

    instruction := intentInstructions[intent]
    if instruction == "" {
        instruction = "Respond naturally and ask a follow-up question."
    }

    // Build conversation context
    var contextBuilder strings.Builder
    contextBuilder.WriteString("INTENT: " + string(intent) + "\n")
    contextBuilder.WriteString("INSTRUCTION: " + instruction + "\n")
    contextBuilder.WriteString("TURN: " + string(rune('0'+turnCount)) + " of 10\n\n")

    if len(conversationHistory) > 0 {
        contextBuilder.WriteString("Recent conversation:\n")
        start := 0
        if len(conversationHistory) > 6 {
            start = len(conversationHistory) - 6
        }
        for i := start; i < len(conversationHistory); i++ {
            if i%2 == 0 {
                contextBuilder.WriteString("Scammer: " + conversationHistory[i] + "\n")
            } else {
                contextBuilder.WriteString("Me: " + conversationHistory[i] + "\n")
            }
        }
    }
    contextBuilder.WriteString("\nScammer's latest message: " + scammerMessage + "\n")
    contextBuilder.WriteString("\nRespond as the naive victim (1-2 sentences, end with a question):")

    messages := []GroqMessage{
        {Role: "system", Content: systemPrompt},
        {Role: "user", Content: contextBuilder.String()},
    }

    reqBody := GroqRequest{
        Model:    "llama-3.1-8b-instant",
        Messages: messages,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        log.Printf("Error marshaling Groq request: %v", err)
        return GetResponse(intent)
    }

    req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error creating Groq request: %v", err)
        return GetResponse(intent)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error calling Groq API: %v", err)
        return GetResponse(intent)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading Groq response: %v", err)
        return GetResponse(intent)
    }

    var groqResp GroqResponse
    err = json.Unmarshal(body, &groqResp)
    if err != nil {
        log.Printf("Error parsing Groq response: %v", err)
        return GetResponse(intent)
    }

    if len(groqResp.Choices) > 0 && groqResp.Choices[0].Message.Content != "" {
        reply := strings.TrimSpace(groqResp.Choices[0].Message.Content)
        // Remove any quotes the LLM might wrap the response in
        reply = strings.Trim(reply, "\"'")
        log.Printf("Groq response: %s", reply)
        return reply
    }

    return GetResponse(intent)
}