FROM golang:1.15.6-alpine3.12 as build-env

RUN mkdir /oauth-graphql-playground
RUN apk --update add ca-certificates build-base
RUN apk add make git
WORKDIR /oauth-graphql-playground
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go install ./...

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/bin/ /usr/local/bin/
WORKDIR /workspace
EXPOSE 5000

ENTRYPOINT ["/usr/local/bin/oauth-graphql-playground"]