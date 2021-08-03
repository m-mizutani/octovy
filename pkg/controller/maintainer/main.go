package maintainer

import "github.com/m-mizutani/octovy/pkg/domain/interfaces"

type Maintainer struct {
	usecase interfaces.Usecases
}

func New(uc interfaces.Usecases) *Maintainer {
	return &Maintainer{
		usecase: uc,
	}
}

func (x *Maintainer) Run() error {
	return nil
}
