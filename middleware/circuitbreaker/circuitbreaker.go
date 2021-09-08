package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

type (
	// Option is circuit breaker option.
	Option func(*options)
	// NewFunc returns a new breaker
	NewFunc func() circuitbreaker.CircuitBreaker
)

var failedErr = errors.New(503, "CIRCUITBREAKER", "request failed due to circuit breaker triggered")

// WithNewBeaker set the New function of circuit breaker implentation
func WithNewFunc(newFunc NewFunc) Option {
	return func(o *options) {
		o.newFunc = newFunc
	}
}

// WithName set the breaker name for target
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

type options struct {
	name    string
	newFunc NewFunc
}

// Client circuitbreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func Client(opts ...Option) middleware.Middleware {
	options := &options{
		newFunc: func() circuitbreaker.CircuitBreaker { return sre.NewBreaker() },
	}

	for _, o := range opts {
		o(options)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// TODO
			return reply, nil
		}
	}
}
