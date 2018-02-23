APP_NAME = aurora-go
GIT_HEAD = `git rev-parse --short HEAD`
ID_RSA   = $(shell cat ~/.ssh/id_rsa | tr '\n' '_')

build-ash:
	docker build --build-arg id_rsa='$(ID_RSA)' -f Dockerfile.build -t aurora-builder .

ash: build-ash
	docker run -v "$$(pwd):/go/src/github.com/nkansal96/$(APP_NAME)" -it aurora-builder

check: build-ash
	docker run -v "$$(pwd):/go/src/github.com/nkansal96/$(APP_NAME)" -it aurora-builder /bin/ash -c "cd /go/src/github.com/nkansal96/$(APP_NAME) && go build . && gofmt -w ."

build:
	docker build --build-arg id_rsa='$(ID_RSA)' -t $(APP_NAME) .

run:
	docker run --net=host -it $(APP_NAME)
