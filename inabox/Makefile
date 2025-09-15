dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')
.PHONY: new-anvil chain localstack exp deploy-all stop-infra run-e2e run-e2e-nochurner run-e2e-nograph clean

new-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"

new-anvil-nochurner:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-nochurner.yaml "testdata/$(dt)/config.yaml"

chain:
	go run ./deploy/cmd -localstack-port 4570 chain

localstack:
	go run ./deploy/cmd -localstack-port 4570 localstack

exp: 
	go run ./deploy/cmd exp

deploy-all:
	go run ./deploy/cmd  -localstack-port 4570 -deploy-resources true  all

stop-infra:
	./bin.sh stop-graph
# Stop anvil
	docker ps -q --filter "ancestor=ghcr.io/foundry-rs/foundry" | xargs -r docker stop 2>/dev/null || true
# Stop localstack based on container name
	docker ps -q --filter "ancestor=localstack/localstack" | xargs -r docker stop 2>/dev/null || true

run-e2e:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

clean:
	rm -rf testdata/*
