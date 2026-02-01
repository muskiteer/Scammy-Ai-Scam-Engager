package internal

import (
	"regexp"
)
var (
	 UPIRegex = regexp.MustCompile(`\b[a-zA-Z0-9.\-_]{2,}@[a-zA-Z]{2,}\b`)
	 PhoneRegex = regexp.MustCompile(`\b(\+91[\s-]?)?[6-9]\d{9}\b`)
	 PhishingLinkRegex = regexp.MustCompile(`(?i)\bhttps?:\/\/(?!.*(google|bank|gov|nic\.in))[^\s]+`)
	 BankAccountRex = regexp.MustCompile(`(?i)\b(account\s*(number|no\.?|#)?[:\s-]*\d{10,18})\b`)
)
func ExtractIntel(input string, confidence int) Intel {
	intel := Intel{}

	if UPIRegex.FindString(input) != "" {
		intel.UPI = append(intel.UPI, UPIRegex.FindString(input))
	}

	if PhoneRegex.FindString(input) != "" {
		intel.Phone = append(intel.Phone, PhoneRegex.FindString(input))
	}

	if PhishingLinkRegex.FindString(input) != "" {
		intel.Link = append(intel.Link, PhishingLinkRegex.FindString(input))
	}

	if BankAccountRex.FindString(input) != "" {
		intel.Bank = append(intel.Bank, BankAccountRex.FindString(input))
	}

	return intel
}