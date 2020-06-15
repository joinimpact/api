FROM golang

# Force go modules usage
ENV GO111MODULE=on

# Define workdir
WORKDIR /app

# Copy files
COPY . .

# Build through make
RUN make build

# Expose port 8080
EXPOSE 8080

# Define our server binary
ENTRYPOINT [ "/app/server" ]