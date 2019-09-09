# FROM golang:latest
FROM golang:1.11-alpine AS build-env

LABEL maintainer "ericotieno99@gmail.com"
LABEL vendor="Ekas Technologies"

WORKDIR /go/src/github.com/ekas-portal-api

ENV GOOS=linux
ENV GOARCH=386
ENV CGO_ENABLED=0

# Copy the project in to the container
ADD . /go/src/github.com/ekas-portal-api

# Go get the project deps
RUN go get github.com/ekas-portal-api

RUN mkdir -p /go/config
ADD ./config/app.yaml /go/config
ADD ./config/errors.yaml /go/config

# Set the working environment.
ENV GO_ENV production

# Go install the project
# RUN go install github.com/ekas-portal-api
RUN go build

# Run the ekas-portal-api command by default when the container starts.
# ENTRYPOINT /go/bin/ekas-portal-api

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /go/src/github.com/ekas-portal-api/ekas-portal-api /app/ekas-portal-api

ENTRYPOINT ./ekas-portal-api

#Expose the port specific to the ekas API Application.
#EXPOSE 8081


# FROM golang as builder
# WORKDIR /go/src/github.com/habibridho/simple-go/
# COPY . ./
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

# FROM alpine:latest
# WORKDIR /app/
# COPY --from=builder /go/src/github.com/habibridho/simple-go/simple-go /app/simple-go
# EXPOSE 8888
# ENTRYPOINT ./simple-go

