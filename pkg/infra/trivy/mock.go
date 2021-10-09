package trivy

import (
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type Mock struct {
	Interface
	ScanMock func(dir string) (*model.TrivyReport, error)
}

func NewMock() *Mock {
	return &Mock{}
}

func (x *Mock) Scan(dir string) (*model.TrivyReport, error) {
	return x.ScanMock(dir)
}
