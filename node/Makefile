GITCOMMIT := $(shell git rev-parse --short HEAD)
GITDATE := $(shell git log -1 --format=%cd --date=unix)

# GitVersion provides the semantic versioning for the project. If docker is not installed, semver fallsback to 0.0.0
SEMVER := $(shell docker run --rm --volume "${PWD}/../:/repo" gittools/gitversion:5.12.0 /repo -output json -showvariable SemVer)
ifeq ($(SEMVER), )
SEMVER = "0.0.0"
endif

build: clean
	go mod tidy
	go build -o ./bin/node ./cmd

clean:
	rm -rf ./bin

docker: docker-node docker-plugin

docker-node:
	cd ../ && docker build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} . -t opr-node:${SEMVER} -t opr-node:latest -f node/cmd/Dockerfile

docker-plugin:
	cd ../ && docker build --build-arg SEMVER=${SEMVER} --build-arg GITCOMMIT=${GITCOMMIT} --build-arg GITDATE=${GITDATE} . -t opr-nodeplugin:${SEMVER} -t opr-nodeplugin:latest -f node/plugin/cmd/Dockerfile

semver:
	echo "${SEMVER}"
