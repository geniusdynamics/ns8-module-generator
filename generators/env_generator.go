package generators

import (
	"fmt"
	"ns8-module-generator/processors"
	"os"
	"strings"
)

// Handle env from the envrioment service
func GenerateEnvFileContents(
	imageName string,
	enviroments []string,
	filePath string,
) (string, error) {
	println("Image Name Enviroment: " + imageName + " filePath: " + filePath)
	// Write to configure vars all variables as per now
	var vars strings.Builder
	envConfig := fmt.Sprintf("%s = { \n", imageName)
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
		return "", fmt.Errorf("Failed to read the file: %v", err)
	}
	fmt.Println("Existing File Contents: ", string(content))

	// Open file in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write to file: %v ", err)
	}
	// Close file later
	defer file.Close()
	// Write the new content to file
	if _, err := file.WriteString(vars.String() + " \n"); err != nil {
		return "", fmt.Errorf("Failed to write to file: %v", err)
	}
	// Write env config
	if _, err := file.WriteString(envConfig + " \n"); err != nil {
		return "", fmt.Errorf("Failed to write env file config: %v", err)
	}

	println("New file content added")
	envFileFlags := " --env " + imageName + ".env"

	// Add ENV to back up
	err = AddToBackup(
		processors.OutputDir+"/imageroot/etc/state-include.conf",
		fmt.Sprintf("state/%s.env \n", imageName),
	)
	if err != nil {
		return "", fmt.Errorf(
			"An error occurred while adding %s.env to back up: %v",
			imageName,
			err,
		)
	}

	// return nil
	return envFileFlags, nil
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

func GenerateGetConfigurationContent(
	imageName string,
	enviroments []string,
	filePath string,
) error {
	fmt.Println("Generating Get Configurations")
	var config strings.Builder
	// string to check if env exists
	line := fmt.Sprintf("if os.path.exists(\"%s.env\"): \n", imageName)
	config.WriteString(line)
	// Read the env file
	config.WriteString(fmt.Sprintf("\tdata = agent.read_envfile(\"%s.env\") \n", imageName))
	// Loop thru the env vars and put in config obj
	for _, env := range enviroments {
		// clean the enviroment string
		cleanEnv := cleanEnviromentString(env)
		// Write string
		config.WriteString(
			fmt.Sprintf("\tconfig[\"%s\"] = data.get(\"%s\") \n", cleanEnv, cleanEnv),
		)
	}
	// Write the else part to return an empty string
	config.WriteString("else: \n")

	// Loop thru the enviroments
	for _, env := range enviroments {
		// Clean env
		cleanEnv := cleanEnviromentString(env)
		// Add empty string
		config.WriteString(fmt.Sprintf("\tconfig[\"%s\"] = \"\" \n", cleanEnv))
	}
	// Read get the configuration file
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
	if _, err := file.WriteString(config.String() + " \n"); err != nil {
		return fmt.Errorf("Failed to write to file: %v", err)
	}
	return nil
}

func AddJsonDump(filePath string) error {
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
	jsonDump := "json.dump(config, fp=sys.stdout)"
	if _, err := file.WriteString(jsonDump + "\n"); err != nil {
		return fmt.Errorf("Failed to add JSON DUMP in %s;", filePath)
	}
	return nil
}
