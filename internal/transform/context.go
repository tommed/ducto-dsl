package transform

type ExecutionContext struct {
	OnError string  // "ignore", "fail", "error"
	Errors  []error // Collected if OnError = "error"
}

func NewExecutionContext(onError string) *ExecutionContext {
	if onError != "ignore" && onError != "fail" && onError != "error" {
		onError = "ignore"
	}
	return &ExecutionContext{OnError: onError}
}

func (ctx *ExecutionContext) Handle(err error) bool {
	switch ctx.OnError {
	case "fail":
		return false
	case "error":
		ctx.Errors = append(ctx.Errors, err)
		return true
	case "ignore":
		return true
	default:
		return true
	}
}
