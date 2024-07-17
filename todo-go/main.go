package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Installer struct {
	Installer string `json:"installer"`
	Id        string `json:"id"`
}

type App struct {
	Name          string      `json:"name"`
	Windows       interface{} `json:"windows"`
	Linux         interface{} `json:"linux"`
	Prerequisites []string    `json:"prerequisites"`
}

type InstallApp struct {
	Name      string `json:"name"`
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
}

type PrettyInstallApp struct {
	title       string
	installed   bool
	description string
}

// Implement the list.Item interface
func (app PrettyInstallApp) FilterValue() string {
	return app.title
}

func (app PrettyInstallApp) Title() string {
	return app.title
}

func (app PrettyInstallApp) Description() string {
	return app.description
}

type Model struct {
	list list.Model
}

func New() *Model {
	return &Model{}
}

func (model *Model) initList(width int, height int) {
	model.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	model.list.Title = "Apps"

	// Convert the array of InstallApp to a slice of list.Item
	// items := make([]list.Item, len(installApps))
	// for i, app := range installApps {
	// 	items[i] = PrettyInstallApp{
	// 		name:      app.Name,
	// 		installed: app.Installed,
	// 		version:   app.Version,
	// 	}
	// }

	model.list.SetItems([]list.Item{
		PrettyInstallApp{
			title:       "GoLang",
			description: "none",
			installed:   false,
		},
	})
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.initList(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	model.list, cmd = model.list.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	return model.list.View()
}

func isInstallApp(installApps []InstallApp, appName string) bool {
	for _, installApp := range installApps {
		if installApp.Name == appName {
			return true
		}
	}
	return false
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
		file.WriteString("[]")
		file.Close()
		fmt.Println("File 'nk-apps.json' created:", filePath)
	} else if err != nil {
		fmt.Println("Error checking file 'nk-apps.json' existence:", err)
		return
	} else {
		fmt.Println("File 'nk-apps.json' already exists:", filePath)
	}

	// Read the existing file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file 'nk-apps.json':", err)
		return
	}

	// Parse the JSON content into a slice of InstallApp structs
	var installApps []InstallApp
	err = json.Unmarshal(data, &installApps)
	if err != nil {
		fmt.Println("Error unmarshalling 'nk-apps.json':", err)
		return
	}

	for _, app := range apps {
		if !isInstallApp(installApps, app.Name) {
			newInstallApp := InstallApp{
				Name:      app.Name,
				Installed: false,
				Version:   "none",
			}
			installApps = append(installApps, newInstallApp)
		}
	}

	model := New()
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Println("Program error:", err)
		os.Exit(1)
	}

	// Marshal the modified data into JSON with indentation
	newData, err := json.MarshalIndent(installApps, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Overwrite the file with the new content
	err = os.WriteFile(filePath, newData, 0644)
	if err != nil {
		fmt.Println("Error writing to file 'nk-apps.json':", err)
		return
	}

	fmt.Println("'nk-apps.json' Updated!")
}
