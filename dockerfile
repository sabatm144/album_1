# The base go-image
# FROM golang:1.14-alpine

# get the nsq for message queue
FROM nsqio/nsq

# Add Maintainer Info
LABEL maintainer="Monalisha sabat<sabatm144@gmail.com>"

# Create a directory for the app
RUN mkdir /app

# Create a directory to store the album
RUN mkdir /app/gallery

# Copy all files from the current directory to the app directory
COPY album /app

# Set working directory
WORKDIR /app

# Run the server executable
CMD [ "/app/album" ]