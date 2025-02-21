package page

type Page struct {
	Name   string
	Title  string
	Layout string
	Path   string
	Data   any
	Error  string
}

func New() *Page {
	return &Page{}
}
