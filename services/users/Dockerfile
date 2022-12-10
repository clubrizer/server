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
RUN go build -o /user-service .

###
### Deploy
###
FROM gcr.io/distroless/base-debian11:sha256-cd1bf874ac029cfca6d6a8221f79bb435c5223a3d3eb717e045ca5b0163f2a6c.sig

ENV PORT=${PORT}

WORKDIR /

COPY --from=build /user-service /user-service

EXPOSE ${PORT}

USER nonroot:nonroot

ENTRYPOINT ["/user-service"]