{{/*
Expand the name of the chart.
*/}}
{{- define "todo-app-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (e.g. DNS names).
*/}}
{{- define "todo-app-chart.fullname" -}}
{{- if .Values.deployment.app.name }}
{{- .Values.deployment.app.name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified app name for MongoDB.
*/}}
{{- define "todo-app-chart.mongodb.fullname" -}}
{{- if .Values.deployment.db.name }}
{{- .Values.deployment.db.name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- printf "%s-%s-mongo" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}