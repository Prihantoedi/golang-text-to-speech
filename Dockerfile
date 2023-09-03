######################################
# STEP 1 build executable binary
######################################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . . 
RUN go build -o tts

######################################
# STEP 2 build small image
######################################
FROM alpine
WORKDIR /app
COPY --from=builder /app/tts /app
CMD ["./tts"]