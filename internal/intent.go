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
	IntentStall          Intent = "STALL"
	IntentNeutral        Intent = "NEUTRAL"
)

type Intel struct {
	UPI   []string
	Phone []string
	Link  []string
	Bank  []string
}

type SessionContext struct {
	ScamDetected bool
	TurnCount    int
	Intel        Intel
	CurrentState State
}

func GetState(ctx SessionContext) State {
	if !ctx.ScamDetected {
		return StateInit
	}

	// Check if we have complete intelligence
	hasAllIntel := len(ctx.Intel.UPI) > 0 && len(ctx.Intel.Phone) > 0 &&
		len(ctx.Intel.Link) > 0 && len(ctx.Intel.Bank) > 0

	if hasAllIntel {
		return StateComplete
	}

	// Check if we have partial intelligence - start extracting
	hasPartialIntel := len(ctx.Intel.UPI) > 0 || len(ctx.Intel.Phone) > 0 ||
		len(ctx.Intel.Link) > 0 || len(ctx.Intel.Bank) > 0

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
func DeriveIntent(state State, intel Intel, turnCount int) Intent {
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
		// Prioritize intelligence gathering
		if len(intel.UPI) == 0 {
			return IntentAskUPI
		}
		if len(intel.Phone) == 0 {
			return IntentAskPhone
		}
		if len(intel.Link) == 0 {
			return IntentAskLink
		}
		if len(intel.Bank) == 0 {
			return IntentConfirmDetails
		}
		return IntentStall

	case StateComplete:
		return IntentStall

	default:
		return IntentNeutral
	}
}
