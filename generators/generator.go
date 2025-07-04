package generators

import (
	"fmt"
	"ns8-module-generator/config"
	"strings"
)

// Generate NS8 Volume flags eg --volume
func GenerateNS8VolumeFlags(volumes []map[string]string) string {
	formattedVolume := ""
	for _, volumeMap := range volumes {
		source := volumeMap["source"]
		target := volumeMap["target"]
		typeOfVolume := volumeMap["type"]

		volumeString := ""
		if typeOfVolume == "bind" {
			volumeString = fmt.Sprintf("%s:%s", source, target)
		} else {
			// Default to bind mount if type is not specified or unknown
			volumeString = fmt.Sprintf("%s:%s", source, target)
		}

		formattedVolume += fmt.Sprintf(" --volume %s", volumeString)

		// For backup, we still need the original volume name (source part)
		volumeName := source
		// Check volume name prefix
		if !strings.HasPrefix(volumeName, "./") && !strings.HasPrefix(volumeName, "/") {
			err := AddToBackup(
				config.Cfg.OutputDir+"/imageroot/etc/state-include.conf",
				fmt.Sprintf("volumes/%s", volumeName),
			)
			if err != nil {
				fmt.Printf("An error occurred while adding volume back up: %v", err)
			}
		}
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


