package satellite

import (
	"sync"

	"github.com/paulstuart/satellite/providers"
)

// NOTE: we don't really care about errors, as they're expected for all but our own CSP
// and the first non-empty string response wins

type determiner func() string

func pickOne(checker ...determiner) string {
	var csp string
	var mu sync.Mutex // perhaps overkill as we only should get a single response
	var wg sync.WaitGroup
	wg.Add(len(checker))
	for _, fn := range checker {
		go func(fx determiner) {
			if s := fx(); s != "" {
				mu.Lock()
				csp = s
				mu.Unlock()
			}
			wg.Done()
		}(fn)
	}
	wg.Wait()
	return csp
}

func DetermineCSP() string {
	checker := []determiner{
		providers.IdentifyAlibaba,
		providers.IdentifyAmazon,
		providers.IdentifyAzure,
		providers.IdentifyDigitalOcean,
		providers.IdentifyGoogle,
		providers.IdentifyOracle,
	}
	if csp := pickOne(checker...); csp != "" {
		return csp
	}
	providers.Logger.Printf("no local CSP identifiers found")
	checker = []determiner{
		providers.IdentifyAlibabaViaMetadataServer,
		providers.IdentifyAmazonViaMetadataServer,
		providers.IdentifyAzureViaMetadataServer,
		providers.IdentifyDigitalOceanViaMetadataServer,
		providers.IdentifyGoogleViaMetadataServer,
		providers.IdentifyOracleViaMetadataServer,
	}
	return pickOne(checker...)
}
