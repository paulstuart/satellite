package satellite

import (
	"log"
	"sync"

	"github.com/paulstuart/satellite/providers"
)

type determiner func() (string, error)

func pickOne(checker ...determiner) string {
	var csp string
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(checker))
	for _, fn := range checker {
		go func(fx determiner) {
			if s, err := fx(); err == nil && s != "" {
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
		providers.IdentifyAsure,
		providers.IdentifyDigitalOcean,
		providers.IdentifyGoogle,
		providers.IdentifyOracle,
	}
	if csp := pickOne(checker...); csp != "" {
		return csp
	}
	log.Printf("no local CSP identifiers found")
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
