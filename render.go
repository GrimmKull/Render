package render

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/eknkc/amber"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	ContentType = "Content-Type"
	ContentText = "text/plain; charset=utf-8"
	ContentJSON = "application/json"
	ContentXML  = "application/xml"
	ContentHTML = "text/html"
)

// Render is a service that can be injected into a Martini handler. Render provides functions for easily writing JSON and
// HTML templates out to a http Response.
type Render interface {
	// Text writes raw text response
	TEXT(status int, v interface{})
	// JSON writes the given status and JSON serialized version of the given value to the http.ResponseWriter.
	JSON(status int, v interface{})
	// XML writes the given data as XML to the http.ResponseWriter
	XML(status int, v interface{})
	// HTML renders a html template specified by the name and writes the result and given status to the http.ResponseWriter.
	HTML(status int, name string, v interface{})
	// MD renders markdown files as html
	MD(status int, name string, v interface{})
	// Error is a convenience function that writes an http status to the http.ResponseWriter.
	Error(status int)
}

// Renderer is a Middleware that maps a render.Render service into the Martini handler chain. Renderer will compile templates
// globbed in the given dir. Templates must have the .tmpl or .amber extension to be compiled.
//
// If MARTINI_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MARTINI_ENV environment variable to "production"
func Renderer(dir string) martini.Handler {
	t := compile(dir)
	return func(res http.ResponseWriter, c martini.Context) {
		// recompile for easy development
		if martini.Env == martini.Dev {
			t = compile(dir)
		}
		c.MapTo(&renderer{res, t}, (*Render)(nil))
		//c.MapTo()
	}
}

func compile(dir string) *template.Template {
	t := template.New(dir)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(r)
		if ext == ".tmpl" {

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			name := (r[0 : len(r)-len(ext)])
			tmpl := t.New(filepath.ToSlash(name))
			// Bomb out if parse fails. We don't want any silent server starts.
			template.Must(tmpl.Parse(string(buf)))
		} else if ext == ".amber" {
			fmt.Println("amber")
			compiler := amber.New()

			err := compiler.ParseFile(path)

			if err != nil {
				panic(err)
			}

			name := (r[0 : len(r)-len(ext)])
			tmpl := t.New(filepath.ToSlash(name))

			template.Must(compiler.CompileWithTemplate(tmpl))
		}

		return nil
	})

	return t
}

type renderer struct {
	http.ResponseWriter
	t *template.Template
}

func (r *renderer) TEXT(status int, v interface{}) {
	result := fmt.Sprintln(v)
	//result, ok := v.(string)
	/*if !ok {
		http.Error(r, "Not of type string.", 500)
		return
	}*/

	r.Header().Set(ContentType, ContentText)
	r.WriteHeader(status)
	r.Write([]byte(result))
}

func (r *renderer) JSON(status int, v interface{}) {
	result, err := json.Marshal(v)
	if err != nil {
		http.Error(r, err.Error(), 500)
		return
	}

	// json rendered fine, write out the result
	r.Header().Set(ContentType, ContentJSON)
	r.WriteHeader(status)
	r.Write(result)
}

func (r *renderer) XML(status int, v interface{}) {
	result, err := xml.Marshal(v)

	if err != nil {
		http.Error(r, err.Error(), 500)
		return
	}

	r.Header().Set(ContentType, ContentXML)
	r.WriteHeader(status)
	r.Write(result)
}

func (r *renderer) HTML(status int, name string, binding interface{}) {
	var buf bytes.Buffer
	if err := r.t.ExecuteTemplate(&buf, name, binding); err != nil {
		http.Error(r, err.Error(), 500)
		return
	}

	// template rendered fine, write out the result
	r.Header().Set(ContentType, ContentHTML)
	r.WriteHeader(status)
	r.Write(buf.Bytes())
}

func (r *renderer) MD(status int, name string, binding interface{}) {
	//input, _ := os.Open(name + ".md")
	inputBytes, err := ioutil.ReadFile("markdown/" + name + ".md")

	if err != nil {
		http.Error(r, err.Error(), 500)
		return
	}

	r.Header().Set(ContentType, ContentHTML)
	r.WriteHeader(status)
	r.Write(blackfriday.MarkdownCommon(inputBytes))
}

func (r *renderer) Error(status int) {
	r.WriteHeader(status)
}
