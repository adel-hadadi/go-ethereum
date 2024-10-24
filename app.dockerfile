FROM ubuntu:22.04

# Install SQLite3 and other necessary libraries
RUN apt-get update && apt-get install -y \
    sqlite3 \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the precompiled Go binary
COPY webApp .

# Copy environment file if needed
COPY .env .

# Expose the necessary port (if needed)
EXPOSE 8080

# Run the binary
CMD ["./webApp"]
