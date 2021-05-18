FROM hashicorp/terraform:latest as terraform
FROM golang:1.16 as builder

COPY --from=terraform /bin/terraform /bin/terraform
COPY --from=terraform /bin/echo /usr/bin/echo
RUN mkdir /src
COPY ./src /src

RUN cd /src/ && go test ./...