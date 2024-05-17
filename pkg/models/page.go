package models

type Page struct {
	Data TaskList
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

func NewPage(tasklist TaskList) Page {
	return Page{
		Data: tasklist,
		Form: NewFormData(),
	}
}
