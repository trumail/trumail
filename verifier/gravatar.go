package verifier

import (
	"fmt"
)

const gravatarBaseURL = "https://en.gravatar.com"

// HasGravatar performs an http HEAD request to check if the email is
// associated with a gravatar account
func (v *Verifier) HasGravatar(a *Address) bool {
	u := fmt.Sprintf("%s/%s.json", gravatarBaseURL, a.MD5Hash)
	return v.client.Head(u) == nil
}
