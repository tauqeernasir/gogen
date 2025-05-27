package templates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Manager handles template loading and management
type Manager struct {
	templates map[string]*template.Template
}

// NewManager creates a new template manager
func NewManager() *Manager {
	return &Manager{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplates loads templates from a directory or uses embedded templates
func (tm *Manager) LoadTemplates(templatesDir string) error {
	// Load TypeScript templates
	if err := tm.loadLanguageTemplates("typescript", templatesDir); err != nil {
		return fmt.Errorf("failed to load TypeScript templates: %w", err)
	}

	// Load Python templates
	if err := tm.loadLanguageTemplates("python", templatesDir); err != nil {
		return fmt.Errorf("failed to load Python templates: %w", err)
	}

	return nil
}

// loadLanguageTemplates loads templates for a specific language
func (tm *Manager) loadLanguageTemplates(language, templatesDir string) error {
	langDir := filepath.Join(templatesDir, language)
	if _, err := os.Stat(langDir); os.IsNotExist(err) {
		// If templates directory doesn't exist, use embedded templates
		return tm.loadEmbeddedTemplates(language)
	}

	return filepath.Walk(langDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(templatesDir, path)
			templateName := strings.TrimSuffix(relPath, ".tmpl")

			tmpl, err := template.New(templateName).Parse(string(content))
			if err != nil {
				return err
			}

			tm.templates[templateName] = tmpl
		}

		return nil
	})
}

// loadEmbeddedTemplates loads embedded templates for a language
func (tm *Manager) loadEmbeddedTemplates(language string) error {
	switch language {
	case "typescript":
		return tm.loadTypeScriptTemplates()
	}
	return nil
}

// GetTemplate retrieves a template by name
func (tm *Manager) GetTemplate(name string) (*template.Template, bool) {
	tmpl, exists := tm.templates[name]
	return tmpl, exists
}
