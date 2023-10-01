package providers

import (
	"strings"

	"github.com/paulstuart/satellite/csp"
)

// Used docs
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-identity-documents.html
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html

const (
	AmazonURL      = "http://169.254.169.254/latest/dynamic/instance-identity/document"
	AmazonFile     = "/sys/class/dmi/id/product_version"
	AmazonContents = "amazon"
)

type instanceIdentityResponse struct {
	ImageID    string `json:"imageId"`
	InstanceID string `json:"instanceId"`
}

func (r *instanceIdentityResponse) IsCSP() string {
	if strings.HasPrefix(r.ImageID, "ami-") &&
		strings.HasPrefix(r.InstanceID, "i-") {
		return csp.Amazon
	}
	return ""
}

// Identify tries to identify Amazon provider by reading the /sys/class/dmi/id/product_version file
func IdentifyAmazon() (string, error) {
	if fileContains(AmazonFile, AmazonContents) {
		return csp.Amazon, nil
	}
	return "", nil
}

// IdentifyAmazonViaMetadataServer tries to identify Amazon via metadata server
func IdentifyAmazonViaMetadataServer() (string, error) {
	var r instanceIdentityResponse
	return IdentifyViaMetadataServer(AmazonURL, &r)
}
