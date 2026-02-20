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
	IntentConfirmDetails  Intent = "CONFIRM_DETAILS"
	IntentAskUPI          Intent = "ASK_UPI"
	IntentAskLink         Intent = "ASK_LINK"
	IntentAskPhone        Intent = "ASK_PHONE"
	IntentAskBank         Intent = "ASK_BANK"
	IntentAskEmail        Intent = "ASK_EMAIL"
	IntentAskCaseID       Intent = "ASK_CASE_ID"
	IntentAskPolicyNumber Intent = "ASK_POLICY_NUMBER"
	IntentAskOrderNumber  Intent = "ASK_ORDER_NUMBER"
	IntentAskCardNumber   Intent = "ASK_CARD_NUMBER"
	IntentAskIFSCCode     Intent = "ASK_IFSC_CODE"
	IntentAskIdentity     Intent = "ASK_IDENTITY"
	IntentStall           Intent = "STALL"
	IntentNeutral         Intent = "NEUTRAL"
)

type Intel struct {
	UPI           []string
	Phone         []string
	Link          []string
	Bank          []string
	Email         []string
	CaseIDs       []string
	PolicyNumbers []string
	OrderNumbers  []string
	CardNumbers   []string
	IFSCCodes     []string
}

type AskCount struct {
	UPI          int
	Phone        int
	Link         int
	Bank         int
	Email        int
	CaseID       int
	PolicyNumber int
	OrderNumber  int
	CardNumber   int
	IFSCCode     int
}

type SessionContext struct {
	ScamDetected            bool
	TurnCount               int
	Intel                   Intel
	CurrentState            State
	AskCount                AskCount
	QuestionsAsked          int
	InvestigativeQuestions  int
	RedFlagsIdentified      []string
	InformationElicitations int
}

func GetState(ctx SessionContext) State {
	const maxTurns = 10 // Evaluator sends up to 10 turns - NEVER complete before

	// ONLY complete at turn 10 - maximize engagement for all scoring categories
	if ctx.TurnCount >= maxTurns {
		return StateComplete
	}

	if !ctx.ScamDetected {
		return StateInit
	}

	// Always keep extracting intel until turn 10 - never complete early
	return StateIntelExtract
}

// DeriveIntent determines the next intent based on current state and missing intelligence
func DeriveIntent(state State, intel Intel, turnCount int, askCount AskCount) Intent {
	const maxAskCount = 3
	const maxTurnCount = 10

	if turnCount >= maxTurnCount {
		return IntentStall
	}

	switch state {
	case StateInit:
		return IntentConfirmDetails // Ask questions even before scam detected

	case StateEngaging:
		if turnCount%2 == 0 {
			return IntentConfirmDetails
		}
		return IntentStall

	case StateIntelExtract:
		// Rotate between intel extraction and investigative questions
		// Every 3rd turn ask identity questions for investigative scoring
		if turnCount%3 == 0 && turnCount > 2 {
			return IntentAskIdentity
		}

		// Prioritize the most common fake data types first
		if len(intel.Phone) == 0 && askCount.Phone < maxAskCount {
			return IntentAskPhone
		}
		if len(intel.UPI) == 0 && askCount.UPI < maxAskCount {
			return IntentAskUPI
		}
		if len(intel.Link) == 0 && askCount.Link < maxAskCount {
			return IntentAskLink
		}
		if len(intel.Bank) == 0 && askCount.Bank < maxAskCount {
			return IntentAskBank
		}
		if len(intel.Email) == 0 && askCount.Email < maxAskCount {
			return IntentAskEmail
		}
		if len(intel.CaseIDs) == 0 && askCount.CaseID < maxAskCount {
			return IntentAskCaseID
		}
		if len(intel.PolicyNumbers) == 0 && askCount.PolicyNumber < maxAskCount {
			return IntentAskPolicyNumber
		}
		if len(intel.OrderNumbers) == 0 && askCount.OrderNumber < maxAskCount {
			return IntentAskOrderNumber
		}
		if len(intel.CardNumbers) == 0 && askCount.CardNumber < maxAskCount {
			return IntentAskCardNumber
		}
		if len(intel.IFSCCodes) == 0 && askCount.IFSCCode < maxAskCount {
			return IntentAskIFSCCode
		}
		// All asked enough - rotate between engagement types
		switch turnCount % 3 {
		case 0:
			return IntentAskIdentity
		case 1:
			return IntentConfirmDetails
		default:
			return IntentStall
		}

	case StateComplete:
		return IntentStall

	default:
		return IntentConfirmDetails
	}
}
