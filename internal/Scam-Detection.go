package internal

import (
	"regexp"
)

type ScamIndicators struct {
	Score          int
	Words         []string
}

 var regexUrgent = regexp.MustCompile(`(?i)\b(urgent|immediately|right\s*now|today|within\s*\d+\s*(minutes?|hours?)|final\s*warning)\b`)
var regexAccountThreatWords = regexp.MustCompile(`(?i)\b(account|upi|bank).*(blocked|suspended|disabled|closed|frozen)\b`)
var regexVerificationWords = regexp.MustCompile(`(?i)\b(verify|verification|kyc|re-?activate|update\s*kyc)\b`)
var regexPaymentWords = regexp.MustCompile(`(?i)\b(pay|payment|transfer|send|deposit)\b`)
var regexOTPWords = regexp.MustCompile(`(?i)\b(otp|one\s*time\s*password|pin|cvv|password)\b`)
var regexImpersonationWords = regexp.MustCompile(`(?i)\b(bank|sbi|hdfc|icici|axis|rbi|customer\s*care|support\s*team)\b`)

func ScamDetection(input string,  indicators *ScamIndicators) {


	if regexUrgent.MatchString(input){
		indicators.Score += 20
		indicators.Words = append(indicators.Words, "Urgent Request" )
	}
	if regexAccountThreatWords.MatchString(input) {
		indicators.Score += 30
		indicators.Words = append(indicators.Words, "Account Threat" )
	}
	if regexVerificationWords.MatchString(input) {
		indicators.Score += 30
		indicators.Words = append(indicators.Words, "Verification Request" )
	}
	if regexPaymentWords.MatchString(input) {
		indicators.Score += 30
		indicators.Words = append(indicators.Words, "Payment Request" )
	}
	if regexOTPWords.MatchString(input) {
		indicators.Score += 50
		indicators.Words = append(indicators.Words, "OTP Request" )
	}
	if regexImpersonationWords.MatchString(input) {
		indicators.Score += 20
		indicators.Words = append(indicators.Words, "Impersonation" )
	} 

}