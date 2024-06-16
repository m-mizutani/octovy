package octovy

import "time"

#IgnoreConfig: {
	Target: string
	Vulns: [...#IgnoreVuln] @go(,[]IgnoreVuln)
}

#IgnoreVuln: {
	ID:        string
	Comment?:  string
	ExpiresAt: time.Time
}

IgnoreList?: [...#IgnoreConfig]
