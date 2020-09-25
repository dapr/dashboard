FROM alpine:latest as alpine

ARG BIN_PATH
ADD $BIN_PATH/dashboard .
ADD $BIN_PATH/web .
ENTRYPOINT ["./dashboard"]

