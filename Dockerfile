FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . /app/
RUN apk update \
    && apk add --no-cache git \
    && apk add --no-cache ca-certificates \
    && apk add --update gcc musl-dev \
    && update-ca-certificates
ENV CGO_ENABLED=1
RUN go build

FROM alpine
WORKDIR /app
ARG hb_username
ARG hb_password
COPY --from=build /app/rating-gain /app
ENV GIN_MODE release
ENV HB_USERNAME $hb_username
ENV HB_PASSWORD $hb_password
CMD [ "/app/rating-gain" ]
