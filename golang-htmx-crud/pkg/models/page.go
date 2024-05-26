package models

type Page struct {
	Data Tasks
	Form FormData
}

func (p *Page) NewFormData() FormData {
	return NewFormData()
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

func NewPage(tasklist Tasks) Page {
	return Page{
		Data: tasklist,
		Form: NewFormData(),
	}
}
