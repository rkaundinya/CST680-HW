Before beginning, cd into each api directory (voter-api, vote-api, and poll-api) and run ./build-docker.sh
This will build docker images for each api

You must do this step before running docker compose

Next, cd into docker directory
From here you can run 'docker compose up' to spin up docker containers

Finally, all testing scripts are in the scripting directory

There is also a make file you can use if you like with help descriptions on the various commands

But, you should be able to test all functionality using written scripts in the scripts directory

Note - I have kept some old code from the voter-api of last hw that deals with Vote History. While the struct
still contains this, I do not use it at all. In order to view votes, we can simply go to the votes url and path
and filter out for specific voters if desired. 