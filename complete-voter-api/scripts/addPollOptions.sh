#!/bin/bash
curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:2080/polls/poll/1/pollOption/1/description/Yes
curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:2080/polls/poll/1/pollOption/2/description/No
curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:2080/polls/poll/1/pollOption/3/description/Maybe