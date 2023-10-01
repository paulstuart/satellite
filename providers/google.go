package providers

import (
	"net/http"

	"github.com/paulstuart/satellite/csp"
)

//Used Doc
//https://cloud.google.com/compute/docs/storing-retrieving-metadata#endpoints

const (
	GoogleURL      = "http://metadata.google.internal/computeMetadata/v1/instance/tags"
	GoogleFile     = "/sys/class/dmi/id/product_name"
	GoogleContents = "Google"
)

// Identify tries to identify Google provider by reading the /sys/class/dmi/id/product_name file
func IdentifyGoogle() string {
	return fileContains(GoogleFile, csp.Google, GoogleContents)
}

// IdentifyGoogleViaMetadataServer tries to identify Google via metadata server
func IdentifyGoogleViaMetadataServer() string {
	headers := map[string]string{"Metadata-Flavor": "Google"}
	resp, err := httpGet(GoogleURL, headers)

	if err != nil {
		Logger.Printf("something happened during the request %v", err)
		return ""
	}
	if resp.StatusCode == http.StatusOK {
		return csp.Google
	}
	Logger.Printf("something happened during the request with status %s", resp.Status)
	return ""
}
