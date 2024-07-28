# Use the official Golang image as the base image
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy files to container
COPY . .

# Install ffmpeg and yt-dlp
RUN apt-get update
RUN apt-get install -y -qq ffmpeg yt-dlp

# Build the Go app
RUN make build

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/bot"]