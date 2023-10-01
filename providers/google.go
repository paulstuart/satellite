package providers

import (
	"fmt"
	"net/http"

	"github.com/paulstuart/satellite/csp"
)

//Used Doc
//https://cloud.google.com/compute/docs/storing-retrieving-metadata#endpoints

// Identify tries to identify Google provider by reading the /sys/class/dmi/id/product_name file
func IdentifyGoogle() (string, error) {
	if fileContains(GoogleFile, GoogleContents) {
		return csp.Google, nil
	}
	return "", nil
}

const (
	GoogleURL      = "http://metadata.google.internal/computeMetadata/v1/instance/tags"
	GoogleFile     = "/sys/class/dmi/id/product_name"
	GoogleContents = "Google"
)

// IdentifyGoogleViaMetadataServer tries to identify Google via metadata server
func IdentifyGoogleViaMetadataServer() (string, error) {
	headers := map[string]string{"Metadata-Flavor": "Google"}
	resp, err := httpGet(GoogleURL, headers)

	if err != nil {
		return "", fmt.Errorf("something happened during the request %w", err)
	}
	if resp.StatusCode == http.StatusOK {
		return csp.Google, nil
	}
	return "", fmt.Errorf("something happened during the request with status %s", resp.Status)
}
