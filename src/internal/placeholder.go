package internal

import (
	"fmt"
	"strings"
)

type Placeholder struct {
	Before string
	After  string
}

func (p *Placeholder) Replace(input string) string {
	return strings.ReplaceAll(input, fmt.Sprintf("{{%s}}", p.Before), p.After)
}
