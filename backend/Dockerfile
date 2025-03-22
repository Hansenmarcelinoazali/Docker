FROM golang:1.22.0

WORKDIR /app

# Copy semua file proyek ke dalam container
COPY . .

# Download dependencies
RUN go mod tidy

# Build aplikasi
RUN go build -o main

# Jalankan aplikasi
CMD ["./main"]
