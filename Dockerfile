FROM golang:1.8

# Set the Current Working Directory inside the container
WORKDIR /go/src/app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["app"]