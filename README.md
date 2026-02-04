# AI Scam Engagement Backend

Agentic Honey-Pot for Scam Detection & Intelligence Extraction

A deterministic, rule-based backend service in Go that acts as an agentic scam honeypot. The system engages with potential scammers, detects scam patterns using regex and scoring, and intelligently extracts critical information like UPI IDs, phone numbers, bank accounts, and phishing links.

> **🎉 Recent Improvements (Feb 2026):** Enhanced output quality with accurate message counting, normalized intelligence data, actual scammer phrase extraction, and comprehensive innocent sender testing. See [IMPROVEMENTS.md](IMPROVEMENTS.md) for details.


## Features

- **Rule-Based Scam Detection**: No ML/LLM required - uses regex patterns and scoring for signals like urgency, account threats, payment requests, OTPs, and phishing links
- **State Machine Architecture**: Manages conversation flow through INIT → ENGAGING → INTEL_EXTRACT → COMPLETE states
- **Intelligence Extraction**: Automatically identifies and collects UPI IDs, bank accounts, phone numbers, phishing links, and suspicious keywords
- **Intent-Driven Responses**: Derives conversational intent based on state and missing intelligence to keep scammers engaged
- **Session Management**: Per-session state tracking with thread-safe in-memory storage
- **Callback Webhooks**: Sends comprehensive final report when sufficient intelligence is gathered
- **Human-Like Responses**: Uses predefined templates with random selection for natural conversation

## Architecture

### State Machine
- **INIT**: Initial state, neutral engagement
- **ENGAGING**: Scam detected (score ≥ 50), building rapport
- **INTEL_EXTRACT**: Actively gathering missing intelligence
- **COMPLETE**: Sufficient data collected, triggers final callback

### Intent System
- **CONFIRM_DETAILS**: Ask for clarification about the scam
- **ASK_UPI**: Request UPI ID information
- **ASK_PHONE**: Request phone number
- **ASK_LINK**: Request verification links
- **STALL**: Buy time while extracting information
- **NEUTRAL**: Acknowledgment responses

### Scam Detection Scoring
- Urgency indicators: +20 points
- Account threats: +30 points
- Verification requests: +30 points
- Payment requests: +30 points
- OTP/password requests: +50 points
- Impersonation: +20 points
- **Threshold**: 50 points = scam detected

## API Endpoints

### POST /api/engage
Main conversation endpoint. Receives messages and returns bot responses.

**Request Body:**
```json
{
  "sessionId": "unique-session-id",
  "message": {
    "sender": "user",
    "text": "Your account will be blocked. Send OTP to verify.",
    "timestamp": "2026-02-02T10:30:00Z"
  },
  "conversationHistory": [
    {
      "sender": "user",
      "text": "Previous message",
      "timestamp": "2026-02-02T10:29:00Z"
    }
  ],
  "metadata": {
    "channel": "whatsapp",
    "language": "en",
    "locale": "en-IN"
  }
}
```

**Response:**
```json
{
  "status": "ok",
  "reply": "Why will my account be blocked?",
  "state": "ENGAGING"
}
```

When conversation is complete:
```json
{
  "status": "complete",
  "reply": "I am checking this now, please wait.",
  "state": "COMPLETE"
}
```

### GET /health
Health check endpoint.

**Response:** `200 OK` with body `"OK"`

## Callback Webhook

When a session reaches the COMPLETE state, a final report is automatically sent to the GUVI evaluation endpoint.

**Callback Endpoint (Default):**
```
POST https://hackathon.guvi.in/api/updateHoneyPotFinalResult
```

**Callback Payload:**
```json
{
  "sessionId": "unique-session-id",
  "scamDetected": true,
  "totalMessagesExchanged": 8,
  "extractedIntelligence": {
    "bankAccounts": ["1234567890123"],
    "upiIds": ["scammer@paytm", "fraud@ybl"],
    "phishingLinks": ["http://fake-bank-verify.com/login"],
    "phoneNumbers": ["9876543210", "+919123456789"],
    "suspiciousKeywords": ["Urgent Request", "Account Threat", "OTP Request"]
  },
  "agentNotes": "Session completed after extracting sufficient intelligence. Detected scam indicators: Urgent Request, Account Threat, OTP Request. Total intelligence items extracted: 5"
}
```

> **Note**: The callback URL defaults to the GUVI endpoint. You can override it by setting `CALLBACK_URL` in your `.env` file for testing purposes.

## Installation & Setup

### Prerequisites
- Go 1.22 or higher
- Git

### Clone and Install
```bash
git clone <repository-url>
cd Ai-Scam-Engagement
go mod download
```

### Configuration
Create a `.env` file (optional):
```bash
cp .env.example .env
```

Edit `.env` to configure:
```env
# API Key for authentication (REQUIRED for production)
API_KEY=your-secret-api-key-here

# Callback URL (defaults to GUVI endpoint if not set)
CALLBACK_URL=https://hackathon.guvi.in/api/updateHoneyPotFinalResult

# Server port (optional, defaults to 8080)
PORT=8080
```

### Run the Server
```bash
go run main.go
```

Server starts on `http://localhost:8080`

## Usage Examples

### Example 1: Scam Detection and Engagement
```bash
curl -X POST http://localhost:8080/api/engage \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "test-session-1",
    "message": {
      "sender": "user",
      "text": "URGENT! Your SBI account will be blocked in 1 hour. Verify at http://fake-sbi.com",
      "timestamp": "2026-02-02T10:00:00Z"
    },
    "conversationHistory": [],
    "metadata": {
      "channel": "whatsapp",
      "language": "en",
      "locale": "en-IN"
    }
  }'
```

**Response:**
```json
{
  "status": "success",
  "reply": "What exactly is the problem with my account?",
}
```

### Example 2: Intelligence Extraction
```bash
curl -X POST http://localhost:8080/api/engage \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "test-session-1",
    "message": {
      "sender": "user",
      "text": "Send payment to scammer@paytm or call 9876543210 immediately",
      "timestamp": "2026-02-02T10:02:00Z"
    },
    "conversationHistory": [],
    "metadata": {
      "channel": "whatsapp",
      "language": "en",
      "locale": "en-IN"
    }
  }'
```

**Response:**
```json
{
  "status": "success",
  "reply": "Can you send the link here?",
}
```

## API Authentication

All requests to `/api/engage` must include authentication via the `x-api-key` header:

```bash
curl -X POST http://localhost:8080/api/engage \
  -H "Content-Type: application/json" \
  -H "x-api-key: your-secret-api-key" \
  -d '...'
```

Set your API key in the `.env` file. If no API key is configured, authentication is bypassed (not recommended for production).

## HTTPS Deployment

For GUVI Hackathon submission, your endpoint must use HTTPS. See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed instructions on:

- Deploying to cloud platforms (Railway, Render, Fly.io)
- Setting up reverse proxy with Nginx/Caddy
- Configuring SSL certificates
- Docker containerization
- Production best practices

## Project Structure

```
.
├── main.go                    # Entry point
├── go.mod                     # Go module definition
├── handler/
│   └── handler.go            # HTTP handlers and business logic
├── internal/
│   ├── Scam-Detection.go     # Rule-based scam detection with scoring
│   ├── Extract.go            # Intelligence extraction (UPI, phone, links, etc.)
│   ├── intent.go             # State machine and intent derivation
│   ├── responses.go          # Response templates and generation
│   ├── session.go            # Session state management
│   └── parsing.go            # Input parsing utilities
└── routes/
    └── routes.go             # Route definitions
```

## Scam Detection Patterns

The system uses regex patterns to detect:

- **UPI IDs**: `username@bank` format (e.g., `scammer@paytm`)
- **Phone Numbers**: Indian format with optional +91 prefix (e.g., `9876543210`)
- **Phishing Links**: HTTP/HTTPS URLs excluding trusted domains
- **Bank Accounts**: Account numbers with 10-18 digits
- **Urgency Keywords**: "urgent", "immediately", "right now", "final warning"
- **Account Threats**: "blocked", "suspended", "disabled", "frozen"
- **Verification Requests**: "verify", "KYC", "reactivate"
- **OTP Requests**: "OTP", "PIN", "CVV", "password"
- **Impersonation**: Bank names, "customer care", "support team"

## Development

### Running Tests

#### Test Scam Detection (Full Intelligence Extraction):
```bash
export API_KEY="YOUR_SECRET_API_KEY"
./test_complete_intel.sh
```

#### Test Innocent Sender (False Positive Check):
```bash
export API_KEY="YOUR_SECRET_API_KEY"
./test_innocent_sender.sh
```

#### Run Unit Tests:
```bash
go test ./...
```

See [QUICK_REFERENCE.md](QUICK_REFERENCE.md) for more testing details.

### Building for Production
```bash
go build -o scam-honeypot
./scam-honeypot
```

## Documentation

- **[IMPROVEMENTS.md](IMPROVEMENTS.md)** - Recent quality improvements and refinements
- **[BEFORE_AFTER_COMPARISON.md](BEFORE_AFTER_COMPARISON.md)** - Output examples showing improvements
- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Quick start guide and testing checklist

## Security Considerations

- Sessions are stored in-memory and cleared after completion
- No persistent storage of sensitive data by default
- Configure CALLBACK_URL to send data to your secure backend
- Consider adding authentication for production deployments
- Rate limiting recommended for public-facing deployments

  