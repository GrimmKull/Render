# Render

Expansion of the Martini middleware/handler for easily rendering serialized JSON, XML, Markdown, HTML and Amber template responses.

You can find the original Render middleware [here](https://github.com/codegangsta/martini-contrib/tree/master/render)

This implementation uses:

* [Amber](https://github.com/eknkc/amber) for templates
* [Blackfriday](https://github.com/russross/blackfriday) for Markdown

Make sure that the markdown files are stored in a markdown folder.

## Example

### Main application file

```go
package main

import (
	"github.com/GrimmKull/render"
	"github.com/codegangsta/martini"
)

type Message struct {
	Name string
	Body string
	Time int64
}

func main() {
	msg := Message{"World", "Hello", 1294706395881547000}

	m := martini.Classic()

	m.Use(render.Renderer("templates"))


	m.Get("/text", func(r render.Render) {
		r.TEXT(200, msg)
	})

	m.Get("/json", func(r render.Render) {
		r.JSON(200, msg)
	})

	m.Get("/xml", func(r render.Render) {
		r.XML(200, msg)
	})

	m.Get("/html", func(r render.Render) {
		r.HTML(200, "hello", "world")
	})

	m.Get("/amber", func(r render.Render) {
		r.HTML(200, "index", msg)
	})

	m.Get("/md", func(r render.Render) {
		r.MD(200, "markdown", msg)
	})

	m.Run()
}

```

### Templates

Go template `hello.tmpl`

```html
<!-- templates/hello.tmpl -->
<h2>Hello {{.}}!</h2>
```

Amber template `index.amber`

```amber
//-templates/index.amber
!!! 5

html
	head
		title Hello #{Name}

		meta[name="description"][value="This is a sample"]

		script[type="text/javascript"]
			var hw = "Hello World!"
			alert(hw)

		style[type="text/css"]
			body {
				background: maroon;
				color: white;
			}

	body
		header#mainHeader
			ul
				li.active
					a[href="/"] Main Page #{Time}
						[title="Main Page"]
		div #{Body}
		footer
			| Hey
			br
			| there
```
