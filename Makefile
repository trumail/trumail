NOVENDOR_PATH = $$(glide novendor)

runtest:
	go clean
	go test ${NOVENDOR_PATH}
