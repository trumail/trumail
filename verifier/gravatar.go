package verifier

import (
	"fmt"
	"net/http"
)

const gravatarBaseURL = "https://en.gravatar.com"

// HasGravatar performs an http HEAD request to check if the email is
// associated with a gravatar account
func (v *Verifier) HasGravatar(a *Address) bool {
	u := fmt.Sprintf("%s/%s.json", gravatarBaseURL, a.MD5())
	resp, err := v.client.Head(u)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
