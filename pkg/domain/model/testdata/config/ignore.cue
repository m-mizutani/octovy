package octovy

IgnoreList: [
	{
		Target: "test.data"
		Vulns: [
			{
				ID:        "CVE-2017-9999"
				Comment:   "This is test data"
				ExpiresAt: "2018-01-01T00:00:00Z"
			},
		]
	},
	{
		Target: "test2.data"
		Vulns: [
			{
				ID:        "CVE-2017-11423"
				Comment:   "Hoge"
				ExpiresAt: "2022-03-04T00:00:00Z"
			},
			{
				ID:        "CVE-2023-11423"
				Comment:   "Hoge"
				ExpiresAt: "2023-03-04T00:00:00Z"
			},
		]
	},
]
