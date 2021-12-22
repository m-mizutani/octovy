package server

type Option func(cfg *config)

func DisableAuth() Option {
	return func(cfg *config) {
		cfg.DisableAuth = true
	}
}

func DisableWebhookGitHub() Option {
	return func(cfg *config) {
		cfg.DisableWebhookGitHub = true
	}
}

func DisableWebhookTrivy() Option {
	return func(cfg *config) {
		cfg.DisableWebhookTrivy = true
	}
}

func DisableFrontend() Option {
	return func(cfg *config) {
		cfg.DisableFrontend = true
	}
}
