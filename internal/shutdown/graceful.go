package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GracefulShutdown struct {
	timeout time.Duration
	cleanup []func(context.Context) error
}

func New(timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		timeout: timeout,
	}
}

// Register accepts BOTH:
//
//	func() error
//	func(context.Context) error
func (gs *GracefulShutdown) Register(fn any) {
	switch f := fn.(type) {

	case func(context.Context) error:
		gs.cleanup = append(gs.cleanup, f)

	case func() error:
		wrapped := func(ctx context.Context) error {
			return f()
		}
		gs.cleanup = append(gs.cleanup, wrapped)

	default:
		panic("Register: invalid function type (must be func() error or func(context.Context) error)")
	}
}

func (gs *GracefulShutdown) Wait(ctx context.Context) (int, error) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan

	// Convert to UNIX exit code
	exitCode := 128 + int(sig.(syscall.Signal))

	// cleanup with timeout
	cleanupCtx, cancel := context.WithTimeout(ctx, gs.timeout)
	defer cancel()

	for _, fn := range gs.cleanup {
		if err := fn(cleanupCtx); err != nil {
			return exitCode, err
		}
	}
	return exitCode, nil
}
