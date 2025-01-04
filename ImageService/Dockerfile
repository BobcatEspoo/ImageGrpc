FROM golang:1.23.1

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
WORKDIR /app

COPY ./go.mod ./go.sum ./ ./

COPY . .

RUN go build -o imageService

CMD ["./imageService"]