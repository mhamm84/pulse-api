
FROM golang:1.18

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
ADD /bin/linux_amd64/api /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
#RUN go build -o server .

# Run the server executable
ENTRYPOINT ["/app/api"]
CMD ["-cors-trusted-origins=$PULSE_CORS_TRUSTED_ORIGIN"]