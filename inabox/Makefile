

dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')


.PHONY: new-anvil new-anvil-nochurner chain localstack exp deploy-all stop-infra run-e2e run-e2e-nochurner run-e2e-nograph clean

new-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"

new-anvil-nochurner:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-nochurner.yaml "testdata/$(dt)/config.yaml"

clean:
	rm -rf testdata/*

deploy-all:
	go run ./deploy/cmd  -localstack-port 4570 -deploy-resources true  all

stop-infra:
# 	We rm forcefully instead of using `docker stop` because it is much faster and we don't care about graceful shutdowns in tests.
	docker rm -f graph-node graph-postgres graph-ipfs anvil-inabox localstack-test 2>/dev/null || true
	docker network rm graph-network 2>/dev/null || true

run-e2e:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

run-e2e-nochurner:
	go test ./tests -v -config=../templates/testconfig-anvil-nochurner.yaml

run-e2e-nograph:
	go test ./tests -v -config=../templates/testconfig-anvil-nograph.yaml

### Debug indidividual deploy targets

# chain deploys both anvil and thegraph
deploy-anvil-and-graph:
	go run ./deploy/cmd -localstack-port 4570 chain

deploy-localstack:
	go run ./deploy/cmd -localstack-port 4570 localstack

deploy-exp: 
	go run ./deploy/cmd exp