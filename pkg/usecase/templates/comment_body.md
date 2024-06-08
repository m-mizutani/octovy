{{ .Signature }}
{{ if eq .Metadata.TotalVulnCount 0 }}
🎉 **No vulnerability detected** 🎉
{{ else if eq .Metadata.FixableVulnCount 0 }}
👍 **No fixable vulnerability detected** 👍
{{ end }}

{{ if .Added }}
## 🚨New Vulnerabilities
{{ range .Added }}
### {{ .Target }}
{{ range .Vulnerabilities }}
<details>
<summary>{{ .VulnerabilityID }}: {{ .Title }} ({{.Severity}})</summary>

- **PkgName**: {{ if .PkgName }}`{{ .PkgName }}`{{ else }}N/A{{ end }}
- **Installed Version**: {{ if .InstalledVersion }}`{{ .InstalledVersion }}`{{ else }}N/A{{ end }}
- **Fixed Version**: {{ if .FixedVersion }}`{{ .FixedVersion }}`{{ else }}N/A{{ end }}
- **Status**: {{ if .Status }}`{{ .Status }}`{{ else }}N/A{{ end }}
- **Severity**: {{ if .Severity }}`{{ .Severity }}`{{ else }}N/A{{ end }}

#### Description

{{ .Description }}

#### References
{{ range .References }}
- [{{ . }}]({{ . }}){{ end }}
</details>
{{ end }}{{ end }}{{ end }}

{{ if .Fixed }}
## ✅Fix Vulnerabilities
{{ range .Fixed }}
### {{ .Target }}
{{ range .Vulnerabilities }}
<details>
<summary>{{ .VulnerabilityID }}: {{ .Title }} ({{.Severity}})</summary>

- **PkgName**: {{ if .PkgName }}`{{ .PkgName }}`{{ else }}N/A{{ end }}
- **Installed Version**: {{ if .InstalledVersion }}`{{ .InstalledVersion }}`{{ else }}N/A{{ end }}
- **Fixed Version**: {{ if .FixedVersion }}`{{ .FixedVersion }}`{{ else }}N/A{{ end }}
- **Status**: {{ if .Status }}`{{ .Status }}`{{ else }}N/A{{ end }}
- **Severity**: {{ if .Severity }}`{{ .Severity }}`{{ else }}N/A{{ end }}

#### Description

{{ .Description }}

#### References

{{ range .References }}
- [{{ . }}]({{ . }}){{ end }}
</details>
{{ end }}{{ end }}{{ end }}

{{ if ne .Metadata.TotalVulnCount 0 }}
## All detected vulnerabilities
{{ range .Report.Results }}
<details>
<summary>{{ .Target }}: ({{ .Vulnerabilities | len }})</summary>

{{ range .Vulnerabilities }}- {{ .VulnerabilityID }}: ( `{{ .PkgName }}` ) {{ .Title }}
{{ end }}
</details>
{{ end }}
{{ end }}
