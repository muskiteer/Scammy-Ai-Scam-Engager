# 🍯 Honeypot API — AI Scam Engagement System

## Description
A hybrid scam detection and engagement system that combines rule-based pattern matching with adaptive conversation strategies. The system uses multi-turn confidence scoring to detect scams, extracts intelligence through strategic questioning, and generates human-like responses using **Groq LLM API** to maintain engagement with scammers — wasting their time and gathering critical evidence.

---

## 🧰 Tech Stack
| Component | Technology |
|-----------|------------|
| **Language** | Go 1.x |
| **AI / LLM Provider** | Groq API (LLaMA-based models for fast, natural response generation) |
| **Pattern Matching** | Regular Expressions |
| **Architecture** | RESTful API with in-memory session management |
| **Deployment** | Render (Cloud Hosting) |
| **Build** | Go Modules |

---

## ⚙️ Setup Instructions

### Prerequisites
- Go 1.21 or higher installed
- A valid [Groq API](https://console.groq.com/) key
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Ai-Scam-Engagement
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export API_KEY="your-api-key"
   export CALLBACK_URL="https://your-callback-url.com"
   export GROQ_API_KEY="your-groq-api-key"
   ```

4. **Run the application**
   ```bash
   go run src/main.go
   # Or build and run
   go build -o build/scam-detector src/main.go
   ./build/scam-detector
   ```

5. **Verify the server is running**
   ```bash
   curl http://localhost:8080/health
   ```

---

## 🌐 API Endpoint

| Field | Value |
|-------|-------|
| **Base URL** | `https://scammy-ai-scam-engager.onrender.com` |
| **Engage Endpoint** | `POST /api/engage` |
| **Health Check** | `GET /health` |
| **Authentication** | `Sc_ODjFW5MFFFOW547mIUPolg1Qrc-BD8Ys` |

### Sample Request
```json
{
  "message": "Dear customer, your bank account has been suspended. Send your OTP to reactivate.",
  "session_id": "abc-123"
}
```

### Sample Response
```json
{
  "response": "Oh no, that sounds serious! Which bank is this regarding? I have accounts with multiple banks.",
  "session_id": "abc-123",
  "scam_detected": true,
  "confidence": 0.75
}
```

---

## 🧠 Approach

### 1. Scam Detection
- **Rule-based analysis** identifies urgency keywords, threats, financial requests, and impersonation attempts using curated regex patterns and keyword dictionaries
- **Confidence scoring** accumulates across multiple message turns — each detected scam indicator (e.g., *"act now"*, *"send money"*, *"your account will be blocked"*) adds weighted points to an overall scam confidence score
- **Threshold activation** triggers engagement mode once confidence exceeds **60%**, transitioning from passive detection to active scam engagement
- **Groq-powered contextual analysis** supplements rule-based detection by leveraging the **Groq LLM API** to understand nuanced scam tactics, interpret ambiguous messages, and validate scam intent when rule-based confidence is borderline — ensuring fewer false positives and smarter escalation decisions
- **Multi-category classification** detects various scam types including bank fraud, UPI fraud, phishing, lottery scams, tech support scams, and impersonation attempts
- **Conversation history awareness** analyzes the full conversation context (not just individual messages) to catch scammers who gradually escalate their tactics over multiple turns

### 2. Intelligence Extraction
- **Regex patterns** extract phone numbers, UPI IDs, bank accounts, email addresses, and phishing links from scammer messages
- **Intent-based questioning** strategically asks for missing information types — if a phone number is already captured, the system pivots to ask for a bank name or UPI ID
- **Session tracking** maintains full context across conversation turns, building a complete intelligence profile of the scammer
- **Data normalization** cleans and standardizes extracted data (e.g., phone number formats, URL deobfuscation)

### 3. Response Generation
- **Intent mapping** determines what type of question or response is needed based on the current conversation state and missing intelligence
- **Groq API integration** generates natural, human-like responses based on intent and conversation tone — the system prompts the Groq LLM with carefully crafted instructions to sound like a genuine, slightly naive victim
- **Adaptive strategy** balances information gathering with maintaining engagement (optimal engagement window: **8–15 turns**)
- **Tone matching** adjusts response style based on scam type — fearful for threat-based scams, excited for lottery scams, confused for tech support scams
- **Anti-detection measures** introduces natural delays, typos, and conversational fillers to avoid detection by sophisticated scammers

### 4. Engagement Metrics
- Tracks conversation duration, turn count, questions asked, and red flags identified
- Calculates scam type (`bank_fraud`, `upi_fraud`, `phishing`) and confidence level
- Submits final intelligence report with extracted data and engagement metrics
- Measures **time wasted** — the primary success metric for keeping scammers occupied

---

## 🔄 System Flow

```
Incoming Message
       │
       ▼
┌──────────────┐
│ Rule-Based   │──── Keywords, regex, threat patterns, Groq API
│ Analysis     │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Confidence   │──── Accumulates score across turns
│ Scoring      │
└──────┬───────┘
       │
       ▼
   Score > 60%? ─── No ──▶ Generic safe response
       │
      Yes
       │
       ▼
┌──────────────┐
│ Intelligence │──── Extract phone, UPI, email, links
│ Extraction   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Groq LLM     │──── Generate human-like engagement reply
│ Response Gen  │
└──────┬───────┘
       │
       ▼
   Response sent back to scammer
```

---

## 📊 Scam Categories Detected

| Category | Indicators |
|----------|------------|
| **Bank Fraud** | Account suspension, OTP requests, KYC updates |
| **UPI Fraud** | Payment requests, QR codes, refund scams |
| **Phishing** | Suspicious links, login page mimics, credential harvesting |
| **Lottery/Prize** | Congratulations messages, prize claims, advance fee requests |
| **Tech Support** | Virus warnings, remote access requests, software installation |
| **Impersonation** | Government official claims, bank representative claims |

---

## 🗂️ Project Structure

```
Ai-Scam-Engagement/
├── src/
│   ├── main.go              # Application entry point & HTTP server
│   ├── handler/             # Scam detection & confidence scoring
│   ├── internal/           # Response generation & strategy
│   ├── middleware/           # Intelligence extraction (regex)
│   └── session/              # Session management
├── build/                    # Compiled binaries
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency checksums
└── README.md
```

---

## 🚀 Why Groq?

The system uses [Groq](https://groq.com/) as its LLM provider for response generation because:

- **Ultra-low latency** — Groq's LPU (Language Processing Unit) delivers responses in milliseconds, critical for real-time scam engagement where delays feel unnatural
- **Cost-effective** — Generous free tier and affordable pricing for high-volume scam interception
- **High-quality output** — Runs LLaMA and Mixtral models that produce convincing, context-aware responses
- **Simple API** — OpenAI-compatible API format makes integration straightforward

---

## 📈 Key Metrics Tracked

| Metric | Description |
|--------|-------------|
| `turn_count` | Number of messages exchanged |
| `confidence` | Scam detection confidence (0.0 – 1.0) |
| `scam_type` | Classified scam category |
| `extracted_data` | Phone numbers, UPI IDs, emails, links found |
| `questions_asked` | Strategic questions posed to the scammer |
| `red_flags` | Specific scam indicators triggered |
| `engagement_duration` | Total time the scammer was kept engaged |

---

## 🛡️ Disclaimer

This project is built for **defensive cybersecurity purposes only**. It is designed to:
- Waste scammers' time, reducing their ability to target real victims
- Gather intelligence on scam operations for reporting to authorities
- Study scam tactics and improve detection systems

This tool should **not** be used for harassment, entrapment, or any illegal activity.
