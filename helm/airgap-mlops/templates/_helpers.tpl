{{/*
Common labels and helpers for the airgap-mlops chart.
*/}}

{{- define "airgap-mlops.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "airgap-mlops.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "airgap-mlops.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "airgap-mlops.labels" -}}
helm.sh/chart: {{ include "airgap-mlops.chart" . }}
app.kubernetes.io/name: {{ include "airgap-mlops.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
agm.airgap.mlops/component: "core"
{{- end -}}

{{- define "airgap-mlops.airgapNamespaceLabel" -}}
airgap.mlops/enforced: "true"
{{- end -}}
