package rule

import (
	"strconv"
	"strings"
)

type RuleConfig string

type CompiledRule struct {
	BirthSet   [9]bool
	SurviveSet [9]bool
}

func (c RuleConfig) Compile() CompiledRule {
	var r CompiledRule

	s := string(c)

	var bPart, sPart string
	for p := range strings.SplitSeq(s, "/") {
		if len(p) == 0 {
			continue
		}
		switch p[0] {
		case 'B':
			bPart = p[1:]
		case 'S':
			sPart = p[1:]
		}
	}

	for _, ch := range bPart {
		if n, err := strconv.Atoi(string(ch)); err == nil {
			r.BirthSet[n] = true
		}
	}
	for _, ch := range sPart {
		if n, err := strconv.Atoi(string(ch)); err == nil {
			r.SurviveSet[n] = true
		}
	}

	return r
}
