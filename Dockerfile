FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /ndugu-backend cmd/server/main.go

EXPOSE 50051

CMD ["/ndugu-backend"] 
