package page

type Page struct {
	Name       string
	Title      string
	Layout     string
	StatusCode int
	Data       any
	Error      string
}

func New() *Page {
	return &Page{}
}
