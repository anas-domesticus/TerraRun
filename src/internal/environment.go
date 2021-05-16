package internal

type Environment struct {
	Name string
}

func (e *Environment) GetEnvPlaceholder() Placeholder {
	return Placeholder{
		Before: "environment",
		After:  e.Name,
	}
}
