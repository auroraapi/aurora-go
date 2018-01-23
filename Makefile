APP_NAME = aurora-go
GIT_HEAD = `git rev-parse --short HEAD`
ID_RSA   = $(shell cat ~/.ssh/id_rsa | tr '\n' '_')

ash:
	docker build --build-arg id_rsa='$(ID_RSA)' -f Dockerfile.build -t aurora-builder .
	docker run -v $$(pwd):/go/src/github.com/nkansal96/$(APP_NAME) -it aurora-builder

build:
	docker build --build-arg id_rsa='$(ID_RSA)' -t $(APP_NAME) .

run:
	docker run --net=host -it $(APP_NAME)
