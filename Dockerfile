FROM golang:latest

RUN apt-get update -y && apt-get install -y make nodejs npm git

# Copy all the bits
COPY . /code

WORKDIR /code

# Build the go side of the project
RUN make build

RUN export NODE_OPTIONS=--openssl-legacy-provider

# Make sure folder for static files exists
RUN mkdir -p /code/public

# Build the UI and copy the build to the staticly served folder
RUN cd /code/warchest-ui && \
    # Install node dependencies
    npm install && \
    # Build the uil into dist/
    npm run build && \
    cp -r dist/* /code/public/

# Kick off service
CMD ["./warchest", "--server"]