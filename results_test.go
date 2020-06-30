package gweb

import "testing"

func TestErrorResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *ErrorResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestNotFindResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *ViewActionMappingResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestViewResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *ViewResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestEmptyResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *EmptyResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestHTMLResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *HTMLResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestJsonResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *JsonResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestTextResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *TextResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestXMLResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *XMLResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestRedirectToUrlResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *RedirectToUrlResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestImageResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *ImageResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}

func TestImageBytesResult_Apply(t *testing.T) {
	type args struct {
		context *Context
	}
	tests := []struct {
		name string
		r    *ImageBytesResult
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Apply(tt.args.context)
		})
	}
}
