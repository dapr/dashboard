FROM scratch
WORKDIR /app
COPY release .
CMD ["./dashboard"]