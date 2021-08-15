package gweb

type IWidget interface {
	Render(ctx *Context) (interface{}, error)
}
