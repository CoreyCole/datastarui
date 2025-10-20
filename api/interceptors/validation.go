package interceptors

import (
	"context"
	"errors"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

// ValidationInterceptor returns a Connect interceptor that validates all
// incoming requests using proto-validate rules defined in proto files.
// Invalid requests are rejected with CodeInvalidArgument and detailed error messages.
func ValidationInterceptor() connect.UnaryInterceptorFunc {
	validator, err := protovalidate.New()
	if err != nil {
		panic("failed to initialize validator: " + err.Error())
	}

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Validate the request message
			if err := validator.Validate(req.Any().(proto.Message)); err != nil {
				connectErr := connect.NewError(connect.CodeInvalidArgument, err)

				// Add detailed validation violations to error
				if validationErr := new(protovalidate.ValidationError); errors.As(err, &validationErr) {
					if detail, err := connect.NewErrorDetail(validationErr.ToProto()); err == nil {
						connectErr.AddDetail(detail)
					}
				}

				return nil, connectErr
			}

			return next(ctx, req)
		}
	}
}
