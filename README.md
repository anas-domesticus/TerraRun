## TerraRun

WARNING: This software is currently very young & immature, potentially unsafe to use. Be careful.

This is a wrapper to Terraform, it allows you to run multiple Terraform stacks in a single command, it does not require any non-standard Terraform HCL code.

It is designed to make Terraform easier to use in automation. Unlike another well-known Terraform wrapper (TerraGrunt), it requires no non-standard HCL & doesn't incite users to box everything into highly parameterised modules.

It understands the concept of environments, which are based on the presence of a tfvars file in the root of the stack in question. This allows for some differences between environments, whilst encouraging consistency.

*i.e. if a stack has a file named "dev.tfvars", but no file called "prod.tfvars", it will run when an environment is not specified, or when "dev" is passed as the environment, but definitely not in "prod"*

The relevant tfvars file will be included in plans / applies to that environment.

It also has a collated HTML output format for plans, currently a little rudimentary, but readable nonetheless. This allows you to view all the changes across all your stacks in a single file. 

### Commands
- list - Lists eligible Terraform stacks
- validate - Performs a terraform validate against each of the stacks
- plan - Runs a plan against each of these stacks, saves a tfplan file in the stack directories
- apply - Runs the plan saved from teh previous step
- Help - Prints out a listing of available commands & flags

### Flags
- -d / --directory - Specifies the directory to base the search for Terraform stacks
- -e / --environment - Specifies an environment to run, if omitted it will ignore the concept of environments totally
- -r / --report - Outputs an HTML report detailing the changes to each of the stacks, only possible when using the "plan" command
*Be warned, the HTML is incredibly ugly, PRs are welcome to improve this :)*

## Examples:
`terrarun plan -d ./terraform-files -r` - Plans all stacks under ./terraform-files, outputs a report to report.html

`terrarun apply -e dev` - Applies plans already created in all Terraform stacks, searching from current directory, will fail if plans do not exist, will use dev.tfvars files

## Building from source:

- Navigate to src/cli
- `go build -o terrarun`

This will create a binary called terrarun, put it whereever you like.

Alternatively, the provided Dockerfile runs the test suite, then builds a clean docker image with the latest version of Terraform built in

## Future stuff:
- Running tasks in parallel
- YAML config for custom commands, linting etc...