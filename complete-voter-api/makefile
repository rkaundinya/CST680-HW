SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "     build           		Build the voter executable"
	@echo "     run             		Run the voter program from code"
	@echo " 	populate-sample-voters	fill database with sample voters and polls"
	@echo "     add-sample-polls       	fill database with sample polls"
	@echo "     add-sample-votes       	fill database with sample votes"
	@echo "     add-vote       			pass voteID=<id>, pollID=<id>, voterID=<id>, voteVal=<id>"
	@echo "     add-poll       			pass pollID=<id>, title='title', question='question'"
	@echo "     add-poll-option       	pass pollID=<id>, optID=<id>, desc='description'"
	@echo "     delete-poll-option      pass pollID=<id>, optID=<id>"
	@echo " 	get-voter 				pass voter id using voterID=<ID>, get voter info for voter ID"
	@echo " 	get-poll 				pass poll id using id=<ID>, get poll info for poll ID"
	@echo " 	update-voter 			pass voterID=<ID>, firstName=<name>, lastName=<name>; updates voter info"
	@echo " 	delete-voter 			pass voterID=<ID>, deletes this voter if found from database"
	@echo " 	delete-poll 			pass pollID=<ID>, deletes poll"
	@echo " 	health-check 			returns a health check for voter api"

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
	curl -d '{ "VoterID": 1, "FirstName": "Ram Eshwar", "LastName": "Kaundinya" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters
	curl -d '{ "VoterID": 2, "FirstName": "John", "LastName": "Doe" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters

.PHONY: add-sample-polls
add-sample-polls:
	curl -d '{ "PollID": 1, "PollTitle": "Testing", "PollQuestion": "Are you going to work"}' -H "Content-Type: application/json" -X POST http://localhost:2080/polls

.PHONY: add-sample-votes
add-sample-votes:
	curl -d '{ "VoteID": 1, "VoterID": 1, "PollID": 1, "VoteValue": 2}' -H "Content-Type: application/json" -X POST http://localhost:3080/votes

.PHONY: add-vote
add-vote:
	curl -d '{ "VoteID": $(voteID), "VoterID": $(voterID), "PollID": $(pollID), "VoteValue": $(voteVal)}' -H "Content-Type: application/json" -X POST http://localhost:3080/votes/voteID/$(voteID)/voterID/$(voterID)/pollID/$(pollID)/voteVal/$(voteVal)

.PHONY: add-poll
add-poll:
	curl -d '{ "PollID": $(pollID), "PollTitle": $(title), "PollQuestion": $(question)}' -H "Content-Type: application/json" -X POST http://localhost:2080/polls

.PHONY: add-poll-option
add-poll-option:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:2080/polls/poll/$(pollID)/pollOption/$(optID)/description/$(desc)

.PHONY: delete-poll-option
delete-poll-option:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:2080/polls/poll/$(pollID)/pollOption/$(optID)

.PHONY: get-voter
get-voter:
	curl -d '{ "VoterID": $(id)}' -H "Content-Type: application/json" -X GET http://localhost:3080/votes/voter/$(id)

.PHONY: get-poll
get-poll:
	curl -d '{ "PollID": $(id)}' -H "Content-Type: application/json" -X GET http://localhost:2080/polls/poll/$(id)

.PHONY: update-voter
update-voter:
	curl -d '{ "VoterID": $(voterID), "FirstName": "$(firstName)", "LastName": "$(lastName)" }' -H "Content-Type: application/json" -X PUT http://localhost:1080/voter-api/voters/$(voterID)

.PHONY: delete-vote
delete-vote:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:3080/votes/vote/$(voteID)

.PHONY: delete-voter
delete-voter:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters/$(voterID)

.PHONY: delete-poll
delete-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:2080/polls/poll/$(pollID)

.PHONY: health-check
health-check:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/health