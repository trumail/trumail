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

// Disposabler contains the map of known disposable email domains
type Disposabler struct {
	sync.RWMutex
	client     *http.Client
	disposable map[string]bool
}

// NewDisposabler creates a new Disposabler and starts a domain farmer
// that retrieves all known disposable domains periodically
func NewDisposabler(client *http.Client) *Disposabler {
	// Generates a new domainMap on a Disposabler and appends domains
	d := &Disposabler{
		client:     client,
		disposable: make(map[string]bool),
	}

	// Retrieves new disposable lists every hour
	go d.farmDomains()
	return d
}

// IsDisposable tests whether a string is among the known set of disposable
// mailbox domains. Returns true if the address is disposable
func (d *Disposabler) IsDisposable(domain string) bool {
	d.RLock()
	defer d.RUnlock()
	return d.disposable[domain]
}

// farmDomains retrieves new disposable domains every set interval
func (d *Disposabler) farmDomains() error {
	for {
		for _, url := range lists {
			// Performs the request for the domain list
			body, err := d.get(url)
			if err != nil {
				continue
			}

			// Adds every domain to our disposable domain map
			d.Lock()
			if strings.Contains(url, "FGRibreau/mailchecker") {
				re := regexp.MustCompile(`//.*`)
				res := re.ReplaceAllString(string(body), "")

				var domains [][]string
				if err := json.Unmarshal([]byte(res), &domains); err != nil {
					continue
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

// get performs a get request using the passed URL
func (d *Disposabler) get(url string) ([]byte, error) {
	// Perform the GET request on the passed URL
	resp, err := d.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode and return the bytes
	return ioutil.ReadAll(resp.Body)
}
