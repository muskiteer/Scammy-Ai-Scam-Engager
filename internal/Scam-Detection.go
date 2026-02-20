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
	HasLottery       bool
	HasTechSupport   bool
	HasGovtThreat    bool
}

// EXPANDED: Added more urgency patterns
var regexUrgent = regexp.MustCompile(`(?i)\b(urgent|immediately|right\s*now|today|within\s*\d+\s*(minutes?|hours?)|final\s*warning)\b`)
var regexAccountThreatWords = regexp.MustCompile(`(?i)\b(account|upi|bank).*(blocked|suspended|disabled|closed|frozen)\b`)
var regexVerificationWords = regexp.MustCompile(`(?i)\b(verify|verification|kyc|re-?activate|update\s*kyc)\b`)
var regexPaymentWords = regexp.MustCompile(`(?i)\b(pay|payment|transfer|send|deposit)\b`)
var regexOTPWords = regexp.MustCompile(`(?i)\b(otp|one\s*time\s*password|pin|cvv|password)\b`)
var regexImpersonationWords = regexp.MustCompile(`(?i)\b(bank|sbi|hdfc|icici|axis|rbi|customer\s*care|support\s*team)\b`)

// Action words
var regexActionWords = regexp.MustCompile(`(?i)\b(click|tap|call|dial|visit|open\s*link|contact|reply|download)\b`)

// Lottery and prize scam patterns
var regexLotteryWords = regexp.MustCompile(`(?i)\b(prize|lottery|winner|won|congratulations|reward|cashback|refund|bonus|gift)\b`)

// Tech support scam patterns
var regexTechSupportWords = regexp.MustCompile(`(?i)\b(virus|malware|hacked|compromised|remote\s*access|technical\s*support|install\s*software)\b`)

// Government threat and legal intimidation patterns
var regexGovtThreatWords = regexp.MustCompile(`(?i)\b(police|cbi|enforcement|income\s*tax|court|arrest|warrant|legal\s*action|government\s*official)\b`)

// Delivery and parcel scam patterns
var regexDeliveryWords = regexp.MustCompile(`(?i)\b(parcel|package|customs|courier|shipment|detained|delivery\s*fee)\b`)

// Job and investment scam patterns
var regexJobScamWords = regexp.MustCompile(`(?i)\b(job\s*offer|work\s*from\s*home|part[\s\-]time|earn\s*money|investment\s*return|profit\s*guarantee|easy\s*income)\b`)

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
	// Action word detection
	if matches := regexActionWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 15
		indicators.Words = append(indicators.Words, matches[0])
	}
	// Lottery/prize scam detection
	if matches := regexLotteryWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 30
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasLottery = true
		indicators.HasFinancial = true
	}
	// Tech support scam detection
	if matches := regexTechSupportWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 25
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasTechSupport = true
	}
	// Government threat / legal intimidation detection
	if matches := regexGovtThreatWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 35
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasGovtThreat = true
		indicators.HasThreat = true
	}
	// Delivery/parcel scam detection
	if matches := regexDeliveryWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 20
		indicators.Words = append(indicators.Words, matches[0])
		indicators.HasFinancial = true
	}
	// Job/investment scam detection
	if matches := regexJobScamWords.FindStringSubmatch(input); len(matches) > 0 {
		indicators.Score += 20
		indicators.Words = append(indicators.Words, matches[0])
	}
}

// IsScam - OPTIMIZED: Lower thresholds to maximize scam detection
func IsScam(indicators *ScamIndicators) bool {
	hasUrgencyOrThreat := indicators.HasUrgency || indicators.HasThreat
	hasFinancialOrCredential := indicators.HasFinancial || indicators.HasCredential

	// CASE 1: Classic pattern (lowered threshold from 60 to 40)
	if hasUrgencyOrThreat && hasFinancialOrCredential && indicators.Score >= 40 {
		return true
	}

	// CASE 2: Impersonation + credential (bank asking for OTP) - lowered from 45 to 30
	if indicators.HasImpersonation && indicators.HasCredential && indicators.Score >= 30 {
		return true
	}

	// CASE 3: Impersonation + financial (fake bank asking for payment) - lowered from 50 to 35
	if indicators.HasImpersonation && indicators.HasFinancial && indicators.Score >= 35 {
		return true
	}

	// CASE 4: Impersonation + urgency/threat - lowered from 55 to 40
	if indicators.HasImpersonation && hasUrgencyOrThreat && indicators.Score >= 40 {
		return true
	}

	// CASE 5: Any combination with decent score - lowered from 100 to 60
	if indicators.Score >= 60 {
		return true
	}

	// CASE 6: Financial + credential even without urgency
	if indicators.HasFinancial && indicators.HasCredential && indicators.Score >= 40 {
		return true
	}

	// CASE 7: Lottery/prize + financial (classic cashback/refund scam)
	if indicators.HasLottery && indicators.HasFinancial {
		return true
	}

	// CASE 8: Government threat with a significant score
	if indicators.HasGovtThreat && indicators.Score >= 35 {
		return true
	}

	// CASE 9: Tech support scam with urgency or high score
	if indicators.HasTechSupport && (indicators.HasUrgency || indicators.Score >= 40) {
		return true
	}

	return false
}
