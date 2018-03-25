GIT_VERSION = $(shell git describe --always)
COMPONENT=registry-aws-ecr

# build:
# 	docker build -t civelocity/${REPOSITORY}:${GIT_VERSION} .
# 	docker tag civelocity/${REPOSITORY}:${GIT_VERSION} civelocity/${REPOSITORY}:latest
dep-install:
	@docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/velocity-ci/${COMPONENT} \
	--workdir /go/src/github.com/velocity-ci/${COMPONENT} \
	golang:1.9 \
	scripts/install-deps.sh

build:
	@docker run --rm \
	--volume ${CURDIR}:/go/src/github.com/velocity-ci/${COMPONENT} \
	--workdir /go/src/github.com/velocity-ci/${COMPONENT} \
	golang:1.9-alpine \
	scripts/build.sh

publish: build
	docker push civelocity/${REPOSITORY}:${GIT_VERSION}
	docker push civelocity/${REPOSITORY}:latest