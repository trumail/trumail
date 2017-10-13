package verifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const updateInterval = 30 * time.Minute

var lists = []string{
	"https://raw.githubusercontent.com/wesbos/burner-email-providers/master/emails.txt",
	"https://gist.githubusercontent.com/adamloving/4401361/raw/db901ef28d20af8aa91bf5082f5197d27926dea4/temporary-email-address-domains",
	"https://gist.githubusercontent.com/michenriksen/8710649/raw/e09ee253960ec1ff0add4f92b62616ebbe24ab87/disposable-email-provider-domains",
	"https://raw.githubusercontent.com/martenson/disposable-email-domains/master/disposable_email_blacklist.conf",
	"https://raw.githubusercontent.com/andreis/disposable/master/domains.txt",
	"https://raw.githubusercontent.com/jamesaustin/disposable-email-domains/master/disposable-email-domains.txt",
	"https://raw.githubusercontent.com/flotwig/disposable-email-addresses/master/domains.txt",
	"https://raw.githubusercontent.com/FGRibreau/mailchecker/master/list.json",
}

// Disposabler defines all functionality for checking if an email
// address is disposable
type Disposabler interface {
	IsDisposable(domain string) bool
}

// disposabler contains the map of known disposable email domains
type disposabler struct {
	sync.RWMutex
	disposable map[string]bool
}

// NewDisposabler creates a new Disposabler and starts a domain farmer
// that retrieves all known disposable domains periodically
func NewDisposabler() Disposabler {
	// Generates a new domainMap on a Disposabler and appends domains
	d := &disposabler{disposable: make(map[string]bool)}

	// Retrieves new disposable lists every hour
	go d.domainFarmer()
	return d
}

// IsDisposable tests whether a string is among the known set of disposable
// mailbox domains. Returns true if the address is disposable
func (d *disposabler) IsDisposable(domain string) bool {
	d.RLock()
	defer d.RUnlock()
	return d.disposable[domain]
}

// domainFarmer retrieves new disposable domains every set interval
func (d *disposabler) domainFarmer() error {
	for {
		for _, url := range lists {
			// Performs the request for the domain list
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Reads the body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			// Adds every domain to our disposable domain map
			d.Lock()
			if strings.Contains(url, "FGRibreau/mailchecker") {
				re := regexp.MustCompile(`//.*`)
				res := re.ReplaceAllString(string(body), "")

				var domains [][]string
				if err := json.Unmarshal([]byte(res), &domains); err != nil {
					return err
				}

				for _, group := range domains {
					for _, domain := range group {
						d.disposable[strings.TrimSpace(domain)] = true
					}
				}
			} else {
				for _, domain := range strings.Split(string(body), "\n") {
					d.disposable[strings.TrimSpace(domain)] = true
				}
			}
			d.Unlock()
		}
		time.Sleep(updateInterval)
	}
}
