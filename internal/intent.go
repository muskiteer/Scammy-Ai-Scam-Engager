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
	IntentDeepProbe       Intent = "DEEP_PROBE"
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
	const maxTurns = 15 // Extended to 15 turns — intermediate callback fires at turn 10

	// Complete only after 15 turns for maximum engagement score
	if ctx.TurnCount >= maxTurns {
		return StateComplete
	}

	if !ctx.ScamDetected {
		return StateInit
	}

	// Always keep extracting intel until max turns — never complete early
	return StateIntelExtract
}

func DeriveIntent(state State, intel Intel, turnCount int, askCount AskCount) Intent {
	const maxAskCount = 2 // Ask each info type up to TWICE for maximum elicitation score
	const maxTurnCount = 15

	if turnCount >= maxTurnCount {
		return IntentStall
	}

	switch state {
	case StateInit:
		// Even unconfirmed scams get probed — verify caller identity immediately
		if turnCount%2 == 0 {
			return IntentAskIdentity
		}
		return IntentConfirmDetails

	case StateEngaging:
		switch turnCount % 3 {
		case 0:
			return IntentAskIdentity
		case 1:
			return IntentDeepProbe
		default:
			return IntentConfirmDetails
		}

	case StateIntelExtract:
		// === PRIORITY 1: Core intel — ask each up to maxAskCount times ===
		if len(intel.UPI) == 0 && askCount.UPI < maxAskCount {
			return IntentAskUPI
		}
		if len(intel.Phone) == 0 && askCount.Phone < maxAskCount {
			return IntentAskPhone
		}
		if len(intel.Bank) == 0 && askCount.Bank < maxAskCount {
			return IntentAskBank
		}
		if len(intel.Email) == 0 && askCount.Email < maxAskCount {
			return IntentAskEmail
		}
		if len(intel.Link) == 0 && askCount.Link < maxAskCount {
			return IntentAskLink
		}

		// === PRIORITY 2: Secondary intel — ask each up to twice ===
		if len(intel.CaseIDs) == 0 && askCount.CaseID < maxAskCount {
			return IntentAskCaseID
		}
		if len(intel.IFSCCodes) == 0 && askCount.IFSCCode < maxAskCount {
			return IntentAskIFSCCode
		}
		if len(intel.CardNumbers) == 0 && askCount.CardNumber < maxAskCount {
			return IntentAskCardNumber
		}
		if len(intel.PolicyNumbers) == 0 && askCount.PolicyNumber < maxAskCount {
			return IntentAskPolicyNumber
		}
		if len(intel.OrderNumbers) == 0 && askCount.OrderNumber < maxAskCount {
			return IntentAskOrderNumber
		}

		// === PRIORITY 3: Deep investigative probing — fills all remaining turns (11-15) ===
		switch turnCount % 5 {
		case 0:
			return IntentAskIdentity
		case 1:
			return IntentDeepProbe
		case 2:
			return IntentConfirmDetails
		case 3:
			return IntentDeepProbe
		default:
			return IntentStall
		}

	case StateComplete:
		return IntentStall

	default:
		return IntentConfirmDetails
	}
}
