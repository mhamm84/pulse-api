
FROM alpine

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
ADD /bin/linux_amd64/pulse /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
#RUN go build -o server .
EXPOSE 9091 9091
# Run the server executable
ENTRYPOINT ["/app/pulse", "run-api"]
#CMD ["run-api", "--cors-trusted-origins=$PULSE_CORS_TRUSTED_ORIGIN"]
