NOVENDOR_PATH = $$(glide novendor)
.PHONY: test

glide:
	-rm glide.lock
	-rm glide.yaml
	-rm -r vendor
	glide cache-clear
	glide init --non-interactive
	glide install

test:
	go clean
	go test ${NOVENDOR_PATH}

run:
	go run main.go
