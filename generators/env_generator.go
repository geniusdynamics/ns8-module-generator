package generators

import (
	"fmt"
	"os"
	"strings"
)

// Handle env from the envrioment service
func GenerateEnvFileContents(imageName string, enviroments []string, filePath string) error {
	println("Image Name Enviroment: " + imageName + " filePath: " + filePath)
	// Write to configure vars all variables as per now
	var vars strings.Builder
	envConfig := imageName + " = { \n"
	// Loop thru the enviroments
	for _, env := range enviroments {
		cleanEnv := cleanEnviromentString(env)
		line := fmt.Sprintf("%s = data.get(\"%s\") \n", cleanEnv, cleanEnv)
		vars.WriteString(line)
		envConfig += fmt.Sprintf(" \"%s\" : %s ,\n", cleanEnv, cleanEnv)
	}
	// Close the env config
	envConfig += "} \n"

	// write to env file
	envConfig += fmt.Sprintf("agent.write_envfile(\"%s.env\", %s)", imageName, imageName)
	// Check if file Exists
	content, err := os.ReadFile(filePath)
	// If error occurs Close
	if err != nil {
		return fmt.Errorf("Failed to read the file: %v", err)
	}
	fmt.Println("Existing File Contents: ", string(content))

	// Open file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to write to file: %v ", err)
	}
	// Close file later
	defer file.Close()
	// Write the new content to file
	if _, err := file.WriteString(vars.String() + " \n"); err != nil {
		return fmt.Errorf("Failed to write to file: %v", err)
	}
	// Write env config
	if _, err := file.WriteString(envConfig + " \n"); err != nil {
		return fmt.Errorf("Failed to write env file config: %v", err)
	}

	println("New file content added")

	// return nil
	return nil
}

// Remove amnything after =
func cleanEnviromentString(env string) string {
	// Split the string after the first =
	parts := strings.SplitN(env, "=", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

func generateNS8EnvFileFlags() string {
	return ""
}
