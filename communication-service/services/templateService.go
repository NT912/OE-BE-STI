package services

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

func RenderWelcomeTemplate(name string) (string, error) {
	return RenderTemplate("welcome.html", map[string]interface{}{"Name": name})
}

func RenderTemplate(templateFile string, data map[string]interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("./templates/%s", templateFile))
	if err != nil {
		return "", fmt.Errorf("could not find absolute path for template: %w", err)
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("could not parse template file: %w", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return "", fmt.Errorf("could not execute template: %w", err)
	}

	return body.String(), nil
}
