{{ .Signature }}
{{ if .Added }}
## 🚨New Vulnerabilities
{{ range .Added }}
- `{{ .Target }}`
{{ range .Vulnerabilities }}    - {{ .VulnerabilityID }}: ({{ .PkgName }}) {{ .Title }}
{{ end }}{{ end }}{{ end }}

{{ if .Fixed }}
## ✅Fix Vulnerabilities
{{ range .Fixed }}
- `{{ .Target }}`
{{ range .Vulnerabilities }}    - {{ .VulnerabilityID }}: ({{ .PkgName }}) {{ .Title }}
{{ end }}{{ end }}{{ end }}

## All detected vulnerabilities
{{ range .Report.Results }}
<details>
<summary>`{{ .Target }}` ({{ .Vulnerabilities | len }})</summary>

{{ range .Vulnerabilities }}- {{ .VulnerabilityID }}: ( `{{ .PkgName }}` ) {{ .Title }}
{{ end }}
</details>
{{ end }}
