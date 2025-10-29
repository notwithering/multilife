package main

import (
	"strconv"
	"strings"
)

type ruleConfig string

type compiledRule struct {
	string     string
	birth      []int
	survive    []int
	birthSet   [9]bool
	surviveSet [9]bool
}

func (c ruleConfig) compile() compiledRule {
	var r compiledRule
	r.string = string(c)

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
			r.birthSet[n] = true
		}
	}
	for _, ch := range sPart {
		if n, err := strconv.Atoi(string(ch)); err == nil {
			survive = append(survive, n)
			r.surviveSet[n] = true
		}
	}

	r.birth = birth
	r.survive = survive
	return r
}
