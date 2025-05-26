package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		userIDValues := md.Get("x-user-id")
		if len(userIDValues) == 0 {
			log.Println("AuthUnaryServerInterceptor: x-user-id not found in metadata")
			return nil, status.Errorf(codes.Unauthenticated, "x-user-id is not provided")
		}

		var userID string
		if len(userIDValues) > 0 {
			userID = userIDValues[0]
			ctx = context.WithValue(ctx, UserIDKey, userID)
			log.Printf("AuthUnaryServerInterceptor: UserID %s extracted and added to context", userID)
		}

		return handler(ctx, req)
	}
}
