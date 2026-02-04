package internal

type State string

const (
	StateInit         State = "INIT"
	StateEngaging     State = "ENGAGING"
	StateIntelExtract State = "INTEL_EXTRACT"
	StateComplete     State = "COMPLETE"
)

type Intent string

const (
	IntentConfirmDetails Intent = "CONFIRM_DETAILS"
	IntentAskUPI         Intent = "ASK_UPI"
	IntentAskLink        Intent = "ASK_LINK"
	IntentAskPhone       Intent = "ASK_PHONE"
	IntentAskBank        Intent = "ASK_BANK"
	IntentStall          Intent = "STALL"
	IntentNeutral        Intent = "NEUTRAL"
)

type Intel struct {
	UPI   []string
	Phone []string
	Link  []string
	Bank  []string
}

type AskCount struct {
	UPI   int
	Phone int
	Link  int
	Bank  int
}

type SessionContext struct {
	ScamDetected bool
	TurnCount    int
	Intel        Intel
	CurrentState State
	AskCount     AskCount // Track how many times we've asked for each intel type
}

func GetState(ctx SessionContext) State {
	const maxAskCount = 3
	const maxTurns = 15

	// Force completion if we've reached max turns
	if ctx.TurnCount >= maxTurns {
		return StateComplete
	}

	if !ctx.ScamDetected {
		return StateInit
	}

	// Check if we have complete intelligence
	hasAllIntel := len(ctx.Intel.UPI) > 0 && len(ctx.Intel.Phone) > 0 &&
		len(ctx.Intel.Link) > 0 && len(ctx.Intel.Bank) > 0

	if hasAllIntel {
		return StateComplete
	}

	// Check if we have sufficient intelligence OR exhausted all attempts
	intelCount := len(ctx.Intel.UPI) + len(ctx.Intel.Phone) +
		len(ctx.Intel.Link) + len(ctx.Intel.Bank)

	// End if we have at least 2 pieces of intel and exhausted attempts for missing ones
	if intelCount >= 2 {
		allMissingExhausted := true
		if len(ctx.Intel.UPI) == 0 && ctx.AskCount.UPI < maxAskCount {
			allMissingExhausted = false
		}
		if len(ctx.Intel.Phone) == 0 && ctx.AskCount.Phone < maxAskCount {
			allMissingExhausted = false
		}
		if len(ctx.Intel.Link) == 0 && ctx.AskCount.Link < maxAskCount {
			allMissingExhausted = false
		}
		if len(ctx.Intel.Bank) == 0 && ctx.AskCount.Bank < maxAskCount {
			allMissingExhausted = false
		}

		if allMissingExhausted {
			return StateComplete
		}
	}

	// Check if we have partial intelligence - start extracting
	hasPartialIntel := intelCount > 0

	if hasPartialIntel {
		return StateIntelExtract
	}

	// Still in early engagement phase
	if ctx.TurnCount < 2 {
		return StateEngaging
	}

	// After 2 turns with no intel, actively extract
	return StateIntelExtract
}

// DeriveIntent determines the next intent based on current state and missing intelligence
func DeriveIntent(state State, intel Intel, turnCount int, askCount AskCount) Intent {
	const maxAskCount = 3
	const maxTurnCount = 15

	// End conversation if max turns reached
	if turnCount >= maxTurnCount {
		return IntentStall
	}

	switch state {
	case StateInit:
		return IntentNeutral

	case StateEngaging:
		// Early engagement - ask for clarification
		if turnCount%2 == 0 {
			return IntentConfirmDetails
		}
		return IntentStall

	case StateIntelExtract:
		// Prioritize intelligence gathering, but only if we haven't asked 3 times already
		if len(intel.UPI) == 0 && askCount.UPI < maxAskCount {
			return IntentAskUPI
		}
		if len(intel.Phone) == 0 && askCount.Phone < maxAskCount {
			return IntentAskPhone
		}
		if len(intel.Link) == 0 && askCount.Link < maxAskCount {
			return IntentAskLink
		}
		if len(intel.Bank) == 0 && askCount.Bank < maxAskCount {
			return IntentAskBank
		}
		return IntentStall

	case StateComplete:
		return IntentStall

	default:
		return IntentNeutral
	}
}
