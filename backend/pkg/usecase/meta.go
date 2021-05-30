package usecase

import "github.com/m-mizutani/octovy/backend/pkg/domain/model"

func (x *Default) GetOctovyMetadata() *model.OctovyMetadata {
	return &model.OctovyMetadata{
		FrontendURL: x.config.FrontendURL,
		AppURL:      x.config.GitHubAppURL,
		HomepageURL: x.config.HomepageURL,
	}
}
