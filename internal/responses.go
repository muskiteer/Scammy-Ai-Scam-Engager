package internal

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var responses = map[Intent][]string{
	IntentConfirmDetails: {
		"Why will my account be blocked?",
		"Can you explain what the issue is?",
		"Which account is this regarding?",
		"What exactly is the problem with my account?",
		"I don't understand, can you clarify?",
	},

	IntentAskUPI: {
		"Can you share the UPI ID again?",
		"Which UPI ID should I use?",
		"Please send the UPI ID once more.",
		"What is the correct UPI ID?",
		"I didn't catch the UPI ID, can you repeat?",
	},

	IntentAskLink: {
		"Can you send the link here?",
		"I did not receive the link properly.",
		"Please share the link again.",
		"Where is the link?",
		"Can you resend the verification link?",
	},

	IntentAskPhone: {
		"Which number should I contact?",
		"Can you share the phone number?",
		"Is there a number I can call?",
		"What is the customer care number?",
		"Please provide the contact number.",
	},

	IntentStall: {
		"I am checking this now, please wait.",
		"Give me a moment to look into this.",
		"I will check and get back shortly.",
		"Let me verify this information.",
		"One moment please, I'm looking into it.",
		"I need to check with my bank first.",
	},

	IntentNeutral: {
		"Okay.",
		"I see.",
		"Understood.",
		"Got it.",
		"Alright.",
	},
}

// GetResponse returns a random response for the given intent
func GetResponse(intent Intent) string {
	templates, exists := responses[intent]
	if !exists || len(templates) == 0 {
		return "I see."
	}

	index := rng.Intn(len(templates))
	return templates[index]
}
