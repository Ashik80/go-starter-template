package renderer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var baseTemplate *template.Template

func InitBaseTemplate() error {
	baseFile := "web/templates/base.html"
	partials, err := getPartialFiles()
	if err != nil {
		return err
	}
	files := append([]string{baseFile}, partials...)
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}
	baseTemplate = tmpl
	return nil
}

func ParseTemplate(page string, layout ...string) *template.Template {
	layoutFile := "web/templates/layouts/main.html"
	if len(layout) > 0 && layout[0] != "" {
		layoutFile = "web/templates/layouts/" + layout[0] + ".html"
	}
	file := "web/templates/pages/" + page + ".html"

	tmpl, err := baseTemplate.Clone()
	if err != nil {
		log.Fatalf("ERROR: failed to clone base template: %v", err)
	}

	tmpl, err = tmpl.ParseFiles(layoutFile, file)
	if err != nil {
		log.Fatalf("ERROR: failed to parse template: %v", err)
	}

	return tmpl
}

func GetBaseTemplate() *template.Template {
	return baseTemplate
}

func RenderString(w http.ResponseWriter, html string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("").Parse(html)
	if err != nil {
		errorMsg := "failed to parse string"
		log.Printf("ERROR: %s\n", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}
	return tmpl.Execute(w, data)
}

func walkTemplateFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func getPartialFiles() ([]string, error) {
	partialDir := "web/templates/partials"
	files, err := walkTemplateFiles(partialDir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no partial files found in %s", partialDir)
	}
	return files, nil
}
