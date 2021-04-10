
jammer:
	go build \
		-o ./target/$@ \
		-trimpath \
		./

.PHONY: docker
docker:
	docker build -t jammer:local -f .docker/Dockerfile .