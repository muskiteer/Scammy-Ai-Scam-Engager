package internal

type State string

const (
	StateInit          State = "INIT"
	StateEngaging      State = "ENGAGING"
	StateIntelExtract  State = "INTEL_EXTRACT"
	StateComplete      State = "COMPLETE"
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
}

func GetState(ctx SessionContext) State {

	
	if !ctx.ScamDetected {
		return StateInit
	}

	
	if ctx.TurnCount < 3 {
		return StateEngaging
	}

	
	if len(ctx.Intel.UPI) == 0 || len(ctx.Intel.Phone) == 0 || len(ctx.Intel.Link) == 0 || len(ctx.Intel.Bank) == 0 {
		return StateIntelExtract
	}

	
	return StateComplete
}