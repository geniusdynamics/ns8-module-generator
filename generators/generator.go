package generators

import (
	"fmt"
	"strings"
)

func GenerateNS8VolumeFlags(volumes []string) string {
	formattedVolume := ""
	for _, volume := range volumes {
		formattedVolume += fmt.Sprintf(" --volume %s", volume)
	}
	return strings.TrimSpace(formattedVolume)
}
