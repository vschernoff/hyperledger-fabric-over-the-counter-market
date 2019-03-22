.PHONY: up down generate clear docker

help:
	@echo "OTC Demo"
	@echo ""
	@echo "generate: generate required certificates and genesis block"
	@echo "up: bring up the network"
	@echo "down: clear the network"
	@echo ""

generate:
	./scripts/docker-frag-otc.sh
	./network.sh -m generate

up:
	./network.sh -m up

down:
	./network.sh -m down

clean:
	./network.sh -m clean
