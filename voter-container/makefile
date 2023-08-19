SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "     build           		Build the voter executable"
	@echo "     run             		Run the voter program from code"
	@echo " 	populate-sample-voters	fill database with sample voters and polls"
	@echo "     add-voter       		pass voterID=<ID> firstName=<firstName> lastName=<lastName>,"
	@echo " 							and add this voter to database"
	@echo " 	"
	@echo "		add-poll 				Add a poll to a voter's polls taken list,"
	@echo "								pass voter id using voterID=<ID> and poll id using pollID=<ID>"
	@echo " 	get-voter 				pass voter id using voterID=<ID>, get voter info for voter ID"
	@echo " 	get-voter-polls 		pass voterID=<ID>, get polls for voter"
	@echo " 	get-voter-specific-poll pass voterID=<ID> and pollID=<ID>, get specific poll for voter"
	@echo " 	update-voter 			pass voterID=<ID>, firstName=<name>, lastName=<name>; updates voter info"
	@echo " 	update-voter-poll 		pass voterID=<ID> and pollID=<ID>, update poll for voter with new timestamp"
	@echo " 	delete-voter 			pass voterID=<ID>, deletes this voter if found from database"
	@echo " 	delete-voter-poll 		pass voterID=<ID> and pollID=<ID>, deletes voter's poll"
	@echo " 	health-check 			returns a health check for voter app"

.PHONY: build
build:
	go build .

.PHONY: build-amd64-linux
build-amd64-linux:
	GOOS=linux GOARCH=amd64 go build -o ./voterapi-linux-amd64 .

.PHONY: build-arm64-linux
build-arm64-linux:
	GOOS=linux GOARCH=arm64 go build -o ./voterapi-linux-arm64 .

.PHONY: run
run:
	go run main.go

.PHONY: populate-sample-voters
populate-sample-voters:
	curl -d '{ "VoterID": 1, "FirstName": "Ram Eshwar", "LastName": "Kaundinya" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api
	curl -d '{ "VoterID": 2, "FirstName": "John", "LastName": "Doe" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api/voters/1/polls/1
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api/voters/1/polls/2
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api/voters/2/polls/3

.PHONY: add-voter
add-voter:
	curl -d '{ "VoterID": $(voterID), "FirstName": "$(firstName)", "LastName": "$(lastName)" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api/voters/$(voterID)/firstName/$(firstName)/lastName/$(lastName)

.PHONY: add-poll
add-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1080/voter-api/voters/$(voterID)/polls/$(pollID)

.PHONY: get-voter
get-voter:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voter-api/voters/$(voterID)

.PHONY: get-voter-polls
get-voter-polls:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voter-api/voters/$(voterID)/polls

.PHONY: get-voter-specific-poll
get-voter-specific-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voter-api/voters/$(voterID)/polls/$(pollID)

.PHONY: update-voter
update-voter:
	curl -d '{ "VoterID": $(voterID), "FirstName": "$(firstName)", "LastName": "$(lastName)" }' -H "Content-Type: application/json" -X PUT http://localhost:1080/voter-api/voters/$(voterID)

.PHONY: update-voter-poll
update-voter-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X PUT http://localhost:1080/voter-api/voters/$(voterID)/polls/$(pollID)

.PHONY: delete-voter
delete-voter:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voter-api/voters/$(voterID)

.PHONY: delete-voter-poll
delete-voter-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voter-api/voters/$(voterID)/polls/$(pollID)

.PHONY: health-check
health-check:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voter-api/voters/health