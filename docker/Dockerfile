FROM gcr.io/distroless/static:nonroot

ARG BIN_PATH
EXPOSE 8080/tcp

WORKDIR /app
ADD $BIN_PATH/dashboard .
ADD $BIN_PATH/web ./web
ENTRYPOINT ["./dashboard"]

