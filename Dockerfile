# FROM golang:latest
FROM golang:latest AS build-env

LABEL maintainer "ericotieno99@gmail.com"
LABEL vendor="Ekas Technologies"

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /go/ekas-portal-api

ENV GOOS=linux
ENV GOARCH=386
ENV CGO_ENABLED=0

# Copy the project in to the container
ADD . /go/ekas-portal-api

RUN go mod download 

# Go get the project deps
RUN go get github.com/ekas-portal-api

# Set the working environment.
ENV GO_ENV production

# Go install the project
# RUN go install github.com/ekas-portal-api
RUN go build

# Run the ekas-portal-api command by default when the container starts.
# ENTRYPOINT /go/bin/ekas-portal-api

FROM golang:latest
WORKDIR /go/

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /go/ekas-portal-api/ekas-portal-api /go/ekas-portal-api
COPY --from=build-env /go/ekas-portal-api/config/app.yaml /go/config/app.yaml
COPY --from=build-env /go/ekas-portal-api/config/errors.yaml /go/config/errors.yaml
COPY --from=build-env /go/ekas-portal-api/cert.pem /go/cert.pem
COPY --from=build-env /go/ekas-portal-api/key.pem /go/key.pem
RUN mkdir p logs  

# Use an unprivileged user.
# USER appuser

# Set the working environment.
ENV GO_ENV production

ENTRYPOINT ./ekas-portal-api

#Expose the port specific to the ekas API Application.
EXPOSE 8081
EXPOSE 8084


# FROM golang as builder
# WORKDIR /go/habibridho/simple-go/
# COPY . ./
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

# FROM alpine:latest
# WORKDIR /app/
# COPY --from=builder /go/habibridho/simple-go/simple-go /app/simple-go
# EXPOSE 8888
# ENTRYPOINT ./simple-go

