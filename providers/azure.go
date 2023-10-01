package providers

import (
	"net/http"

	"github.com/paulstuart/satellite/csp"
)

//Used docs
// https://azure.microsoft.com/en-us/blog/announcing-general-availability-of-azure-instance-metadata-service/

const (
	AzureURL      = "http://169.254.169.254/metadata/instance?api-version=2017-12-01"
	AzureFile     = "/sys/class/dmi/id/sys_vendor"
	AzureContents = "Microsoft Corporation"
)

// Identify tries to identify Azure provider by reading the /sys/class/dmi/id/sys_vendor file
func IdentifyAzure() string {
	return fileContains(AzureFile, csp.Azure, AzureContents)
}

// IdentifyAzureViaMetadataServer tries to identify Azure via metadata server
func IdentifyAzureViaMetadataServer() string {
	headers := map[string]string{"Metadata": "true"}
	resp, err := httpGet(AzureURL, headers)
	if err != nil {
		Logger.Printf("something happened during the request %v", err)
		return ""
	}
	if resp.StatusCode == http.StatusOK {
		return csp.Azure
	}
	Logger.Printf("something happened during the request with status %s", resp.Status)
	return ""
}
