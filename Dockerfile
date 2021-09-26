FROM node AS build-node
ADD . /app
WORKDIR /app/assets
RUN npm run export

FROM golang AS build-go
COPY --from=build-node /app /app
WORKDIR /app
RUN go build .

## gcr.io/distroless/static is not enough because of github.com/mattn/go-sqlite3
FROM gcr.io/distroless/base
COPY --from=build-go /app/octovy /
WORKDIR /
ENV OCTOVY_ADDR="0.0.0.0"
EXPOSE 9080
ENTRYPOINT ["/octovy"]
