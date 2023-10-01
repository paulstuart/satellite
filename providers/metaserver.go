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
	dialTimeout = time.Second * 2
)

func XhttpReq(url string) (*http.Request, error) {
	ctx, _ := context.WithTimeout(context.Background(), dialTimeout)

	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func httpGet(url string, headers map[string]string) (*http.Response, error) {
	ctx, cncl := context.WithTimeout(context.Background(), dialTimeout)
	defer cncl()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return http.DefaultClient.Do(req)

	// transport := http.Transport{
	// 	Dial: dialTimeout,
	// }

	// client := http.Client{
	// 	Transport: &transport,
	// }

	// return client.Get(url)
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not create a proper http request %s", err.Error())
	// }
	// return http.DefaultClient.Do(req)
}

func IdentifyViaMetadataServer(url string, c CSP) (string, error) {
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return "", fmt.Errorf("could not create a proper http request %s", err.Error())
	// }
	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return "", fmt.Errorf("Something happened during the request %w", err)
	// }
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
