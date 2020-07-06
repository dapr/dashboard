# Build the code using build.sh file before doing docker build
FROM scratch
WORKDIR /app
COPY release/linux_amd64 .
CMD ["./dashboard"]