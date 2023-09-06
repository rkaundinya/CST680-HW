#!/bin/bash
curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:2080/polls/poll/1/pollOption/3
curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:2080/polls/poll/1/pollOption/3