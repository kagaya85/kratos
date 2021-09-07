package circuitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Option is circuit breaker option.
type Option func(*options)

// WithNewBeaker set the New function of circuit breaker implentation
func WithNewFunc(newFunc circuitbreaker.NewFunc) Option {
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
	name          string
	newFunc       circuitbreaker.NewFunc
	fallbackFuncs []circuitbreaker.FallbackFunc
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

	// create breaker group
	group := &circuitbreaker.Group{
		New: options.newFunc,
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if gerr := group.Do(options.name, func() error {
				reply, err = handler(ctx, req)
				return err
			}, options.fallbackFuncs...); gerr != nil {
				// rejected
				return reply, errors.New(503, "CIRCUITBREAKER", "request failed due to circuit breaker triggered")
			}
			return reply, err
		}
	}
}
