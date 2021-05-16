package internal

type Parameter interface {
	GetValue() string
	AddPlaceholders([]Placeholder)
}

type SimpleParameter struct {
	Value string
}

func (p *SimpleParameter) GetValue() string {
	return p.Value
}

func (p *SimpleParameter) AddPlaceholders(ph []Placeholder) {
	return
}

type ParameterWithPlaceholders struct {
	Value        string
	Placeholders []Placeholder
}

func (pwp *ParameterWithPlaceholders) GetValue() string {
	outString := pwp.Value
	for _, placeholder := range pwp.Placeholders {
		outString = placeholder.Replace(outString)
	}
	return outString
}
func (pwp *ParameterWithPlaceholders) AddPlaceholders(ph []Placeholder) {
	pwp.Placeholders = append(pwp.Placeholders, ph...)
	return
}
