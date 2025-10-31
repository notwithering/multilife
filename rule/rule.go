package rule

import (
	"strconv"
	"strings"
)

type RuleConfig string

type CompiledRule struct {
	String     string
	Birth      []int
	Survive    []int
	BirthSet   [9]bool
	SurviveSet [9]bool
}

func (c RuleConfig) Compile() CompiledRule {
	var r CompiledRule
	r.String = string(c)

	s := string(c)
	birth := []int{}
	survive := []int{}

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
			birth = append(birth, n)
			r.BirthSet[n] = true
		}
	}
	for _, ch := range sPart {
		if n, err := strconv.Atoi(string(ch)); err == nil {
			survive = append(survive, n)
			r.SurviveSet[n] = true
		}
	}

	r.Birth = birth
	r.Survive = survive
	return r
}
