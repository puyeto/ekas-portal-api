# FROM golang:latest
FROM golang:1.10.3 as builder

LABEL maintainer "ericotieno99@gmail.com"
LABEL vendor="Ekas Technologies"

# Copy the project in to the container
ADD . /go/src/github.com/ekas-portal-api

# Go get the project deps
RUN go get github.com/ekas-portal-api

# Go install the project
RUN go install github.com/ekas-portal-api

FROM alpine:latest

RUN mkdir -p /go/app
# WORKDIR /go/app/
COPY --from=builder /go/bin/ekas-portal-api /go/app/ekas-portal-api

RUN mkdir -p /go/config
COPY --from=builder /go/src/github.com/ekas-portal-api/config/app.yaml /go/config/app.yaml
COPY --from=builder /go/src/github.com/ekas-portal-api/config/errors.yaml /go/config/errors.yaml

# Set the working environment.
ENV GO_ENV production

# Run the ekas-portal-api command by default when the container starts.
ENTRYPOINT /go/app/ekas-portal-api

#Expose the port specific to the ekas API Application.
EXPOSE 8081

