{{ .Arranger }}
{{ .UpstreamDNS }}
{{ .Hostname }}
{{ .IP }}
{{ .GatewayIP }}
{{ .Netmask }}
{{ .K8sMasterIP }}
{{ .DockerRegistry }}
{{- range .HaPeer}}
{{ .IP }}
{{ .Hostname }}
{{- end }}