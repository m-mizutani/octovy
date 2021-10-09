package trivy

import "github.com/aquasecurity/trivy/pkg/report"

type Mock struct {
	Interface
	ScanMock func(dir string) (*report.Report, error)
}

func NewMock() *Mock {
	return &Mock{}
}

func (x *Mock) Scan(dir string) (*report.Report, error) {
	return x.ScanMock(dir)
}
