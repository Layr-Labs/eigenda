build: clean
	go mod tidy
	go build -o ./bin/node ./cmd

clean:
	rm -rf ./bin

docker: docker-node docker-plugin

docker-node:
	cd ../ && docker build . -t opr-node -f node/cmd/Dockerfile

docker-plugin:
	cd ../ && docker build . -t opr-nodeplugin -f node/plugin/cmd/Dockerfile
