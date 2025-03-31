package transform

import (
	"context"
)

type ExecutionContext struct {
	Ctx context.Context

	OnError string   // "ignore", "fail", "error"
	Errors  []string // Collected if OnError = "error"
}

func NewExecutionContext(ctx context.Context, onError string) *ExecutionContext {
	if onError != "ignore" && onError != "fail" && onError != "error" {
		onError = "ignore"
	}

	return &ExecutionContext{
		Ctx:     ctx,
		OnError: onError,
		Errors:  []string{},
	}
}

func (ctx *ExecutionContext) HandleError(err error) bool {
	if err == nil {
		return true // no-op
	}

	switch ctx.OnError {
	case "fail":
		return false
	case "error":
		ctx.Errors = append(ctx.Errors, err.Error())
		return true
	case "ignore":
		return true
	default:
		return true
	}
}
