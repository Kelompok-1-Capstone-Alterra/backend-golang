# syntax=docker/dockerfile:1

# specify the base image to be used for the application, alpine or ubuntu
FROM golang:1.17-alpine

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./
COPY go.sum ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . ./

# compile application
RUN go build -o /backend-golang

# command to be used to execute when the image is used to start a container
CMD [ "/backend-golang" ]