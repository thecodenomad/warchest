FROM golang:latest

RUN apt-get update -y && apt-get install -y make nodejs npm

# Copy all the bits
COPY . /code

WORKDIR /code

# Build the go side of the project
RUN make build

RUN export NODE_OPTIONS=--openssl-legacy-provider

# Build the UI and copy the build to the staticly served folder
RUN cd /code/warchest-ui && \
    # Install node dependencies
    npm install

RUN cd /code/warchest-ui && \
    # Build the uil into dist/
    npm run build && \
    cp -r dist/* /code/public/

# Kick off service
CMD ["./warchest", "--server"]