package providers

import (
	"strings"

	"github.com/paulstuart/satellite/csp"
)

//Used doc
//https://docs.cloud.oracle.com/iaas/Content/Compute/Tasks/gettingmetadata.htm

const (
	OracleURL      = "http://169.254.169.254/opc/v1/instance/metadata/"
	OracleFile     = "/sys/class/dmi/id/chassis_asset_tag"
	OracleContents = "OracleCloud"
)

type oracleMetadataResponse struct {
	OkeTM string `json:"oke-tm"`
}

func (or *oracleMetadataResponse) IsCSP() string {
	if strings.Contains(or.OkeTM, "oke") {
		return csp.Oracle
	}
	return ""
}

// Identify tries to identify Oracle provider by reading the /sys/class/dmi/id/chassis_asset_tag file
func IdentifyOracle() (string, error) {
	if fileContains(OracleFile, OracleContents) {
		return csp.Oracle, nil
	}
	return "", nil
}

// IdentifyOracleViaMetadataServer tries to identify Oracle via metadata server
func IdentifyOracleViaMetadataServer() (string, error) {
	var r oracleMetadataResponse
	return IdentifyViaMetadataServer(OracleURL, &r)
}
