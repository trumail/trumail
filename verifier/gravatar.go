package verifier

import (
	"fmt"
)

const gravatarBaseURL = "https://en.gravatar.com"

// HasGravatar performs an http HEAD request to check if the email is
// associated with a gravatar account
func (v *Verifier) HasGravatar(md5Hash string) bool {
	u := fmt.Sprintf("%s/%s.json", gravatarBaseURL, md5Hash)
	return v.client.Head(u) == nil
}
