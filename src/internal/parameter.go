package internal

type Parameter interface {
	GetValue() string
}

type SimpleParameter struct {
	Value string
}

func (p *SimpleParameter) GetValue() string {
	return p.Value
}

type ParameterWithPlaceholders struct {
	Value        string
	Placeholders []Placeholder
}

func (p *ParameterWithPlaceholders) GetValue() string {
	outString := p.Value
	for _, placeholder := range p.Placeholders {
		outString = placeholder.Replace(outString)
	}
	return outString
}
