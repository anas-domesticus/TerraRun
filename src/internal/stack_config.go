package internal

type StackConfig struct {
	Depends []Dependency
}

type Dependency string

func (d *Dependency) AsStack() TerraformStack {
	return TerraformStack{
		d.AsString(),
	}
}

func (d *Dependency) AsString() string {
	return string(*d)
}
