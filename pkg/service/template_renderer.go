package service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go-starter-template/pkg/page"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
	partials  *template.Template
	mu        sync.Mutex
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

func getTemplateKey(layoutPath, pagePath string) string {
	layoutName := strings.TrimSuffix(filepath.Base(layoutPath), ".html")

	relPath := strings.TrimSuffix(pagePath, ".html")
	parts := strings.Split(relPath, string(os.PathSeparator))
	pageName := strings.Join(parts[2:], ":")

	return fmt.Sprintf("%s:%s", layoutName, pageName)
}

func NewTemplateRenderer(layoutDir, pagesDir, partialsDir string) (*TemplateRenderer, error) {
	templates := make(map[string]*template.Template)

	layouts, err := walkTemplateFiles(layoutDir)
	if err != nil {
		errorMsg := "failed to parse layouts"
		log.Printf("ERROR: %s: %v\n", errorMsg, err)
		return nil, fmt.Errorf("%s: %v", errorMsg, err)
	}
	if len(layouts) == 0 {
		errorMsg := fmt.Sprintf("no layout files found in %s", layoutDir)
		log.Printf("ERROR: %s\n", errorMsg)
		return nil, fmt.Errorf("%s", errorMsg)
	}

	pages, err := walkTemplateFiles(pagesDir)
	if err != nil {
		errorMsg := "failed to parse pages"
		log.Printf("ERROR: %s: %v\n", errorMsg, err)
		return nil, fmt.Errorf("%s: %v", errorMsg, err)
	}
	if len(pages) == 0 {
		errorMsg := fmt.Sprintf("no page files found in %s", pagesDir)
		log.Printf("ERROR: %s\n", errorMsg)
		return nil, fmt.Errorf("%s", errorMsg)
	}

	partials, err := walkTemplateFiles(partialsDir)
	if err != nil {
		errorMsg := "failed to parse partials"
		log.Printf("ERROR: %s: %v\n", errorMsg, err)
		return nil, fmt.Errorf("%s: %v", errorMsg, err)
	}
	if len(partials) == 0 {
		errorMsg := fmt.Sprintf("no partial files found in %s", partialsDir)
		log.Printf("ERROR: %s\n", errorMsg)
		return nil, fmt.Errorf("%s", errorMsg)
	}

	partialTmpl, err := template.ParseFiles(partials...)
	if err != nil {
		errMsg := "error parsing partials directory"
		log.Printf("ERROR: %s: %v\n", errMsg, err)
		return nil, fmt.Errorf("%s: %v", errMsg, err)
	}

	for _, layoutFile := range layouts {
		for _, pageFile := range pages {
			key := getTemplateKey(layoutFile, pageFile)
			files := append([]string{layoutFile, pageFile}, partials...)
			tmpl, err := template.ParseFiles(files...)
			if err != nil {
				errorMsg := fmt.Sprintf("failed to parse template %s", key)
				log.Printf("ERROR: %s: %v\n", errorMsg, err)
				return nil, fmt.Errorf("%s: %v", errorMsg, err)
			}
			templates[key] = tmpl
		}
	}

	return &TemplateRenderer{
		templates: templates,
		partials:  partialTmpl,
	}, nil
}

func (t *TemplateRenderer) Render(w http.ResponseWriter, p *page.Page) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if p.Name == "" {
		errorMsg := "template name is required"
		log.Printf("ERROR: %s\n", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}

	if p.Layout == "" {
		p.Layout = "main"
	}

	if p.StatusCode == 0 {
		p.StatusCode = http.StatusOK
	}

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(p.StatusCode)

	key := fmt.Sprintf("%s:pages:%s", p.Layout, p.Name)
	tmpl, ok := t.templates[key]
	if !ok {
		errorMsg := fmt.Sprintf("template not found for key: %s", key)
		log.Printf("ERROR: %s\n", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}

	templateName := filepath.Base(p.Layout)
	return tmpl.ExecuteTemplate(w, templateName, p)
}

func (t *TemplateRenderer) RenderPartial(w http.ResponseWriter, statusCode int, partial string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	return t.partials.ExecuteTemplate(w, partial, data)
}

func (t *TemplateRenderer) RenderString(w http.ResponseWriter, statusCode int, html string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	tmpl, err := template.New("").Parse(html)
	if err != nil {
		errorMsg := "failed to parse string"
		log.Printf("ERROR: %s\n", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}
	return tmpl.Execute(w, data)
}
