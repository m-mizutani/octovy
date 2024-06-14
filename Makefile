all: pkg/domain/mock/infra.go pkg/domain/mock/usecase.go

cmd=go run github.com/matryer/moq@v0.3.4

pkg/domain/mock/infra.go: ./pkg/domain/interfaces/infra.go
	$(cmd) -out pkg/domain/mock/infra.go -pkg mock ./pkg/domain/interfaces GitHub BigQuery

pkg/domain/mock/usecase.go: ./pkg/domain/interfaces/usecase.go
	$(cmd) -out pkg/domain/mock/usecase.go -pkg mock ./pkg/domain/interfaces UseCase
