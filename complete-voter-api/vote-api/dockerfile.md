# syntax=docker/dockerfile:1

# Get go language container version 1.20
FROM golang:1.20 AS build-stage

# Creates directory called app and stores builds in there
WORKDIR /app

# copy files to directory
COPY . .

#downloads dependencies
RUN go mod download

#Build (CGO_ENABLED will statically link libraries)
RUN CGO_ENABLED=0 GOOS=linux go build -o /vote-api

FROM alpine:latest AS run-stage

# put alpine in the root
WORKDIR /

#Copy the binary from the build stage
COPY --from=build-stage /vote-api /vote-api

#Expose the port
EXPOSE 3080

ENV REDIS_URL=host.docker.internal:6379

#Run
CMD ["/vote-api"]