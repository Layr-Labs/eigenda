

dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')

setup-geth:
	cd geth && rm -rf data && ./run.sh setup 10 ../secrets/ecdsa_keys/private_key_hex.txt ./secret/geth-account-password

start-geth:
	cd geth && docker compose up -d 
	./wait-for http://0.0.0.0:8545 -- echo "Geth up"

stop-geth:
	cd geth && docker compose down


.PHONY: new-anvil new-docker-anvil infra resources exp deploy-all stop-infra run-e2e run-e2e-nochurner run-e2e-nograph clean

new-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"

new-docker-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-docker.yaml "testdata/$(dt)/config.yaml"

infra:
	go run ./deploy/cmd -localstack-port 4570 infra

resources:
	go run ./deploy/cmd -localstack-port 4570 resources

exp: 
	go run ./deploy/cmd exp

deploy-all:
	go run ./deploy/cmd -localstack-port 4570 all

stop-infra:
	./bin.sh stop-anvil
	./bin.sh stop-graph
	docker stop $(shell docker container ls --all --filter=ancestor=localstack/localstack:latest --format "{{.ID}}")

run-e2e:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

run-e2e-nochurner:
	go test ./tests -v -config=../templates/testconfig-anvil-nochurner.yaml

run-e2e-nograph:
	go test ./tests -v -config=../templates/testconfig-anvil-nograph.yaml

clean:
	rm -rf testdata/*
