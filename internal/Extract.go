package internal

import (
	"regexp"
	"strings"
)

// EXPANDED: Multiple UPI patterns to catch different formats
var (
	// UPI patterns: user@bank, user.name@upi, etc.
	UPIRegex = regexp.MustCompile(`(?i)\b[a-zA-Z0-9][a-zA-Z0-9.\-_]{1,49}@[a-zA-Z]{2,}(?:bank|upi|paytm|ybl|okaxis|okhdfcbank|oksbi|apl|axl|ibl|sbi|hdfc|icici|axis|kotak|pnb|bob|union|canara|indian|indus|federal|rbl|yes|idfc|bandhan|au|equitas|ujjivan)\b`)

	// EXPANDED: Phone patterns - multiple formats
	// +91XXXXXXXXXX, +91-XXXX-XXXXXX, 91XXXXXXXXXX, XXXXXXXXXX, XXX-XXX-XXXX, etc.
	PhoneRegex = regexp.MustCompile(`(?:(?:\+|00)?91[\s\-\.]?)?[6-9]\d{2}[\s\-\.]?\d{3}[\s\-\.]?\d{4}\b|\b[6-9]\d{9}\b`)

	// EXPANDED: Phishing links - http, https, www, shortened URLs
	PhishingLinkRegex = regexp.MustCompile(`(?i)(?:https?://|www\.)[^\s<>"'\)\]\}]+|(?:bit\.ly|goo\.gl|tinyurl\.com|t\.co|is\.gd|buff\.ly|ow\.ly|rebrand\.ly|shorturl\.at)/[^\s<>"'\)\]\}]+`)

	// EXPANDED: Bank account patterns - with/without labels
	BankAccountRex = regexp.MustCompile(`(?i)(?:(?:a/?c|account|acct)[\s\.\-:]*(?:no|number|num|#)?[\s\.\-:]*)?(\d{9,18})`)

	// NEW: Card number pattern (16 digits with optional spaces/dashes)
	CardNumberRegex = regexp.MustCompile(`\b(?:\d{4}[\s\-]?){3}\d{4}\b`)

	// NEW: IFSC code pattern
	IFSCRegex = regexp.MustCompile(`(?i)\b[A-Z]{4}0[A-Z0-9]{6}\b`)

	// NEW: Email pattern for extraction
	EmailRegex = regexp.MustCompile(`(?i)\b[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}\b`)

	// NEW: Generic ID patterns (employee ID, reference ID, etc.)
	ReferenceIDRegex = regexp.MustCompile(`(?i)(?:ref(?:erence)?|id|ticket|case|complaint)[\s\.\-:#]*([A-Z0-9]{6,20})`)
)

// EXPANDED: More trusted domains to filter
var trustedDomains = []string{
	"google.com", "google.co.in", ".gov.in", "nic.in",
	"sbi.co.in", "onlinesbi.com", "icicibank.com", "hdfcbank.com", "axisbank.com",
	"kotak.com", "pnbindia.in", "bankofbaroda.in", "canarabank.com",
	"rbi.org.in", "npci.org.in", "upi.org", "bhimupi.org",
	"paytm.com", "phonepe.com", "gpay.com", "amazonpay.in",
	"microsoft.com", "apple.com", "amazon.in", "flipkart.com",
}

// Common UPI suffixes for validation
var upiSuffixes = []string{
	"@ybl", "@paytm", "@okaxis", "@okhdfcbank", "@oksbi", "@apl", "@axl", "@ibl",
	"@sbi", "@hdfc", "@icici", "@axis", "@kotak", "@pnb", "@bob", "@upi",
	"@axisbank", "@hdfcbank", "@sbiupi", "@icicipay", "@aubank", "@equitas",
	"@federal", "@indus", "@rbl", "@yes", "@idfc", "@bandhan", "@ujjivan",
}

// ExtractIntel extracts intelligence data from input text
func ExtractIntel(input string, confidence int) Intel {
	intel := Intel{
		UPI:   []string{},
		Phone: []string{},
		Link:  []string{},
		Bank:  []string{},
	}

	// Normalize input for better matching
	normalizedInput := strings.ToLower(input)

	// ============ EXTRACT UPI IDs ============
	// Method 1: Standard UPI regex
	upiMatches := UPIRegex.FindAllString(input, -1)
	for _, upi := range upiMatches {
		if isValidUPI(upi) && !isEmail(upi) {
			intel.UPI = append(intel.UPI, strings.ToLower(upi))
		}
	}

	// Method 2: Look for @suffix patterns explicitly
	for _, suffix := range upiSuffixes {
		if idx := strings.Index(normalizedInput, suffix); idx > 0 {
			// Extract username before @
			start := idx - 1
			for start > 0 && isValidUPIChar(rune(normalizedInput[start-1])) {
				start--
			}
			if start < idx {
				upiID := normalizedInput[start : idx+len(suffix)]
				if len(upiID) > 3 && !containsString(intel.UPI, upiID) {
					intel.UPI = append(intel.UPI, upiID)
				}
			}
		}
	}

	// ============ EXTRACT PHONE NUMBERS ============
	phoneMatches := PhoneRegex.FindAllString(input, -1)
	phoneSet := make(map[string]bool)
	for _, phone := range phoneMatches {
		normalized := normalizePhone(phone)
		if len(normalized) == 10 && !phoneSet[normalized] {
			intel.Phone = append(intel.Phone, "+91"+normalized)
			phoneSet[normalized] = true
		} else if len(normalized) == 12 && strings.HasPrefix(normalized, "91") && !phoneSet[normalized[2:]] {
			intel.Phone = append(intel.Phone, "+"+normalized)
			phoneSet[normalized[2:]] = true
		}
	}

	// Method 2: Look for phone patterns with labels
	phonePatterns := []string{
		`(?i)(?:call|contact|phone|mobile|whatsapp|reach)[\s:@\-]*(\+?91)?[\s\-]?([6-9]\d{9})`,
		`(?i)(?:no|number|num)[\s:.\-]*(\+?91)?[\s\-]?([6-9]\d{9})`,
	}
	for _, pattern := range phonePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(input, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				phone := match[len(match)-1]
				normalized := normalizePhone(phone)
				if len(normalized) == 10 && !phoneSet[normalized] {
					intel.Phone = append(intel.Phone, "+91"+normalized)
					phoneSet[normalized] = true
				}
			}
		}
	}

	// ============ EXTRACT PHISHING LINKS ============
	linkMatches := PhishingLinkRegex.FindAllString(input, -1)
	for _, link := range linkMatches {
		cleanLink := cleanURL(link)
		if !isTrustedDomain(cleanLink) && len(cleanLink) > 10 {
			intel.Link = append(intel.Link, cleanLink)
		}
	}

	// Method 2: Look for suspicious URL patterns
	suspiciousPatterns := []string{
		`(?i)(?:click|visit|open|go\s*to)[\s:]*([^\s<>"']+\.[a-z]{2,}[^\s<>"']*)`,
		`(?i)link[\s:]*([^\s<>"']+\.[a-z]{2,}[^\s<>"']*)`,
	}
	for _, pattern := range suspiciousPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(input, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				link := match[1]
				if strings.Contains(link, ".") && !isTrustedDomain(link) {
					// Add http:// if missing
					if !strings.HasPrefix(strings.ToLower(link), "http") {
						link = "http://" + link
					}
					if !containsString(intel.Link, link) {
						intel.Link = append(intel.Link, link)
					}
				}
			}
		}
	}

	// ============ EXTRACT BANK ACCOUNTS ============
	// Method 1: Standard bank account regex
	bankMatches := BankAccountRex.FindAllStringSubmatch(input, -1)
	for _, match := range bankMatches {
		var digits string
		if len(match) > 1 && match[1] != "" {
			digits = match[1]
		} else {
			digits = extractDigits(match[0])
		}
		// Skip if it's a phone number or too short
		if len(digits) >= 11 && len(digits) <= 18 && !phoneSet[digits] && !isPhoneLike(digits) {
			if !containsString(intel.Bank, digits) {
				intel.Bank = append(intel.Bank, digits)
			}
		}
	}

	// Method 2: Look for account numbers with explicit labels
	accountPatterns := []string{
		`(?i)(?:account|a/?c|acct)[\s\.\-:#]*(?:no|number|num)?[\s\.\-:#]*(\d{11,18})`,
		`(?i)(?:bank|saving|current)[\s\.\-:#]*(?:a/?c|account)?[\s\.\-:#]*(\d{11,18})`,
		`(?i)(?:deposit|transfer)[\s\w]*(?:to|into)?[\s:]*(\d{11,18})`,
	}
	for _, pattern := range accountPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(input, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				digits := match[1]
				if len(digits) >= 11 && len(digits) <= 18 && !phoneSet[digits] {
					if !containsString(intel.Bank, digits) {
						intel.Bank = append(intel.Bank, digits)
					}
				}
			}
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

// MergeIntel combines two Intel structs, avoiding duplicates and limiting to 3 items per type
func MergeIntel(existing Intel, new Intel) Intel {
	const maxIntelPerType = 3

	merged := Intel{
		UPI:   limitItems(deduplicate(append(existing.UPI, new.UPI...)), maxIntelPerType),
		Phone: limitItems(deduplicate(append(existing.Phone, new.Phone...)), maxIntelPerType),
		Link:  limitItems(deduplicate(append(existing.Link, new.Link...)), maxIntelPerType),
		Bank:  limitItems(deduplicate(append(existing.Bank, new.Bank...)), maxIntelPerType),
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

// limitItems limits the slice to maxItems elements
func limitItems(items []string, maxItems int) []string {
	if len(items) > maxItems {
		return items[:maxItems]
	}
	return items
}

// ============ HELPER FUNCTIONS ============

// normalizePhone removes all non-digit characters from phone number
func normalizePhone(phone string) string {
	var result strings.Builder
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// isValidUPI checks if a string is a valid UPI ID
func isValidUPI(s string) bool {
	s = strings.ToLower(s)
	// Must contain @
	if !strings.Contains(s, "@") {
		return false
	}
	// Check for known UPI suffixes
	for _, suffix := range upiSuffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	// Also accept generic patterns
	if strings.Contains(s, "@upi") || strings.Contains(s, "@bank") {
		return true
	}
	return false
}

// isValidUPIChar checks if a character is valid in a UPI username
func isValidUPIChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') || c == '.' || c == '-' || c == '_'
}

// isEmail checks if a string looks like an email (not UPI)
func isEmail(s string) bool {
	s = strings.ToLower(s)
	emailDomains := []string{"gmail.com", "yahoo.com", "hotmail.com", "outlook.com",
		"rediffmail.com", "mail.com", "proton.me", "icloud.com", "live.com"}
	for _, domain := range emailDomains {
		if strings.HasSuffix(s, "@"+domain) || strings.Contains(s, "@"+domain) {
			return true
		}
	}
	return false
}

// isTrustedDomain checks if a URL belongs to a trusted domain
func isTrustedDomain(link string) bool {
	linkLower := strings.ToLower(link)
	for _, domain := range trustedDomains {
		if strings.Contains(linkLower, domain) {
			return true
		}
	}
	return false
}

// cleanURL removes trailing punctuation and quotes from URL
func cleanURL(url string) string {
	// Remove trailing punctuation
	url = strings.TrimRight(url, ".,;:!?\"')}]>")
	return url
}

// containsString checks if a slice contains a string
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, s) {
			return true
		}
	}
	return false
}

// isPhoneLike checks if a number looks like a phone number
func isPhoneLike(digits string) bool {
	// Indian phone numbers start with 6-9 and are 10 digits
	if len(digits) == 10 && (digits[0] >= '6' && digits[0] <= '9') {
		return true
	}
	// With country code
	if len(digits) == 12 && strings.HasPrefix(digits, "91") {
		return digits[2] >= '6' && digits[2] <= '9'
	}
	return false
}
