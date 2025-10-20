package interceptors

import (
	"context"
	"log"
	"time"

	"connectrpc.com/connect"
)

// LoggingInterceptor logs all RPC calls with method name and duration.
func LoggingInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			method := req.Spec().Procedure

			log.Printf("RPC: %s started", method)

			resp, err := next(ctx, req)

			duration := time.Since(start)
			if err != nil {
				log.Printf("RPC: %s failed in %v: %v", method, duration, err)
			} else {
				log.Printf("RPC: %s succeeded in %v", method, duration)
			}

			return resp, err
		}
	}
}
