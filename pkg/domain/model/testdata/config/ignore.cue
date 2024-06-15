package octovy

IgnoreTargets: [
	{
		File: "test.data"
		Vulns: [
			{
				ID:          "CVE-2017-9999"
				Description: "This is test data"
				ExpiresAt:   "2018-01-01T00:00:00Z"
			},
		]
	},
	{
		File: "test2.data"
		Vulns: [
			{
				ID:          "CVE-2017-11423"
				Description: "Hoge"
				ExpiresAt:   "2022-03-04T00:00:00Z"
			},
						{
				ID:          "CVE-2023-11423"
				Description: "Hoge"
				ExpiresAt:   "2023-03-04T00:00:00Z"
			},
		]
	},
]
