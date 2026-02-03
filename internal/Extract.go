package internal

import (
	"regexp"
	"strings"
)

var (
	UPIRegex          = regexp.MustCompile(`\b[a-zA-Z0-9.\-_]{2,}@[a-zA-Z]{2,}\b`)
	PhoneRegex        = regexp.MustCompile(`(\+91[\s-]?)?[6-9]\d{9}\b`)
	PhishingLinkRegex = regexp.MustCompile(`https?://[^\s]+`)
	BankAccountRex    = regexp.MustCompile(`(?i)(?:account\s*(?:number|no\.?|#)?[:\s-]*)?\\b\d{10,18}\b`)
)

var trustedDomains = []string{"google.com", ".gov.in", "nic.in", "sbi.co.in", "icicibank.com", "hdfcbank.com", "axisbank.com"}

// ExtractIntel extracts intelligence data from input text
func ExtractIntel(input string, confidence int) Intel {
	intel := Intel{
		UPI:   []string{},
		Phone: []string{},
		Link:  []string{},
		Bank:  []string{},
	}

	// Extract all UPI IDs
	upiMatches := UPIRegex.FindAllString(input, -1)
	intel.UPI = append(intel.UPI, upiMatches...)

	// Extract all phone numbers first (to avoid them being treated as bank accounts)
	phoneMatches := PhoneRegex.FindAllString(input, -1)
	phoneSet := make(map[string]bool)
	for _, phone := range phoneMatches {
		intel.Phone = append(intel.Phone, phone)
		// Normalize phone to just digits for comparison
		normalized := strings.ReplaceAll(strings.ReplaceAll(phone, "+91", ""), " ", "")
		normalized = strings.ReplaceAll(normalized, "-", "")
		phoneSet[normalized] = true
		// Also add with +91 prefix digits
		if strings.HasPrefix(phone, "+91") {
			phoneSet[normalized[2:]] = true // Last 10 digits
		}
	}

	// Extract all links and filter out trusted domains
	linkMatches := PhishingLinkRegex.FindAllString(input, -1)
	for _, link := range linkMatches {
		isTrusted := false
		linkLower := strings.ToLower(link)
		for _, domain := range trustedDomains {
			if strings.Contains(linkLower, domain) {
				isTrusted = true
				break
			}
		}
		if !isTrusted {
			intel.Link = append(intel.Link, link)
		}
	}

	// Extract all bank account numbers (excluding phone numbers)
	bankMatches := BankAccountRex.FindAllString(input, -1)
	for _, bank := range bankMatches {
		// Extract just the digits
		digits := extractDigits(bank)
		// Skip if it's a phone number (exactly 10 digits) or less than 10 digits
		if len(digits) >= 10 && len(digits) != 10 && !phoneSet[digits] {
			intel.Bank = append(intel.Bank, bank)
		}
	}

	return intel
}

// extractDigits returns only the digit characters from a string
func extractDigits(s string) string {
	var result strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// MergeIntel combines two Intel structs, avoiding duplicates
func MergeIntel(existing Intel, new Intel) Intel {
	merged := Intel{
		UPI:   deduplicate(append(existing.UPI, new.UPI...)),
		Phone: deduplicate(append(existing.Phone, new.Phone...)),
		Link:  deduplicate(append(existing.Link, new.Link...)),
		Bank:  deduplicate(append(existing.Bank, new.Bank...)),
	}
	return merged
}

// deduplicate removes duplicate strings from a slice
func deduplicate(items []string) []string {
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
