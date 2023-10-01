package providers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/paulstuart/satellite/csp"
)

//Used docs
//https://www.alibabacloud.com/help/faq-detail/49122.htm

const (
	AlibabaURL      = "http://100.100.100.200/latest/meta-data/instance/instance-type"
	AlibabaFile     = "/sys/class/dmi/id/product_name"
	AlibabaContents = "Alibaba Cloud"
)

// Identify tries to identify Alibaba provider by reading the /sys/class/dmi/id/product_name file
func IdentifyAlibaba() (string, error) {
	if fileContains(AlibabaFile, AlibabaContents) {
		return csp.Alibaba, nil
	}
	return "", nil
}

// IdentifyAlibabaViaMetadataServer tries to identify Alibaba via metadata server
func IdentifyAlibabaViaMetadataServer() (string, error) {
	resp, err := httpGet(AlibabaURL, nil)
	// req, err := http.NewRequest("GET", AlibabaURL, nil)
	// if err != nil {
	// 	return "", fmt.Errorf("could not create proper http request %w", err)
	// }
	// resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("something happened during the request %w", err)
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("something happened during parsing the response body %w", err)
		}
		if strings.HasPrefix(string(body), "ecs.") {
			return csp.Alibaba, nil
		}
	}
	return "", fmt.Errorf("something happened during the request with status %s", resp.Status)
}
