package generators

import (
	"fmt"
	"strings"
)

// Generate NS8 Volume flags eg --volume
func GenerateNS8VolumeFlags(volumes []string) string {
	formattedVolume := ""
	for _, volume := range volumes {
		formattedVolume += fmt.Sprintf(" --volume %s", volume)
	}
	return strings.TrimSpace(formattedVolume)
}

// GenerateNS8AfterServices by using docker compose depends_on
func GenerateNS8AfterServices(services interface{}, allServices, mainService string) string {
	var dependsOn string = mainService + " "
	// Check the Depends On type since they can be two types []string or map[string]map[string]string
	switch service := services.(type) {
	case []interface{}:
		for _, s := range service {
			dependsOn += s.(string) + "-app.service "
		}
	case map[string]interface{}:
		for name := range service {
			dependsOn += name + "-app.service "
		}
	default:
		fmt.Printf("Unknown type for this: %T \n", service)
	}
	// return a string
	return strings.TrimSpace(dependsOn)
}
