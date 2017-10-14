# Trumail 

[![CircleCI](https://circleci.com/gh/sdwolfe32/trumail.svg?style=svg)](https://circleci.com/gh/sdwolfe32/trumail)
[![GoDoc](https://godoc.org/github.com/sdwolfe32/trumail/verifier?status.svg)](https://godoc.org/github.com/sdwolfe32/trumail/verifier)

Trumail is a free and open source email validation/verification system. It is available in three forms, the Golang client library `verifier` for use in your own Go projects, a public API endpoint (more info: https://trumail.io), and a public Docker image on DockerHub (see: https://hub.docker.com/r/sdwolfe32/trumail/). 

Our own API endpoint allows up to 1RPS per IP. This is to prevent abuse and to leave the tubes open for other users. If you perform requests in excess of the specified rate-limit a `429 Too Many Requests` will be returned instead.

NOTE: It is highly recommended (due to potential heroku IP blacklisting resulting in failed validations) that you host the service yourself either using a Docker image or by forking and serving this project on your own instance. However, self-hosting Trumail requires bidirectional communication on port 25 which most residential ISPs restrict. AWS and Digitalocean both allow this.

## Using the API (public or self-hosted)

Using the API is very simple. All that's needed to validate an address is to send a `GET` request using the below URL with one of our two supported formats (json/xml).
```
trumail.io/{format}/{email}
```

## Using the library

```
package main

import (
	"log"

	trumail "github.com/sdwolfe32/trumail/verifier"
)

func main() {
	v := trumail.NewVerifier(20, "YOUR_HOSTNAME.COM", "YOUR_EMAIL@DOMAIN.COM")
	res := v.Verify("test@gmail.com")
	log.Println(*res[0])
}
```

## Running with Docker

```
docker run -p 8000:8000 -e SOURCE_ADDR=my.email@gmail.com sdwolfe32/trumail
```

The BSD 3-clause License
========================

Copyright (c) 2017, Steven Wolfe. All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

 - Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

 - Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

 - Neither the name of Trumail nor the names of its contributors may
   be used to endorse or promote products derived from this software without
   specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
