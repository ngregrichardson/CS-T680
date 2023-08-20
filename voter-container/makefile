SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build the Voter API executable"
	@echo "	   run					Run the Voter API from code. Requires Redis to be running on redis:6379 or REDIS_URL to be set"
	@echo "	   run-bin				Run the Voter API executable. Requires Redis to be running on redis:6379 or REDIS_URL to be set"
	@echo "	   docker-start			Build and start the Voter API as a Docker container, alongside Redis"
	@echo "	   docker-run			Run the Voter API as a Docker container, alongside Redis"
	@echo "	   docker-build			Build the Voter API as a Docker container"
	@echo "	   get-voters			Get all voters"
	@echo "	   get-voter			Get a voter using id=<id>"
	@echo "	   create-voter			Create a voter using id=<id>, firstName=<firstName>, lastName=<lastName>"
	@echo "	   update-voter			Update a voter using id=<id>, firstName=<firstName>, lastName=<lastName>"
	@echo "	   delete-voter			Delete a voter using id=<id>"
	@echo "	   get-voter-history	Get a voter's vote history using id=<id>"
	@echo "	   get-vote				Get a voter's vote using id=<id>, pollId=<pollId>"
	@echo "	   create-vote			Create a voter's vote using id=<id>, pollId=<pollId>"
	@echo "	   update-vote			Update a voter's vote using id=<id>, pollId=<pollId>"
	@echo "	   delete-vote			Delete a voter's vote using id=<id>, pollId=<pollId>"
	@echo "	   health				Health check using hang=<hang> (optional)"

.PHONY: build
build:
	go build .
	
.PHONY: run
run:
	go run main.go

.PHONY: docker-start
docker-start:
	docker compose up --build

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-run
docker-run:
	docker compose up

.PHONY: run-bin
run-bin:
	./voter-api

.PHONY: get-voters
get-voters:
	curl --location 'http://localhost:1080/v1/voters'

.PHONY: get-voter
get-voter:
	curl --location 'http://localhost:1080/v1/voters/$(id)'

.PHONY: create-voter
create-voter:
	curl --request POST --location 'http://localhost:1080/v1/voters/$(id)' --header 'Content-Type: application/json' --data '{ "firstName": "$(firstName)", "lastName": "$(lastName)" }'

.PHONY: update-voter
update-voter:
	curl --request PUT --location 'http://localhost:1080/v1/voters/$(id)' --header 'Content-Type: application/json' --data '{ "firstName": "$(firstName)", "lastName": "$(lastName)" }'

.PHONY: delete-voter
delete-voter:
	curl --request DELETE --location 'http://localhost:1080/v1/voters/$(id)'

.PHONY: get-voter-history
get-voter-history:
	curl --location 'http://localhost:1080/v1/voters/$(id)/votes'

.PHONY: get-vote
get-vote:
	curl --location 'http://localhost:1080/v1/voters/$(id)/votes/$(pollId)'

.PHONY: create-vote
create-vote:
	curl --request POST --location 'http://localhost:1080/v1/voters/$(id)/votes/$(pollId)'

.PHONY: update-vote
update-vote:
	curl --request PUT --location 'http://localhost:1080/v1/voters/$(id)/votes/$(pollId)'

.PHONY: delete-vote
delete-vote:
	curl --request DELETE --location 'http://localhost:1080/v1/voters/$(id)/votes/$(pollId)'

.PHONY: health
health:
	curl --location 'http://localhost:1080/v1/voters/health?hang=$(hang)'