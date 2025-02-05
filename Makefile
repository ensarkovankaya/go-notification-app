API_FILE = api.yaml

# Generate api models
.PHONY: models
models:
	rm -rf models && mkdir -p models && \
	docker run --rm --user `id -u`:`id -g` -e GOPATH=`go env GOPATH`:/go -v `pwd`:`pwd` -w `pwd` "quay.io/goswagger/swagger:v0.30.4" generate model --spec=$(API_FILE)
