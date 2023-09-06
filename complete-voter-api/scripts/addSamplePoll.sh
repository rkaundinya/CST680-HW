#!/bin/bash
curl -d '{ "PollID": 1, "PollTitle": "Testing", "PollQuestion": "Are you going to work"}' -H "Content-Type: application/json" -X POST http://localhost:2080/polls