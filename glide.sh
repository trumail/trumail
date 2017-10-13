rm glide.lock
rm glide.yaml
rm -r vendor
glide cache-clear
glide init --non-interactive
glide install