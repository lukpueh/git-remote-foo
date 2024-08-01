FROM golang:1.22

WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download

COPY . .

RUN go install github.com/lukpueh/git-remote-foo

COPY sshpass.sh test.sh .

ENTRYPOINT ["./test.sh"]