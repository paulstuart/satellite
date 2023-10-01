package providers

import (
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
func IdentifyAlibaba() string {
	return fileContains(AlibabaFile, csp.Alibaba, AlibabaContents)
}

// IdentifyAlibabaViaMetadataServer tries to identify Alibaba via metadata server
func IdentifyAlibabaViaMetadataServer() string {
	resp, err := httpGet(AlibabaURL, nil)
	if err != nil {
		Logger.Printf("something happened during the get %v", err)
		return ""
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			Logger.Printf("something happened during parsing the response body %v", err)
			return ""
		}
		if strings.HasPrefix(string(body), "ecs.") {
			return csp.Alibaba
		}
	}
	Logger.Printf("something happened during the request with status %s", resp.Status)
	return ""
}
