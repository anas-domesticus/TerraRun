package internal

type Environment struct {
	Name string
}

func (e *Environment) GetEnvPlaceholder() Placeholder {
	return Placeholder{
		Before: "Environment",
		After:  e.Name,
	}
}
