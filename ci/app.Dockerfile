FROM golang:1.16-alpine

# Copy and download dependency using go mod
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .
COPY ci/config.yml .
RUN go build -o ./out/dobby .


EXPOSE 3000
ENTRYPOINT ["./out/dobby"]
