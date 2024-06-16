package octovy

import "time"

#IgnoreConfig: {
	Target: string
	Vulns: [...#IgnoreVuln] @go(,[]IgnoreVuln)
}

#IgnoreVuln: {
	ID:           string
	Description?: string
	ExpiresAt:    time.Time
}

IgnoreList?: [...#IgnoreConfig]
