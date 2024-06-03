ifeq ($(wildcard ../.git/*),)
$(warning semver disabled - building from release zip)
GITCOMMIT := ""
GITDATE := ""
SEMVER := $(shell basename $(CURDIR))
else
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITDATE := $(shell git log -1 --format=%cd --date=unix)
SEMVER := $(shell docker run --rm --volume "$(PWD)/../:/repo" gittools/gitversion:5.12.0 /repo -output json -showvariable SemVer)
ifeq ($(SEMVER), )
$(warning semver disabled - docker not installed)
SEMVER := "0.0.0"
endif
endif

RELEASE_TAG := $(or $(RELEASE_TAG),latest)

build: clean
	go mod tidy
	go build -o ./bin/node ./cmd

clean:
	rm -rf ./bin

docker: docker-node docker-plugin

docker-node:
	cd ../ && docker build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} . -t opr-node:${SEMVER} -t opr-node:${RELEASE_TAG} -f node/cmd/Dockerfile

docker-plugin:
	cd ../ && docker build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} . -t opr-nodeplugin:${SEMVER} -t opr-nodeplugin:${RELEASE_TAG} -f node/plugin/cmd/Dockerfile

semver:
	echo "${SEMVER}"
