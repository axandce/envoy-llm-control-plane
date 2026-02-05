package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
    "flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
    "time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/grpc/reflection"
)

type server struct {
	authv3.UnimplementedAuthorizationServer
	requiredAPIKey string
	mode           string // allow | deny_missing_key | deny_quota
}

func main() {
    addr := env("LISTEN_ADDR", "0.0.0.0:9002")
    required := env("REQUIRED_API_KEY", "demo-key")
    mode := env("POLICY_MODE", "allow")

    // Add a distroless-friendly healthcheck mode.
    healthcheck := flag.Bool("healthcheck", false, "run healthcheck and exit")
    flag.Parse()

    if *healthcheck {
        if err := tcpHealthcheck(addr, 500*time.Millisecond); err != nil {
            log.Printf("healthcheck failed: %v", err)
            os.Exit(1)
        }
        fmt.Print("ok")
        return
    }

    lis, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("listen: %v", err)
    }

    s := grpc.NewServer()
    authv3.RegisterAuthorizationServer(s, &server{
        requiredAPIKey: required,
        mode:           mode,
    })
    reflection.Register(s)

    log.Printf("policy-engine-stub listening on %s (mode=%s)", addr, mode)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("serve: %v", err)
    }
}

func tcpHealthcheck(listenAddr string, timeout time.Duration) error {
    _, port, err := net.SplitHostPort(listenAddr)
    if err != nil {
        return fmt.Errorf("invalid LISTEN_ADDR %q: %w", listenAddr, err)
    }
    conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", port), timeout)
    if err != nil {
        return err
    }
    _ = conn.Close()
    return nil
}

func (s *server) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	headers := req.GetAttributes().GetRequest().GetHttp().GetHeaders()

	apiKey := firstHeader(headers, "x-api-key")
	if s.mode == "deny_quota" {
		return deny(429, codes.ResourceExhausted, "quota_exceeded", "Quota exceeded"), nil
	}
	if s.mode == "deny_missing_key" || strings.TrimSpace(apiKey) == "" {
		return deny(401, codes.Unauthenticated, "missing_api_key", "Missing API key"), nil
	}
	if apiKey != s.requiredAPIKey {
		return deny(403, codes.PermissionDenied, "invalid_api_key", "Invalid API key"), nil
	}

	decisionID := randID(8)

	// Allow and inject a couple of headers so downstream can correlate/audit.
	ok := &authv3.OkHttpResponse{
		Headers: []*corev3.HeaderValueOption{
			hdr("x-llm-decision-id", decisionID),
			hdr("x-llm-policy-version", "stub-v1"),
		},
	}

	return &authv3.CheckResponse{
		Status: grpcstatus.New(codes.OK, "ok").Proto(),
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: ok,
		},
	}, nil
}

func deny(httpStatus int, grpcCode codes.Code, code, message string) *authv3.CheckResponse {
	body := fmt.Sprintf(`{"error":{"message":%q,"type":"policy_error","code":%q}}`, message, code)

	denied := &authv3.DeniedHttpResponse{
		Status: &typev3.HttpStatus{Code: typev3.StatusCode(httpStatus)},
		Headers: []*corev3.HeaderValueOption{
			hdr("content-type", "application/json"),
			hdr("x-llm-deny-code", code),
		},
		Body: body,
	}

	return &authv3.CheckResponse{
		Status: grpcstatus.New(grpcCode, message).Proto(),
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: denied,
		},
	}
}

func hdr(k, v string) *corev3.HeaderValueOption {
	return &corev3.HeaderValueOption{
		Header: &corev3.HeaderValue{Key: k, Value: v},
	}
}

func firstHeader(hdrs map[string]string, key string) string {
	// Envoy normalizes header keys to lowercase in ext_authz.
	return hdrs[strings.ToLower(key)]
}

func randID(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func env(k, def string) string {
	if v := os.Getenv(k); strings.TrimSpace(v) != "" {
		return v
	}
	return def
}
