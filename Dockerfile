FROM scratch
WORKDIR /app
COPY release .
CMD ["./dapr-ui"]