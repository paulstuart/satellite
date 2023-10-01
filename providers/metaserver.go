package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CSP interface {
	IsCSP() string
}

var (
	Timeout = time.Second * 2
	logger  = log.Default()
)

func httpGet(url string, headers map[string]string) (*http.Response, error) {
	ctx, cncl := context.WithTimeout(context.Background(), Timeout)
	defer cncl()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return http.DefaultClient.Do(req)
}

func IdentifyViaMetadataServer(url string, c CSP) (string, error) {
	resp, err := httpGet(url, nil)
	if err != nil {
		return "", fmt.Errorf("get failed for %q: %w", url, err)
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed parsing the response body %w", err)
		}
		err = json.Unmarshal(body, c)
		if err != nil {
			return "", fmt.Errorf("failed unmarshalling the response body %w", err)
		}
		return c.IsCSP(), nil
	}
	return "", fmt.Errorf("something happened during the request with status %s", resp.Status)
}

func fileContains(file string, what ...string) bool {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Printf("failed to read file %q -- %v", file, err)
		return false
	}
	data := string(b)
	for _, item := range what {
		if strings.Contains(data, item) {
			return true
		}
	}
	return false
}
