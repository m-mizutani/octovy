{{ .Signature }}
{{ if eq .Metadata.TotalVulnCount 0 }}
🎉 **No vulnerability detected** 🎉
{{ else if eq .Metadata.FixableVulnCount 0 }}
👍 **No fixable vulnerability detected** 👍
{{ end }}
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

{{ if ne .Metadata.TotalVulnCount 0 }}
## All detected vulnerabilities
{{ range .Report.Results }}
<details>
<summary>`{{ .Target }}` ({{ .Vulnerabilities | len }})</summary>

{{ range .Vulnerabilities }}- {{ .VulnerabilityID }}: ( `{{ .PkgName }}` ) {{ .Title }}
{{ end }}
</details>
{{ end }}
{{ end }}
