package internal

import (
	"regexp"
)

type ScamIndicators struct {
	Score            int
	Words            []string
	HasUrgency       bool
	HasThreat        bool
	HasFinancial     bool
	HasCredential    bool
	HasImpersonation bool
}

var regexUrgent = regexp.MustCompile(`(?i)\b(urgent|immediately|right\s*now|today|within\s*\d+\s*(minutes?|hours?)|final\s*warning)\b`)
var regexAccountThreatWords = regexp.MustCompile(`(?i)\b(account|upi|bank).*(blocked|suspended|disabled|closed|frozen)\b`)
var regexVerificationWords = regexp.MustCompile(`(?i)\b(verify|verification|kyc|re-?activate|update\s*kyc)\b`)
var regexPaymentWords = regexp.MustCompile(`(?i)\b(pay|payment|transfer|send|deposit)\b`)
var regexOTPWords = regexp.MustCompile(`(?i)\b(otp|one\s*time\s*password|pin|cvv|password)\b`)
var regexImpersonationWords = regexp.MustCompile(`(?i)\b(bank|sbi|hdfc|icici|axis|rbi|customer\s*care|support\s*team)\b`)

func ScamDetection(input string, indicators *ScamIndicators) {

	if matches := regexUrgent.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 40
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasUrgency = true
	}
	if matches := regexAccountThreatWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 40
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasThreat = true
	}
	if matches := regexVerificationWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 20
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasFinancial = true
	}
	if matches := regexPaymentWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 25
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasFinancial = true
	}
	if matches := regexOTPWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 35
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasCredential = true
	}
	if matches := regexImpersonationWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 10
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasImpersonation = true
	}

}

// IsScam determines if the indicators represent a scam based on combination logic
// Requires: (Urgency OR Threat) AND (Financial OR Credential) AND score >= 100
func IsScam(indicators *ScamIndicators) bool {
	// Must have urgency/threat AND financial/credential indicators
	hasUrgencyOrThreat := indicators.HasUrgency || indicators.HasThreat
	hasFinancialOrCredential := indicators.HasFinancial || indicators.HasCredential

	// Combination logic: need multiple categories + score threshold
	return hasUrgencyOrThreat && hasFinancialOrCredential && indicators.Score >= 100
}
