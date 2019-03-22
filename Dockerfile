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

WORKDIR /app/
COPY --from=builder /go/bin/ekas-portal-api /app/ekas-portal-api

# Set the working environment.
ENV GO_ENV production

#Expose the port specific to the ekas API Application.
EXPOSE 8081

# Run the ekas-portal-api command by default when the container starts.
ENTRYPOINT ./ekas-portal-api


# FROM golang as builder
# WORKDIR /go/src/github.com/habibridho/simple-go/
# COPY . ./
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

# FROM alpine:latest
# WORKDIR /app/
# COPY --from=builder /go/src/github.com/habibridho/simple-go/simple-go /app/simple-go
# EXPOSE 8888
# ENTRYPOINT ./simple-go

