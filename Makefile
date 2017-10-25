NOVENDOR_PATH = $$(glide novendor)
.PHONY: test

glide:
	rm glide.lock
	rm glide.yaml
	rm -r vendor
	glide cache-clear
	glide init --non-interactive
	glide install

test:
	go clean
	go test ${NOVENDOR_PATH}

run:
	export SERVE_WEB='true'
	export ENVIRONMENT='local'
	export RATE_LIMIT='true'
	export PORT=8000
	go run main.go
