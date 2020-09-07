package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Matt-Gleich/statuser/v2"
	"github.com/Matt-Gleich/texsch/pkg/commands/configure"
	"github.com/Matt-Gleich/texsch/pkg/configuration"
	"github.com/Matt-Gleich/texsch/pkg/location"
	"github.com/Matt-Gleich/texsch/pkg/status"
	"github.com/dustin/go-humanize"
)

type DocumentOutline struct {
	Name  string
	Type  string
	Class string
}

// Create a document
func Document() {
	// Asking information
	classNames := []string{}
	classConfig := configuration.GetClasses()
	for _, classConfiguration := range classConfig {
		classNames = append(classNames, classConfiguration.Name)
	}

	questions := []*survey.Question{
		{
			Name:      "name",
			Prompt:    &survey.Input{Message: "What is the name of the document?"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "What is the type for the document?",
				Options: []string{
					"Paper",
					"Notes",
					"Practice",
					"Assessment",
					"Project",
					"Presentation",
					"Other",
				},
			},
		},
		{
			Name: "class",
			Prompt: &survey.Select{
				Message: "What class is this for?",
				Options: classNames,
			},
		},
	}
	var answers DocumentOutline
	err := survey.Ask(questions, &answers)
	if err != nil {
		statuser.Error("Failed to ask document questions", err, 1)
	}
	path := createFolder(answers)
	createFile(answers, path)
}

// Create the folder for a file
func createFolder(answers DocumentOutline) string {
	path := fmt.Sprintf(
		"LaTeX/%v/%v/%v/",
		strings.ReplaceAll(answers.Class, " ", "-"),
		time.Now().Month(),
		answers.Type,
	)
	err := os.MkdirAll(path, 0700)
	if err != nil {
		statuser.Error("Failed to create folder for file", err, 1)
	}
	return path
}

// Create the actual document file
func createFile(answers DocumentOutline, folderPath string) {
	// Reading from template
	files, err := ioutil.ReadDir("./texsch/templates")
	if err != nil {
		statuser.Error("Failed to get a list of all the templates", err, 1)
	}
	templates := []string{}
	for _, file := range files {
		if !file.IsDir() {
			templates = append(templates, strings.TrimSuffix(file.Name(), ".txt"))
		}
	}
	templateName := "default"
	if len(templates) > 1 {
		prompt := &survey.Select{
			Message: "What template would you like to use?",
			Options: templates,
		}
		err := survey.AskOne(prompt, templateName)
		if err != nil {
			statuser.Error("Failed to ask what template you want to use", err, 1)
		}
	}
	templateContent, err := ioutil.ReadFile("./texsch/templates/" + templateName + ".txt")
	if err != nil {
		statuser.Error("Failed to read from template file", err, 1)
	}

	// Replacing document variables
	var class configure.Class
	for _, classInstance := range configuration.GetClasses() {
		if classInstance.Name == answers.Class {
			class = classInstance
		}
	}
	today := time.Now()
	toReplace := map[string]string{
		"ASSIGNMENT_NAME":  answers.Name,
		"AUTHOR_FULL_NAME": configuration.GetGeneral().Full_Name,
		"CLASS_NAME":       class.Name,
		"CLASS_TIME":       class.Time,
		"CLASS_TEACHER":    class.Teacher_Name,
		"DATE": fmt.Sprintf(
			"%v, %v %v\\textsuperscript{%v}, %v",
			today.Weekday(),
			today.Month(),
			today.Day(),
			strings.TrimLeft(humanize.Ordinal(today.Day()), fmt.Sprint(today.Day())),
			today.Year(),
		),
		"YEAR_NUMBER": fmt.Sprint(today.Year()),
	}
	filledInDocument := string(templateContent)
	for key, value := range toReplace {
		filledInDocument = strings.ReplaceAll(filledInDocument, key, value)
	}

	// Creating the actual file
	filePath := folderPath + strings.ReplaceAll(answers.Name, " ", "-") + ".tex"
	_, err = os.Stat(filePath)
	if !os.IsNotExist(err) {
		var override bool
		prompt := &survey.Confirm{
			Message: "A file with that name already exits. Do you want to override it?",
		}
		err := survey.AskOne(prompt, &override)
		if err != nil {
			statuser.Error("Failed ask if you want to override the file", err, 1)
		}
		if !override {
			os.Exit(0)
		}
	}
	err = ioutil.WriteFile(
		filePath,
		[]byte(filledInDocument),
		0700,
	)
	if err != nil {
		statuser.Error("Failed to write to file", err, 1)
	}
	status.Success("Created file in " + location.GetProjectRoot() + "/" + folderPath)
}