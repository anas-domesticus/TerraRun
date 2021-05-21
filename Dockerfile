FROM hashicorp/terraform:latest as terraform
FROM golang:1.16 as tester

COPY --from=terraform /bin/terraform /bin/terraform
# Need /bin/echo for tests
COPY --from=terraform /bin/echo /usr/bin/echo
RUN mkdir /src
COPY ./src /src
WORKDIR /src
ARG GOOS=linux
ARG GOARCH=amd64
ENV GOOS=$GOOS
ENV GOARCH=$GOARCH
ENV CGO_ENABLED=0

RUN go install honnef.co/go/tools/cmd/staticcheck@latest

FROM tester as builder
RUN go build -o /src/terrarun cli/main.go

FROM alpine:3.13 as final
COPY --from=terraform /bin/terraform /usr/local/bin/terraform
COPY --from=builder  /src/terrarun /usr/local/bin/terrarun
RUN chmod a+x /usr/local/bin/terrarun

ENTRYPOINT ["/usr/local/bin/terrarun"]