package usecase

import "github.com/m-mizutani/octovy/backend/pkg/domain/model"

func (x *Default) GetOctovyMetadata() *model.Metadata {
	return &x.config.Metadata
}
