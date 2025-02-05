API_FILE = docs/api.yaml
WEBHOOK_API = docs/webhook.yaml

# Generate service api models
.PHONY: models
models:
	rm -rf models && mkdir -p models && \
	docker run --rm --user `id -u`:`id -g` -e GOPATH=`go env GOPATH`:/go -v `pwd`:`pwd` -w `pwd` "quay.io/goswagger/swagger:v0.30.4" generate model --spec=$(API_FILE)

.PHONY: client_webhook
client_webhook:
	rm -rf clients/webhook && mkdir -p clients/webhook  && \
	docker run --rm --user `id -u`:`id -g` -e GOPATH=`go env GOPATH`:/go -v `pwd`:`pwd` -w `pwd` "quay.io/goswagger/swagger:v0.30.4" generate client --spec=$(WEBHOOK_API) --target=clients/webhook
