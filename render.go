package marmoset

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"runtime"
)

// Renderer ...
type Renderer interface {
	JSON(int, interface{}) error
	HTML(string, interface{}) error
}

var templates *template.Template

// Render ...
func Render(w http.ResponseWriter, pretty ...bool) Renderer {
	if len(pretty) != 0 && pretty[0] {
		return PrettyRenderer{
			w: w,
		}
	}
	return PrettyRenderer{
		w: w,
	}
}

// PrettyRenderer ...
type PrettyRenderer struct {
	w http.ResponseWriter
}

// JSON ...
func (r PrettyRenderer) JSON(status int, data interface{}) error {
	r.w.WriteHeader(status)
	r.w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if _, err := r.w.Write(b); err != nil {
		return err
	}
	return nil
}

// RenderJSON ...
func RenderJSON(w http.ResponseWriter, status int, data interface{}) error {
	return Render(w, true).JSON(status, data)
}

// HTML ...
func (r PrettyRenderer) HTML(name string, data interface{}) error {
	if templates == nil {
		return fmt.Errorf("templates not loaded")
	}
	return templates.ExecuteTemplate(r.w, name+".html", data)
}

// LoadViews ...
func LoadViews(p string) error {
	_, f, _, _ := runtime.Caller(1)
	tpl, err := template.ParseGlob(path.Join(path.Dir(f), p, "*.html"))
	if err != nil {
		panic(err)
	}
	templates = tpl
	return nil
}