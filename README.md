# Octovy [![Go Report Card](https://goreportcard.com/badge/github.com/m-mizutani/octovy)](https://goreportcard.com/report/github.com/m-mizutani/octovy) [![Build Status](https://travis-ci.com/m-mizutani/octovy.svg?branch=master)](https://travis-ci.com/m-mizutani/octovy)

![SampleView](https://user-images.githubusercontent.com/605953/120887167-48f7eb80-c62c-11eb-877d-79f081367c81.png)
https://octovy.io

`Octovy` is a GitHub App to scan vulnerability of package system (such as RubyGems, NPM, etc.) for GitHub repository. It detects a package lock file such as `Gemfile.lock` and checks if the package includes vulnerability based on package version. After that, Octovy stores scan report to database that can be accessed via Web UI and sends a result to [GitHub Check](https://docs.github.com/en/rest/reference/checks) as CI. A conclusion of GitHub Check is only `success` (No vulnerable packages) or `neutral` (Vulnerable package found) for now.

![GitHub Check](https://user-images.githubusercontent.com/605953/120887551-82c9f180-c62e-11eb-8049-1f5e448b4dc5.png)

Basic idea of Octovy is based on [Trivy](https://github.com/aquasecurity/trivy).


## Acknowledge

`Octovy` is massively inspired by [Trivy](https://github.com/aquasecurity/trivy) and has a similar mechanism with trivy to detect vulnerability. Additionally Octovy leverages [trivy-db](https://github.com/aquasecurity/trivy-db) as vulnerability/advisory database. I appreciate trivy authors for publishing great OSS.
