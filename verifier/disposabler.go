package verifier

import (
	"strings"
	"sync"
	"time"

	"github.com/sdwolfe32/trumail/httpclient"
)

// updateInterval is how often we should reach out to update
// the disposable address map
const updateInterval = 30 * time.Minute

// Disposabler contains the map of known disposable email domains
type Disposabler struct {
	client  *httpclient.Client
	dispMap *sync.Map
}

// NewDisposabler creates a new Disposabler and starts a domain farmer
// that retrieves all known disposable domains periodically
func NewDisposabler(client *httpclient.Client) *Disposabler {
	d := &Disposabler{client, &sync.Map{}}
	go d.farmDomains(updateInterval)
	return d
}

// IsDisposable tests whether a string is among the known set of disposable
// mailbox domains. Returns true if the address is disposable
func (d *Disposabler) IsDisposable(domain string) bool {
	_, ok := d.dispMap.Load(domain)
	return ok
}

// farmDomains retrieves new disposable domains every set interval
func (d *Disposabler) farmDomains(interval time.Duration) error {
	for {
		for _, url := range lists {
			// Perform the request for the domain list
			body, err := d.client.GetString(url)
			if err != nil {
				continue
			}

			// Split
			for _, domain := range strings.Split(body, "\n") {
				d.dispMap.Store(strings.TrimSpace(domain), true)
			}
		}
		time.Sleep(interval)
	}
}

// list is a slice of disposable email address lists
var lists = []string{
	"https://raw.githubusercontent.com/wesbos/burner-email-providers/master/emails.txt",
	"https://gist.githubusercontent.com/adamloving/4401361/raw/db901ef28d20af8aa91bf5082f5197d27926dea4/temporary-email-address-domains",
	"https://gist.githubusercontent.com/michenriksen/8710649/raw/e09ee253960ec1ff0add4f92b62616ebbe24ab87/disposable-email-provider-domains",
	"https://raw.githubusercontent.com/martenson/disposable-email-domains/master/disposable_email_blacklist.conf",
	"https://raw.githubusercontent.com/jamesaustin/disposable-email-domains/master/disposable-email-domains.txt",
	"https://raw.githubusercontent.com/flotwig/disposable-email-addresses/master/domains.txt",
	// "https://raw.githubusercontent.com/FGRibreau/mailchecker/master/list.json",
}
