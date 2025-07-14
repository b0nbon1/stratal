package tasks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTPRequestTask performs HTTP requests
func HTTPRequestTask(ctx context.Context, params map[string]string) (string, error) {
	// Required parameters
	url, exists := params["url"]
	if !exists || url == "" {
		return "", fmt.Errorf("missing required parameter: url")
	}

	// Optional parameters with defaults
	method := strings.ToUpper(params["method"])
	if method == "" {
		method = "GET"
	}

	contentType := params["content_type"]
	if contentType == "" {
		contentType = "application/json"
	}

	body := params["body"]
	timeout := params["timeout"]

	// Parse timeout
	var timeoutDuration time.Duration = 30 * time.Second
	if timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err == nil {
			timeoutDuration = parsed
		}
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeoutDuration,
	}

	// Create request
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if contentType != "" && bodyReader != nil {
		req.Header.Set("Content-Type", contentType)
	}

	// Add custom headers (header_* parameters)
	for key, value := range params {
		if strings.HasPrefix(key, "header_") {
			headerName := strings.TrimPrefix(key, "header_")
			req.Header.Set(headerName, value)
		}
	}

	// Perform request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Format output
	output := fmt.Sprintf("Status: %s\n", resp.Status)
	output += fmt.Sprintf("Status Code: %d\n", resp.StatusCode)
	output += "Headers:\n"
	for key, values := range resp.Header {
		output += fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", "))
	}
	output += fmt.Sprintf("\nBody:\n%s", string(respBody))

	// Check if we should fail on non-2xx status
	failOnError := params["fail_on_error"] == "true"
	if failOnError && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		return output, fmt.Errorf("HTTP request returned non-success status: %d", resp.StatusCode)
	}

	return output, nil
}
