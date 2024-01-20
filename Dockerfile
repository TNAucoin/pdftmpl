# Build Stage - Golang dependencies and app build
FROM golang:1.21-alpine AS go-builder
WORKDIR /app
# Copy Go Module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy source files and build the app
COPY . .
RUN go mod tidy && \
    go build -o ./bin/api ./cmd/api

# Final Stage - Build a small runtime image
FROM python:3.11-slim-bullseye
WORKDIR /app
# Update and install system libraries
RUN apt-get update && apt-get install -y \
    # essential libraries for WeasyPrint
    build-essential \
    python3-dev \
    python3-setuptools \
    python3-wheel \
    python3-pip \
    python3-cffi \
    libcairo2 \
    libpango1.0-0 \
    libpangocairo-1.0-0 \
    libgdk-pixbuf2.0-0 \
    libffi-dev \
    shared-mime-info \
    fontconfig \
    fonts-dejavu \
    # clean up unused files
    && rm -rf /var/lib/apt/lists/*
# Install WeasyPrint
RUN pip install WeasyPrint
# Copy necessary files from build stages
COPY --from=go-builder /app/bin/api /api
COPY --from=go-builder /app/templates /templates
# Custom fonts
COPY /weasyprint/fonts/*.ttf /usr/share/fonts/
COPY /weasyprint/fonts/*.otf /usr/share/fonts/
# Update font cache
RUN fc-cache -f -v
# Expose required port
EXPOSE 4000
# Run the binary
CMD ["/api"]