package formatters

import (
	"ns8-module-generator/parser"
	"strings"
)

var (
	HostDefault = "docker.io/"
	ImageSuffix = "_IMAGE"
)

// GetImagesWithRepository Get Images with repository details and tag
func GetImagesWithRepository() []string {
	images := parser.GetImages()
	// Check if  image has hosting eg, docker.io, ghcr.io etc
	var imagesWithRepository []string
	for _, image := range images {
		imageParts := strings.Split(image, "/")
		// Check if the image has a repository
		if strings.Contains(imageParts[0], ".") {
			// Image has a repository (eg, docker.io, ghcr.io)
			imagesWithRepository = append(imagesWithRepository, image)
		} else {
			// Image does not have a repository
			imagesWithRepository = append(imagesWithRepository, HostDefault+image)
		}
	}
	return imagesWithRepository
}

// FormatImageNames Convert the Images into a UPPERCASE string and return it as array
func FormatImageNames() []string {
	var formattedImages []string
	images := parser.GetImages()
	for _, image := range images {
		parts := strings.Split(image, ":")
		imageName := parts[0]

		// Split by "/" and get the last element to handle the case where the image is in a repository
		nameParts := strings.Split(imageName, "/")
		var formattedName string
		if len(nameParts) > 1 {
			formattedName = strings.ToUpper(nameParts[len(nameParts)-1]) + ImageSuffix
		} else {
			formattedName = strings.ToUpper(imageName) + ImageSuffix
		}
		formattedImages = append(formattedImages, formattedName)
	}
	return formattedImages
}

func stripImageName(imageName string) string {
	// Split Parts by Colon
	parts := strings.Split(imageName, ":")
	// Get Image name
	nameParts := strings.Split(parts[0], "/")
	var name string
	// Check length of nameParts
	if len(nameParts) > 1 {
		name = strings.ToLower(nameParts[len(nameParts)-1])
	} else {
		name = strings.ToLower(parts[0])
	}
	return name
}

func GetImagesCompatibleServiceNames() []string {
	var serviceNames []string
	images := parser.GetImages()
	for _, image := range images {
		// Format Image name to lower case
		imageName := stripImageName(image)
		// Append to serviceNames
		serviceNames = append(serviceNames, imageName)
	}
	return serviceNames
}
