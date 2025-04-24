package middleware

import (
	"context"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor injects user_id from metadata into context.
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.Printf("[info] AuthInterceptor: metadata=%v", md)
			if userIDVals := md.Get("x-user-id"); len(userIDVals) > 0 {
				if uid, err := strconv.Atoi(userIDVals[0]); err == nil {
					ctx = context.WithValue(ctx, "user_id", int32(uid))
				}
			}
		}
		return handler(ctx, req)
	}
}
