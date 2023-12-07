dt := $(shell date '+%YY-%mM-%dD-%HH-%MM-%SS')

new-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil.yaml "testdata/$(dt)/config.yaml"

new-docker-anvil:
	mkdir -p "testdata/$(dt)"
	cp ./templates/testconfig-anvil-docker.yaml "testdata/$(dt)/config.yaml"

config: 
	go run ./gen-env/cmd

devnet-up:
	$(eval CURRENT_TESTDATA := $(shell ls -1r testdata/ | head -n 1))
	docker-compose -f testdata/$(CURRENT_TESTDATA)/docker-compose.yml up --detach --build

devnet-logs:
	$(eval CURRENT_TESTDATA := $(shell ls -1r testdata/ | head -n 1))
	docker-compose -f testdata/$(CURRENT_TESTDATA)/docker-compose.yml logs --follow

devnet-down:
	$(eval CURRENT_TESTDATA := $(shell ls -1r testdata/ | head -n 1))
	docker-compose -f testdata/$(CURRENT_TESTDATA)/docker-compose.yml down --volumes

run-e2e:
	go test ./tests -v -config=../templates/testconfig-anvil.yaml

# run-e2e-nochurner:
# 	go test ./tests -v -config=../templates/testconfig-anvil-nochurner.yaml

# run-e2e-nograph:
# 	go test ./tests -v -config=../templates/testconfig-anvil-nograph.yaml

clean:
	rm -rf testdata/*


.PHONY: config