
FROM golang:latest


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN apt-get update && apt-get install -y default-mysql-client


ENV DB_HOST=192.168.246.77  

ENV DB_USER=root
ENV DB_PASSWORD=9994570668@sri
ENV DB_NAME=app

RUN go build -o main .


EXPOSE 8081


CMD ["./main"]
