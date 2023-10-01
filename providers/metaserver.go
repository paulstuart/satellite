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
	Logger  = log.Default()
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

func IdentifyViaMetadataServer(url string, c CSP) string {
	resp, err := httpGet(url, nil)
	if err != nil {
		Logger.Printf("get failed for %q: %v", url, err)
		return ""
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			Logger.Printf("failed parsing the response body %v", err)
			return ""
		}
		err = json.Unmarshal(body, c)
		if err != nil {
			Logger.Printf("failed unmarshalling the response body %v", err)
			return ""
		}
		return c.IsCSP()
	}
	Logger.Printf("something happened during the request with status %s", resp.Status)
	return ""
}

func fileContains(file, csp string, what ...string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		Logger.Printf("failed to read file for %s %q -- %v", csp, file, err)
		return ""
	}
	data := string(b)
	for _, item := range what {
		if strings.Contains(data, item) {
			return csp
		}
	}
	return ""
}
