package gweb

import (
	"html/template"
)

type IWidget interface {
	Render(ctx *Context) template.HTML
}
