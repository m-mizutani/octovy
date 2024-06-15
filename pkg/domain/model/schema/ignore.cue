package octovy

import "time"

#IgnoreTarget: {
	File: string
	Vulns: [...#IgnoreVuln] @go(,[]IgnoreVuln)
}

#IgnoreVuln: {
	ID:           string
	Description?: string
	ExpiresAt:    time.Time
}

IgnoreTargets?: [...#IgnoreTarget]
