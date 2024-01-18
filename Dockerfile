# Specify Go version and Build Stage
FROM golang:1.21-alpine AS build-env

# Set Work Directory
WORKDIR /myapp

# Copy source files
COPY . .

# Install dependencies and build the app
RUN go mod tidy && \
    go build -o ./bin/api ./cmd/api/

# Final Stage - production environment
FROM alpine

# Create and switch to a user
#RUN adduser -D -g '' appuser
#USER appuser

# Copy only the built binary from the build stage
COPY --from=build-env /myapp/bin/api /bin/api

# Copy static assets
COPY --from=build-env /myapp/templates /templates

# Expose port for the app
EXPOSE 4000

# Run the binary
CMD ["/bin/api"]