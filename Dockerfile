FROM golang:1.21-bullseye AS build-go
COPY . /app
WORKDIR /app
# ENV CGO_ENABLED=0
RUN go get -v
RUN go build .

FROM gcr.io/distroless/base:nonroot
COPY --from=build-go /app/octovy /octovy
COPY --from=aquasec/trivy:0.45.1 /usr/local/bin/trivy /trivy
WORKDIR /
ENV OCTOVY_ADDR="0.0.0.0:8000"
ENV OCTOVY_TRIVY_PATH=/trivy
EXPOSE 8000
ENTRYPOINT [ "/octovy" ]
