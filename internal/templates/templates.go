package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nghiadaulau/opsbrew/internal/config"
)

// Template represents a project template
type Template struct {
	Name        string
	Description string
	Files       []TemplateFile
}

// TemplateFile represents a file in a template
type TemplateFile struct {
	Path     string
	Content  string
	IsDir    bool
	Mode     os.FileMode
}

// GetAvailableTemplates returns all available templates
func GetAvailableTemplates() []Template {
	return []Template{
		{
			Name:        "github-actions",
			Description: "GitHub Actions workflow template",
			Files:       getGitHubActionsFiles(),
		},
		{
			Name:        "k8s-deployment",
			Description: "Kubernetes Deployment manifest",
			Files:       getK8sDeploymentFiles(),
		},
		{
			Name:        "k8s-service",
			Description: "Kubernetes Service manifest",
			Files:       getK8sServiceFiles(),
		},
		{
			Name:        "k8s-pod",
			Description: "Kubernetes Pod manifest",
			Files:       getK8sPodFiles(),
		},
		{
			Name:        "k8s-configmap",
			Description: "Kubernetes ConfigMap manifest",
			Files:       getK8sConfigMapFiles(),
		},
		{
			Name:        "dockerfile",
			Description: "Multi-stage Dockerfile template",
			Files:       getDockerfileFiles(),
		},
	}
}

// InitializeTemplate initializes a new project from template
func InitializeTemplate(templateName, projectName, outputDir string, force bool, cfg *config.Config) error {
	// Find template
	var selectedTemplate *Template
	templates := GetAvailableTemplates()
	for _, t := range templates {
		if t.Name == templateName {
			selectedTemplate = &t
			break
		}
	}

	if selectedTemplate == nil {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Determine output directory
	if outputDir == "" {
		if projectName != "" {
			outputDir = projectName
		} else {
			outputDir = "."
		}
	}

	// Create output directory if it doesn't exist
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Template data
	data := map[string]interface{}{
		"ProjectName": projectName,
		"ModuleName":  strings.ToLower(strings.ReplaceAll(projectName, "-", "")),
		"ServiceName": projectName,
	}

	// Create files
	for _, file := range selectedTemplate.Files {
		filePath := filepath.Join(outputDir, file.Path)
		
		// Check if file exists
		if _, err := os.Stat(filePath); err == nil && !force {
			return fmt.Errorf("file %s already exists (use --force to overwrite)", filePath)
		}

		if file.IsDir {
			// Create directory
			if err := os.MkdirAll(filePath, file.Mode); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filePath, err)
			}
		} else {
			// Create file
			dir := filepath.Dir(filePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			// Parse and execute template
			tmpl, err := template.New(filePath).Parse(file.Content)
			if err != nil {
				return fmt.Errorf("failed to parse template for %s: %w", filePath, err)
			}

			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", filePath, err)
			}
			defer f.Close()

			if err := tmpl.Execute(f, data); err != nil {
				return fmt.Errorf("failed to execute template for %s: %w", filePath, err)
			}

			// Set file permissions
			if err := os.Chmod(filePath, file.Mode); err != nil {
				return fmt.Errorf("failed to set permissions for %s: %w", filePath, err)
			}
		}
	}

	return nil
}

func getGitHubActionsFiles() []TemplateFile {
	return []TemplateFile{
		{
			Path: ".github/workflows/ci.yml",
			Content: `name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Download dependencies
      run: go mod download

    - name: Run tests
      run: go test -v ./...

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v4

    - name: Build application
      run: go build -o {{.ServiceName}} .

    - name: Upload artifact
      uses: actions/upload-artifact@v3
      with:
        name: {{.ServiceName}}-binary
        path: {{.ServiceName}}`,
			Mode: 0644,
		},
	}
}

func getK8sDeploymentFiles() []TemplateFile {
	return []TemplateFile{
		{
			Path: "deployment.yaml",
			Content: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ServiceName}}
  labels:
    app: {{.ServiceName}}
spec:
  replicas: 2
  selector:
    matchLabels:
      app: {{.ServiceName}}
  template:
    metadata:
      labels:
        app: {{.ServiceName}}
    spec:
      containers:
      - name: {{.ServiceName}}
        image: {{.ServiceName}}:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: ENVIRONMENT
          value: "development"
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: http
          initialDelaySeconds: 5
          periodSeconds: 5`,
			Mode: 0644,
		},
	}
}

func getK8sServiceFiles() []TemplateFile {
	return []TemplateFile{
		{
			Path: "service.yaml",
			Content: `apiVersion: v1
kind: Service
metadata:
  name: {{.ServiceName}}-service
  labels:
    app: {{.ServiceName}}
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: {{.ServiceName}}`,
			Mode: 0644,
		},
	}
}

func getK8sPodFiles() []TemplateFile {
	return []TemplateFile{
		{
			Path: "pod.yaml",
			Content: `apiVersion: v1
kind: Pod
metadata:
  name: {{.ServiceName}}-pod
  labels:
    app: {{.ServiceName}}
spec:
  containers:
  - name: {{.ServiceName}}
    image: {{.ServiceName}}:latest
    ports:
    - containerPort: 8080
      name: http
    env:
    - name: ENVIRONMENT
      value: "development"
    - name: LOG_LEVEL
      value: "info"
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"`,
			Mode: 0644,
		},
	}
}

func getK8sConfigMapFiles() []TemplateFile {
	return []TemplateFile{
		{
			Path: "configmap.yaml",
			Content: `apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.ServiceName}}-config
data:
  config.yaml: |
    port: 8080
    env: development
    
    database:
      host: localhost
      port: 5432
      name: {{.ServiceName}}
    
    logging:
      level: info
      format: json
    
    features:
      debug: true
      metrics: true`,
			Mode: 0644,
		},
	}
}

func getDockerfileFiles() []TemplateFile {
	return []TemplateFile{
		{
			Path: "Dockerfile",
			Content: `# Multi-stage build for {{.ServiceName}}

# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{.ServiceName}} .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/{{.ServiceName}} .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./{{.ServiceName}}"]`,
			Mode: 0644,
		},
	}
}
