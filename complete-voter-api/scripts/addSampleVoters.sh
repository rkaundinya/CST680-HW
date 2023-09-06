#!/bin/bash
curl -d '{ "VoterID": 1, "FirstName": "Ram Eshwar", "LastName": "Kaundinya" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters
curl -d '{ "VoterID": 2, "FirstName": "John", "LastName": "Doe" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters