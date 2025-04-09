package transform

import (
	"context"
)

type contextKey string

const ContextKeyDebug contextKey = "ducto-debug"

type ExecutionContext struct {
	Ctx     context.Context
	Dropped bool // Don't send to the output writer

	OnError string   // "ignore", "fail", "capture"
	Errors  []string // Collected if OnError = "capture"

	// Debug flag
	Debug bool

	// Future: you could add TraceID, Logger, Metrics, etc.
}

func NewExecutionContext(ctx context.Context, onError string) *ExecutionContext {
	if onError != "ignore" && onError != "fail" && onError != "capture" {
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
		return true
	}

	switch ctx.OnError {
	case "fail":
		return false
	case "capture":
		ctx.Errors = append(ctx.Errors, err.Error())
		return true
	case "ignore":
		return true
	default:
		return true
	}
}
