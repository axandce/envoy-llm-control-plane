package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type policyMeta struct {
	DecisionID    string `json:"x-llm-decision-id,omitempty"`
	PolicyVersion string `json:"x-llm-policy-version,omitempty"`
}

type chatReq struct {
	Model    string `json:"model"`
	Stream   bool   `json:"stream"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type chatResp struct {
	ID      string     `json:"id"`
	Object  string     `json:"object"`
	Created int        `json:"created"`
	Model   string     `json:"model"`
	Policy  policyMeta `json:"policy,omitempty"`
	Choices []struct {
		Index        int `json:"index"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}


func main() {
	addr := env("LISTEN_ADDR", "0.0.0.0:9000")

	healthcheck := flag.Bool("healthcheck", false, "run healthcheck and exit")
	flag.Parse()

	if *healthcheck {
		if err := httpHealthcheck(addr, "/healthz", 500*time.Millisecond); err != nil {
			log.Printf("healthcheck failed: %v", err)
			os.Exit(1)
		}
		fmt.Print("ok")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/v1/chat/completions", handleChatCompletions)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("mock-upstream listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func httpHealthcheck(listenAddr, path string, timeout time.Duration) error {
	_, port, err := net.SplitHostPort(listenAddr)
	if err != nil {
		return fmt.Errorf("invalid LISTEN_ADDR %q: %w", listenAddr, err)
	}

	u := fmt.Sprintf("http://127.0.0.1:%s%s", port, path)

	c := &http.Client{Timeout: timeout}
	resp, err := c.Get(u)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-2xx: %s", resp.Status)
	}
	return nil
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req chatReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if req.Stream {
		streamSSE(w)
		return
	}

	// Minimal non-streaming response shape (enough for demos).
	// Echo policy metadata that Envoy injected into the upstream request.
	resp := chatResp{
		ID:      "mockcmpl-1",
		Object:  "chat.completion",
		Created: 0,
		Model:   "mock",
		Policy: policyMeta{
			DecisionID:    r.Header.Get("x-llm-decision-id"),
			PolicyVersion: r.Header.Get("x-llm-policy-version"),
		},
	}

	// Populate the single choice (matches your anonymous struct type).
	resp.Choices = make([]struct {
		Index        int `json:"index"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	}, 1)

	resp.Choices[0].Index = 0
	resp.Choices[0].Message.Role = "assistant"
	resp.Choices[0].Message.Content = "hello from mock upstream"
	resp.Choices[0].FinishReason = "stop"

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(&resp)

}

func streamSSE(w http.ResponseWriter) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	emit := func(data string) {
		_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	emit(`{"id":"mockcmpl-1","object":"chat.completion.chunk","created":0,"model":"mock","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}`)
	time.Sleep(200 * time.Millisecond)

	emit(`{"id":"mockcmpl-1","object":"chat.completion.chunk","created":0,"model":"mock","choices":[{"index":0,"delta":{"content":"hello"},"finish_reason":null}]}`)
	time.Sleep(200 * time.Millisecond)

	emit(`{"id":"mockcmpl-1","object":"chat.completion.chunk","created":0,"model":"mock","choices":[{"index":0,"delta":{"content":" from envoy demo"},"finish_reason":null}]}`)
	time.Sleep(200 * time.Millisecond)

	emit(`[DONE]`)
}

func env(k, def string) string {
	if v := os.Getenv(k); strings.TrimSpace(v) != "" {
		return v
	}
	return def
}
