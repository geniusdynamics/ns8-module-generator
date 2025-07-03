package parser

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseDockerCompose(t *testing.T) {
	// Create a temporary Docker Compose file for testing
	tempFile, err := os.CreateTemp("", "docker-compose-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	composeContent := `
services:
  web:
    image: nginx:latest
    environment:
      - APP_ENV=production
      - APP_DEBUG=false
    volumes:
      - ./data:/var/www/html
  db:
    image: postgres:13
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - type: bind
        source: /data/db
        target: /var/lib/postgresql/data
volumes:
  data:
networks:
  default:
`
	if _, err := tempFile.WriteString(composeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	services, volumes, networks, err := ParseDockerCompose(tempFile.Name())
	if err != nil {
		t.Fatalf("ParseDockerCompose failed: %v", err)
	}

	// Test services
	if len(services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(services))
	}

	// Test 'web' service
	webService, ok := services["web"]
	if !ok {
		t.Error("Expected 'web' service, but not found")
	}
	if webService.Image != "nginx:latest" {
		t.Errorf("Expected web image 'nginx:latest', got '%s'", webService.Image)
	}
	expectedWebEnv := map[string]string{
		"APP_ENV":   "production",
		"APP_DEBUG": "false",
	}
	if !reflect.DeepEqual(webService.ParsedEnvironment, expectedWebEnv) {
		t.Errorf("Expected web environment %+v, got %+v", expectedWebEnv, webService.ParsedEnvironment)
	}
	expectedWebVolumes := []map[string]string{
		{"source": "./data", "target": "/var/www/html"},
	}
	if !reflect.DeepEqual(webService.ParsedVolumes, expectedWebVolumes) {
		t.Errorf("Expected web volumes %+v, got %+v", expectedWebVolumes, webService.ParsedVolumes)
	}

	// Test 'db' service
	dbService, ok := services["db"]
	if !ok {
		t.Error("Expected 'db' service, but not found")
	}
	if dbService.Image != "postgres:13" {
		t.Errorf("Expected db image 'postgres:13', got '%s'", dbService.Image)
	}
	expectedDbEnv := map[string]string{
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
	}
	if !reflect.DeepEqual(dbService.ParsedEnvironment, expectedDbEnv) {
		t.Errorf("Expected db environment %+v, got %+v", expectedDbEnv, dbService.ParsedEnvironment)
	}
	expectedDbVolumes := []map[string]string{
		{"type": "bind", "source": "/data/db", "target": "/var/lib/postgresql/data"},
	}
	if !reflect.DeepEqual(dbService.ParsedVolumes, expectedDbVolumes) {
		t.Errorf("Expected db volumes %+v, got %+v", expectedDbVolumes, dbService.ParsedVolumes)
	}

	// Test volumes
	if len(volumes) != 1 {
		t.Errorf("Expected 1 volume, got %d", len(volumes))
	}
	if _, ok := volumes["data"]; !ok {
		t.Error("Expected 'data' volume, but not found")
	}

	// Test networks
	if len(networks) != 1 {
		t.Errorf("Expected 1 network, got %d", len(networks))
	}
	if _, ok := networks["default"]; !ok {
		t.Error("Expected 'default' network, but not found")
	}
}

func TestServiceUnmarshalYAML_Environment(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected map[string]string
	}{
		{
			name: "Sequence of key=value strings",
			yaml: `environment:
      - VAR1=value1
      - VAR2=value2`,
			expected: map[string]string{"VAR1": "value1", "VAR2": "value2"},
		},
		{
			name: "Mapping of key-value pairs",
			yaml: `environment:
      VAR3: value3
      VAR4: value4`,
			expected: map[string]string{"VAR3": "value3", "VAR4": "value4"},
		},
		{
			name: "Mixed (should parse as sequence if first element is sequence)",
			yaml: `environment:
      - VAR5=value5
      VAR6: value6`,
			expected: map[string]string{"VAR5": "value5"},
		},
		{
			name:     "Empty environment",
			yaml:     `environment: {}`,
			expected: map[string]string{},
		},
		{
			name: "Environment with no value",
			yaml: `environment:
      - VAR_NO_VALUE`,
			expected: map[string]string{"VAR_NO_VALUE": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s struct {
				Environment       yaml.Node `yaml:"environment"`
				ParsedEnvironment map[string]string
			}

			// Manually unmarshal to the anonymous struct to get the yaml.Node
			if err := yaml.Unmarshal([]byte(tt.yaml), &s); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Now, create a Service and manually call its UnmarshalYAML with the extracted node
			var service Service
			service.Environment = s.Environment // Assign the parsed yaml.Node
			if err := service.UnmarshalYAML(&yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Value: "environment"}, &s.Environment,
			}}); err != nil {
				t.Fatalf("Service UnmarshalYAML failed: %v", err)
			}

			if !reflect.DeepEqual(service.ParsedEnvironment, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, service.ParsedEnvironment)
			}
		})
	}
}

func TestServiceUnmarshalYAML_Volumes(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected []map[string]string
	}{
		{
			name: "Short syntax volumes",
			yaml: `volumes:
      - /host/path:/container/path
      - named_volume:/app/data`,
			expected: []map[string]string{
				{"source": "/host/path", "target": "/container/path"},
				{"source": "named_volume", "target": "/app/data"},
			},
		},
		{
			name: "Long syntax volumes (bind)",
			yaml: `volumes:
      - type: bind
        source: /data/src
        target: /app/src`,
			expected: []map[string]string{
				{"type": "bind", "source": "/data/src", "target": "/app/src"},
			},
		},
		{
			name: "Mixed syntax volumes",
			yaml: `volumes:
      - /host/path:/container/path
      - type: bind
        source: /data/src
        target: /app/src`,
			expected: []map[string]string{
				{"source": "/host/path", "target": "/container/path"},
				{"type": "bind", "source": "/data/src", "target": "/app/src"},
			},
		},
		{
			name:     "Empty volumes",
			yaml:     `volumes: []`,
			expected: []map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s struct {
				Volumes       yaml.Node `yaml:"volumes"`
				ParsedVolumes []map[string]string
			}

			// Manually unmarshal to the anonymous struct to get the yaml.Node
			if err := yaml.Unmarshal([]byte(tt.yaml), &s); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Now, create a Service and manually call its UnmarshalYAML with the extracted node
			var service Service
			service.Volumes = s.Volumes // Assign the parsed yaml.Node
			if err := service.UnmarshalYAML(&yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
				{Value: "volumes"}, &s.Volumes,
			}}); err != nil {
				t.Fatalf("Service UnmarshalYAML failed: %v", err)
			}

			if !reflect.DeepEqual(service.ParsedVolumes, tt.expected) {
				t.Errorf("Expected %+v, got %+v", tt.expected, service.ParsedVolumes)
			}
		})
	}
}
