package providers

import "github.com/paulstuart/satellite/csp"

//Used docs
// https://developers.digitalocean.com/documentation/metadata/#metadata-in-json

const (
	DigitalOceanURL      = "http://169.254.169.254/metadata/v1.json"
	DigitalOceanFile     = "/sys/class/dmi/id/sys_vendor"
	DigitalOceanContents = "DigitalOcean"
)

type digitalOceanMetadataResponse struct {
	DropletID int `json:"droplet_id"`
}

func (do *digitalOceanMetadataResponse) IsCSP() string {
	if do.DropletID > 0 {
		return csp.DigitalOcean
	}
	return ""
}

// Identify tries to identify DigitalOcean provider by reading the /sys/class/dmi/id/sys_vendor file
func IdentifyDigitalOcean() string {
	return fileContains(DigitalOceanFile, csp.DigitalOcean, DigitalOceanContents)
}

// IdentifyDigitalOceanViaMetadataServer tries to identify DigitalOcean via metadata server
func IdentifyDigitalOceanViaMetadataServer() string {
	var do digitalOceanMetadataResponse
	return IdentifyViaMetadataServer(DigitalOceanURL, &do)
}
