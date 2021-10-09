FROM node:16.10.0-buster AS build-node
ADD . /app
WORKDIR /app/assets
RUN npm i
RUN npm run export
#
FROM golang:1.16 AS build-go
COPY --from=build-node /app /app
WORKDIR /app
RUN go build .

#gcr.io/distroless/static is not enough because of github.com/mattn/go-sqlite3
FROM gcr.io/distroless/base
COPY --from=build-go /app/octovy /octovy
COPY --from=aquasec/trivy:0.20.0 /usr/local/bin/trivy /trivy
WORKDIR /
ENV OCTOVY_ADDR="0.0.0.0"
ENV OCTOVY_TRIVY_PATH=/trivy
EXPOSE 9080
ENTRYPOINT ["/octovy"]
