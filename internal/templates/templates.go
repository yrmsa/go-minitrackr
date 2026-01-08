package templates

import (
	"embed"
	"html/template"
	"sync"
)

//go:embed *.html
var templateFS embed.FS

var (
	tmpl *template.Template
	once sync.Once
)

func Load() (*template.Template, error) {
	var err error
	once.Do(func() {
		tmpl, err = template.ParseFS(templateFS, "*.html")
	})
	return tmpl, err
}
