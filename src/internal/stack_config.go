package internal

type StackConfig struct {
	Depends []Dependency
}

type Dependency string

func (d *Dependency) AsString() string {
	return string(*d)
}
