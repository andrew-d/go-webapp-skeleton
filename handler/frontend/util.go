package frontend

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
	"golang.org/x/net/context"

	"github.com/andrew-d/go-webapp-skeleton/handler/frontend/layouts"
	"github.com/andrew-d/go-webapp-skeleton/handler/frontend/templates"
	"github.com/andrew-d/go-webapp-skeleton/log"
)

type M map[string]interface{}

var (
	// Buffer pool for rendering templates
	bufpool *bpool.BufferPool

	// Map of templates
	templatesMap map[string]*template.Template

	// Extra functions
	templateFuncs = template.FuncMap{}
)

func init() {
	// Create buffer pool
	bufpool = bpool.NewBufferPool(64)

	// Get the contents of all layouts.
	layoutData := make(map[string]string)
	for _, lname := range layouts.AssetNames() {
		d, _ := layouts.Asset(lname)
		layoutData[lname] = string(d)
	}

	// For each template, we parse it.
	templatesMap = make(map[string]*template.Template)
	for _, aname := range templates.AssetNames() {
		tname := filepath.Base(aname)

		// Create new template with functions
		tmpl := template.New(tname).Funcs(templateFuncs)

		// Get the template's data
		d, _ := templates.Asset(aname)

		// Parse the main template, then all the layouts.
		tmpl = template.Must(tmpl.Parse(string(d)))
		for _, layout := range layouts.AssetNames() {
			tmpl = template.Must(tmpl.Parse(layoutData[layout]))
		}

		// Insert
		templatesMap[tname] = tmpl
	}
}

// renderTemplate is a wrapper around template.ExecuteTemplate.  It writes into
// a bytes.Buffer before writing to the http.ResponseWriter to catch any errors
// resulting from populating the template.
func renderTemplate(ctx context.Context, w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templatesMap[name]
	if !ok {
		log.FromContext(ctx).WithField("name", name).Error("template does not exist")
		return fmt.Errorf("The template %s does not exist", name)
	}

	// Create a buffer to temporarily write to and check if any errors were encounted.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		log.FromContext(ctx).WithField("err", err).Error("could not render template")
		return err
	}

	// Set the header and write the buffer to the http.ResponseWriter
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
