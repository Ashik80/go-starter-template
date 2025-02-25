package infrastructure

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

type TemplateRenderer interface {
	GetTemplate(key string) (*template.Template, error)
	Render(w http.ResponseWriter, p *page.Page) error
	RenderPartial(w http.ResponseWriter, partial string, data any) error
	RenderString(w http.ResponseWriter, html string, data any) error
}

type TemplateRendererImpl struct {
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
	pageName := strings.Join(parts[3:], ":")

	return fmt.Sprintf("%s:%s", layoutName, pageName)
}

// NewTemplateRenderer returns a new TemplateRenderer
// that can be used to render templates
// to a http.ResponseWriter
//
// Parameters:
//   - baseTemplateFile: the path to the base template file
//   - layoutDir: the path to the layout directory
//   - pagesDir: the path to the pages directory
//   - partialsDir: the path to the partials directory
//
// Returns:
//   - *TemplateRenderer: a new TemplateRenderer
//   - error: an error
func NewTemplateRenderer(baseTemplateFile, layoutDir, pagesDir, partialsDir string) (TemplateRenderer, error) {
	templates := make(map[string]*template.Template)

	layouts, err := walkTemplateFiles(layoutDir)
	if err != nil {
		errorMsg := "failed to parse layouts"
		log.Printf("ERROR: %s: %v\n", errorMsg, err)
		return nil, fmt.Errorf("%s: %w", errorMsg, err)
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
		return nil, fmt.Errorf("%s: %w", errorMsg, err)
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
		return nil, fmt.Errorf("%s: %w", errorMsg, err)
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
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	for _, layoutFile := range layouts {
		for _, pageFile := range pages {
			key := getTemplateKey(layoutFile, pageFile)
			files := append([]string{baseTemplateFile, layoutFile, pageFile}, partials...)
			tmpl, err := template.ParseFiles(files...)
			if err != nil {
				errorMsg := fmt.Sprintf("failed to parse template %s", key)
				log.Printf("ERROR: %s: %v\n", errorMsg, err)
				return nil, fmt.Errorf("%s: %w", errorMsg, err)
			}
			templates[key] = tmpl
		}
	}

	return &TemplateRendererImpl{
		templates: templates,
		partials:  partialTmpl,
	}, nil
}

func (t *TemplateRendererImpl) Render(w http.ResponseWriter, p *page.Page) error {
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

	w.Header().Add("Content-Type", "text/html")

	key := fmt.Sprintf("%s:%s", p.Layout, p.Name)
	tmpl, ok := t.templates[key]
	if !ok {
		errorMsg := fmt.Sprintf("template not found for key: %s", key)
		log.Printf("ERROR: %s", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}

	return tmpl.ExecuteTemplate(w, "base", p)
}

func (t *TemplateRendererImpl) RenderPartial(w http.ResponseWriter, partial string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	return t.partials.ExecuteTemplate(w, partial, data)
}

func (t *TemplateRendererImpl) RenderString(w http.ResponseWriter, html string, data any) error {
	w.Header().Add("Content-Type", "text/html")
	tmpl, err := template.New("").Parse(html)
	if err != nil {
		errorMsg := "failed to parse string"
		log.Printf("ERROR: %s\n", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}
	return tmpl.Execute(w, data)
}

func (t *TemplateRendererImpl) GetTemplate(key string) (*template.Template, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	temp, ok := t.templates[key]
	if !ok {
		errorMsg := fmt.Sprintf("template not found for key: %s", key)
		log.Printf("ERROR: %s\n", errorMsg)
		return nil, fmt.Errorf("%s", errorMsg)
	}

	return temp, nil
}
