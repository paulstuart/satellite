package providers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/paulstuart/satellite/csp"
)

//Used docs
// https://azure.microsoft.com/en-us/blog/announcing-general-availability-of-azure-instance-metadata-service/

const (
	AzureURL  = "http://169.254.169.254/metadata/instance?api-version=2017-12-01"
	AzureFile = "/sys/class/dmi/id/sys_vendor"
)

// Identify tries to identify Azure provider by reading the /sys/class/dmi/id/sys_vendor file
func IdentifyAsure() (string, error) {
	data, err := os.ReadFile(AzureFile)
	if err != nil {
		return "", fmt.Errorf("something happened during reading a file: %w", err)
	}
	if strings.Contains(string(data), "Microsoft Corporation") {
		return csp.Azure, nil
	}
	return "", nil
}

// IdentifyAzureViaMetadataServer tries to identify Azure via metadata server
func IdentifyAzureViaMetadataServer() (string, error) {
	headers := map[string]string{"Metadata": "true"}
	resp, err := httpGet(AzureURL, headers)
	if err != nil {
		return "", fmt.Errorf("something happened during the request %w", err)
	}
	if resp.StatusCode == http.StatusOK {
		return csp.Azure, nil
	}
	return "", fmt.Errorf("Something happened during the request with status %s", resp.Status)
}
