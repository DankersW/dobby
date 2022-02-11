FROM golang:1.16-alpine

# Copy and download dependency using go mod
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .
COPY docker/config.yml .
RUN go build -o ./out/home-automation-backend .


EXPOSE 3000
ENTRYPOINT ["./out/home-automation-backend"]
