package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// App struct to hold each app's data
type App struct {
	Name          string   `json:"name"`
	WindowInstall string   `json:"window_install"`
	LinuxInstall  string   `json:"linux_install"`
	Prerequisites []string `json:"prerequisites"`
}

func main() {
	// Open the JSON file
	appsData, err := os.ReadFile("../apps.json")
	if err != nil {
		fmt.Println("Error opening file apps.json:", err)
		return
	}

	// Decode JSON array into a slice of App structs
	var apps []App
	err = json.Unmarshal(appsData, &apps)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	folderPath := filepath.Join(homeDir, ".config")
	// Create the folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755) // 0755 permission gives read and write access to owner, and read access to group and others
		if err != nil {
			fmt.Println("Error creating folder '.config':", err)
			return
		}
		fmt.Println("Folder '.config' created:", folderPath)
	} else if err != nil {
		fmt.Println("Error checking folder '.config' existence:", err)
		return
	} else {
		fmt.Println("Folder '.config' already exists:", folderPath)
	}

	filePath := filepath.Join(folderPath, "nk-apps.json")
	// Create the file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}

		// file.WriteString()

		file.Close()
		fmt.Println("File 'nk-apps.json' created:", filePath)
	} else if err != nil {
		fmt.Println("Error checking file 'nk-apps.json' existence:", err)
		return
	} else {
		fmt.Println("File 'nk-apps.json' already exists:", filePath)
	}
}
