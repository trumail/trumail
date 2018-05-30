# Trumail

[![CircleCI](https://circleci.com/gh/sdwolfe32/trumail.svg?style=svg)](https://circleci.com/gh/sdwolfe32/trumail)
[![GoDoc](https://godoc.org/github.com/sdwolfe32/trumail/verifier?status.svg)](https://godoc.org/github.com/sdwolfe32/trumail/verifier)

Trumail is a free and open source email validation/verification system. It is available in three forms, the Golang client library `verifier` for use in your own Go projects, a public API endpoint (more info: https://trumail.io), and a public Docker image on DockerHub (see: https://hub.docker.com/r/sdwolfe32/trumail/). 

NOTE: While we do offer a managed, enterprise level service to paying customers, it is highly recommended that you host the service yourself either using a Docker image or by forking and serving this project on your own instance. Please keep in mind, self-hosting Trumail requires bidirectional communication on port 25 which most residential ISPs restrict - AWS and Digitalocean both permit this sort of communication.

## Using the API (public or self-hosted)

Using the API is very simple. All that's needed to validate an address is to send a `GET` request using the below URL with one of our three supported formats (json/jsonp(with "callback" (all lowercase) queryparam)/xml).
```
https://api.trumail.io/v2/lookups/{format}?email={email}&token={token}
```

## Using the library

```go
package main

import (
	"log"

	trumail "github.com/sdwolfe32/trumail/verifier"
)

func main() {
  v := trumail.NewVerifier("YOUR_HOSTNAME.COM", "YOUR_EMAIL@DOMAIN.COM")
  
  // Validate a single address
  log.Println(v.Verify("test@gmail.com"))
}
```

## Running with Go

```
go get -d github.com/sdwolfe32/trumail/...
go install github.com/sdwolfe32/trumail
trumail
```

## Running with Docker

```
docker run -p 8080:8080 -e SOURCE_ADDR=my.email@gmail.com sdwolfe32/trumail
```

## How it Works

Verifying the deliverability of an email address isn't a very complicated process. In fact, the process Trumail takes to verify an address is really only half that of sending a standard email transmission and is outlined below...
```
First a TCP connection is formed with the MX server on port 25.

HELO my-domain.com              // We identify ourselves as my-domain.com (set via environment variable)
MAIL FROM: me@my-domain.com     // Set the FROM address being our own
RCPT TO: test-email@example.com // Set the recipient and receive a (200, 500, etc..) from the server
QUIT                            // Cancel the transaction, we have all the info we need
```
As you can see we first form a tcp connection with the mail server on port 25. We then identify ourselves as example.com and set a reply-to email of admin@example.com (both these are configured via the SOURCE_ADDR environment variable). The last, and obviously most important step in this process is the RCPT command. This is where, based on the response from the mail server, we are able to conclude the deliverability of a given email address. A 200 implies a valid inbox and anything else implies either an error with our connection to the mail server, or a problem with the address requested.

The BSD 3-clause License
========================

Copyright (c) 2018, Steven Wolfe. All rights reserved.

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
