package internal

func GetTerraformInit() Command {
	return Command{
		Binary: "terraform",
		Parameters: []Parameter{
			&SimpleParameter{Value: "init"},
			&SimpleParameter{Value: "-input=false"},
		},
	}
}

func GetTerraformPlan() Command {
	return Command{
		Binary: "terraform",
		Parameters: []Parameter{
			&SimpleParameter{Value: "plan"},
			&SimpleParameter{Value: "-out=plan.tfplan"},
			&SimpleParameter{Value: "-input=false"},
		},
	}
}

func GetTerraformApply() Command {
	return Command{
		Binary: "terraform",
		Parameters: []Parameter{
			&SimpleParameter{Value: "apply"},
			&SimpleParameter{Value: "plan.tfplan"},
			&SimpleParameter{Value: "-input=false"},
		},
	}
}

func GetTerraformDestroy() Command {
	return Command{
		Binary: "terraform",
		Parameters: []Parameter{
			&SimpleParameter{Value: "destroy"},
			&SimpleParameter{Value: "-input=false"},
		},
	}
}
