FROM golang:1.19.3 AS build-go
ADD . /app
WORKDIR /app
RUN go build .

#gcr.io/distroless/static is not enough because of github.com/mattn/go-sqlite3
FROM gcr.io/distroless/base
COPY --from=build-go /app/octovy /octovy
COPY --from=aquasec/trivy:0.44.1 /usr/local/bin/trivy /trivy
WORKDIR /
ENV OCTOVY_ADDR="0.0.0.0"
ENV OCTOVY_TRIVY_PATH=/trivy
EXPOSE 9080
ENTRYPOINT ["/octovy"]
