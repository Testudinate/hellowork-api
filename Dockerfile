FROM golang:1.7-alpine

# Install make and curl
RUN apk --update add make curl git

# Set apps home directory.
ENV APP_DIR ${GOPATH}/src/github.com/italolelis/hellowork-api

# Creates the application directory
RUN mkdir -p $APP_DIR

# Add sources.
COPY . $APP_DIR

# Define current working directory.
WORKDIR $APP_DIR

# Build the go binary
RUN make

# Clean apk cache
RUN rm -rf /var/cache/apk/*

CMD hellowork

