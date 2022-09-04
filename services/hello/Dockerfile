# syntax=docker/dockerfile:1

# See https://docs.docker.com/language/golang/build-images/#multi-stage-builds

##
## Build
##
FROM golang:1.19-buster AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o /hello-service .

###
### Deploy
###
FROM gcr.io/distroless/base-debian10

ENV PORT=${PORT}

WORKDIR /

COPY --from=build /hello-service /hello-service

EXPOSE ${PORT}

USER nonroot:nonroot

ENTRYPOINT ["/hello-service"]