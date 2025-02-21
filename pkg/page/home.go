package page

func NewHomePage() *Page {
	p := New()
	p.Title = "Home"
	p.Name = "home"
	p.Path = "/"
	return p
}
