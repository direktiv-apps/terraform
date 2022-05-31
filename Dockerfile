FROM golang:1.18-buster as build

COPY go.mod src/go.mod
COPY go.sum src/go.sum
RUN cd src/ && go mod download

COPY cmd src/cmd/
COPY models src/models/
COPY restapi src/restapi/

RUN cd src && \
    export CGO_LDFLAGS="-static -w -s" && \
    go build -tags osusergo,netgo -o /application cmd/terraform-server/main.go; 

FROM ubuntu:21.04

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install ca-certificates graphviz -y

RUN apt-get install -y gnupg software-properties-common curl git wget
RUN curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -
RUN apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
RUN apt-get update && apt-get install terraform

# DON'T CHANGE BELOW 
COPY --from=build /application /bin/application

EXPOSE 8080

CMD ["/bin/application", "--port=8080", "--host=0.0.0.0"]

