FROM hashicorp/terraform:latest as terraform
FROM golang:1.16 as builder

COPY --from=terraform /bin/terraform /bin/terraform
# Need /bin/echo for tests
COPY --from=terraform /bin/echo /usr/bin/echo
RUN mkdir /src
COPY ./src /src

RUN cd /src/ && go test ./...

RUN cd /src/cli && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /src/terrarun

FROM alpine:3.13
COPY --from=terraform /bin/terraform /usr/local/bin/terraform
COPY --from=builder  /src/terrarun /usr/local/bin/terrarun
RUN chmod a+x /usr/local/bin/terrarun

ENTRYPOINT ["/usr/local/bin/terrarun"]