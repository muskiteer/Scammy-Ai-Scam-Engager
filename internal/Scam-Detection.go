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

// EXPANDED: Added more urgency patterns
var regexUrgent = regexp.MustCompile(`(?i)\b(urgent|immediately|right\s*now|today|within\s*\d+\s*(minutes?|hours?)|final\s*warning)\b`)
var regexAccountThreatWords = regexp.MustCompile(`(?i)\b(account|upi|bank).*(blocked|suspended|disabled|closed|frozen)\b`)
var regexVerificationWords = regexp.MustCompile(`(?i)\b(verify|verification|kyc|re-?activate|update\s*kyc)\b`)
var regexPaymentWords = regexp.MustCompile(`(?i)\b(pay|payment|transfer|send|deposit)\b`)
var regexOTPWords = regexp.MustCompile(`(?i)\b(otp|one\s*time\s*password|pin|cvv|password)\b`)
var regexImpersonationWords = regexp.MustCompile(`(?i)\b(bank|sbi|hdfc|icici|axis|rbi|customer\s*care|support\s*team)\b`)

// NEW: Action words
var regexActionWords = regexp.MustCompile(`(?i)\b(click|tap|call|dial|visit|open\s*link|contact|reply|download)\b`)

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
        indicators.Score += 15 // Increased from 10
        indicators.Words = append(indicators.Words, matches[0])
        indicators.HasImpersonation = true
    }
    // NEW: Action word detection
    if matches := regexActionWords.FindStringSubmatch(input); len(matches) > 0 {
        indicators.Score += 15
        indicators.Words = append(indicators.Words, matches[0])
    }
}

// IsScam - FIXED: Lower threshold + impersonation logic
func IsScam(indicators *ScamIndicators) bool {
    hasUrgencyOrThreat := indicators.HasUrgency || indicators.HasThreat
    hasFinancialOrCredential := indicators.HasFinancial || indicators.HasCredential

    // CASE 1: Classic pattern (lowered threshold from 100 to 60)
    if hasUrgencyOrThreat && hasFinancialOrCredential && indicators.Score >= 60 {
        return true
    }

    // CASE 2: Impersonation + credential (bank asking for OTP)
    if indicators.HasImpersonation && indicators.HasCredential && indicators.Score >= 45 {
        return true
    }

    // CASE 3: Impersonation + financial (fake bank asking for payment)
    if indicators.HasImpersonation && indicators.HasFinancial && indicators.Score >= 50 {
        return true
    }

    // CASE 4: Impersonation + urgency/threat
    if indicators.HasImpersonation && hasUrgencyOrThreat && indicators.Score >= 55 {
        return true
    }

    // CASE 5: Very high score
    if indicators.Score >= 100 {
        return true
    }

    return false
}