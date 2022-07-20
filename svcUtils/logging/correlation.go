package logging

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	correlationIDKey = "correlation-id"
)

type (
	correlationKey struct{}
	ctxLogger      struct{}
)

func WithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	cid := CorrelationIDFromContext(ctx)
	if cid == "" {
		cid = generateCorrelationID()
		ctx = context.WithValue(ctx, correlationKey{}, cid)
	}
	return context.WithValue(ctx, ctxLogger{}, logger.WithField(correlationIDKey, cid))
}

func FromContext(ctx context.Context) *logrus.Entry {
	if lgr := ctx.Value(ctxLogger{}); lgr != nil {
		return lgr.(*logrus.Entry)
	}
	return NewLogger().WithField(correlationIDKey, CorrelationIDFromContext(ctx))
}

func WithContext(ctx context.Context, logger logrus.FieldLogger) *logrus.Entry {
	cid := CorrelationIDFromContext(ctx)
	return logger.WithField(correlationIDKey, cid)
}

func CorrelationIDFromContext(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		correlationIDs := md[correlationIDKey]
		if len(correlationIDs) > 0 {
			return correlationIDs[0]
		}
	}
	if cid := ctx.Value(correlationKey{}); cid != nil {
		return cid.(string)
	}
	return ""
}

func ContextWithCorrelationID(ctx context.Context) context.Context {
	cid := CorrelationIDFromContext(ctx)
	if cid == "" {
		cid := generateCorrelationID()
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md[correlationIDKey] = []string{cid}
		} else {
			md = metadata.Pairs(correlationIDKey, cid)
		}
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, correlationIDKey, cid)
	return ctx
}

const logLead = "x-log-"

func UnaryClientInterceptor(ignore ...string) grpc.UnaryClientInterceptor {
	ig := make(map[string]bool, len(ignore))
	for _, v := range ignore {
		ig[v] = true
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		lg := FromContext(ctx)
		md := metautils.ExtractOutgoing(ctx)
		for k, v := range lg.Data {
			if ig[k] {
				continue
			}
			md.Add(logLead+k, fmt.Sprint(v))
		}
		if md.Get(correlationIDKey) == "" {
			md.Set(correlationIDKey, CorrelationIDFromContext(ctx))
		}
		ctx = md.ToOutgoing(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func UnaryServerInterceptor(log *logrus.Entry, clear bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = ContextWithCorrelationID(ctx)
		ctx = WithLogger(ctx, log)
		ctx = forwardLogger(ctx, clear)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(log *logrus.Entry, clear bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx := ContextWithCorrelationID(ss.Context())
		newCtx = WithLogger(newCtx, log)
		newCtx = forwardLogger(newCtx, clear)
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}

func forwardLogger(ctx context.Context, clear bool) context.Context {
	md := metautils.ExtractIncoming(ctx)
	if md.Get(correlationIDKey) == "" {
		md.Set(correlationIDKey, CorrelationIDFromContext(ctx))
	}
	lg := FromContext(ctx)
	for k, v := range md {
		if strings.HasPrefix(k, logLead) {
			md.Del(k)
			if len(v) > 0 && !clear {
				lg = lg.WithField(strings.TrimPrefix(k, logLead), fmt.Sprint(v[0]))
			}
		}
	}
	grpc_ctxtags.Extract(ctx).Set(correlationIDKey, md.Get(correlationIDKey))
	ctx = md.ToIncoming(ctx)
	lg = WithContext(ctx, lg)
	return WithLogger(ctx, lg)
}

func generateCorrelationID() string {
	return uuid.New().String()
}

func LoggerMiddleware(logger logrus.FieldLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(WithLogger(r.Context(), logger)))
		})
	}
}

func LoggerMiddlewareFunc(logger logrus.FieldLogger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r.WithContext(WithLogger(r.Context(), logger)))
		}
	}
}
