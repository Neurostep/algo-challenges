package radar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	_neOp      = "!="
	_eqOp      = "=="
	_ltOp      = "<"
	_lteOp     = "<="
	_gtOp      = ">"
	_gteOp     = ">="
	_boolOpAnd = "AND"
	_boolOpOr  = "OR"
	_chargeOp  = "CHARGE"
	_allowOp   = "ALLOW"
	_blockOp   = "BLOCK"
)

var (
	_glueExpr  = regexp.MustCompile(fmt.Sprintf("(.*)(%s|%s)(.*)", _boolOpAnd, _boolOpOr))
	// the order is important to preserve the arithmetic operations match (ex.: <= vs <)
	_arithExpr = regexp.MustCompile(fmt.Sprintf("(.*)(%s|%s|%s|%s|%s|%s)(.*)", _neOp, _eqOp, _gteOp, _gtOp , _lteOp, _ltOp))
)

type (
	Node interface {
		Eval(p map[string]string) bool
	}
	And struct {
		Lft, Rgt Node
	}
	Or struct {
		Lft, Rgt Node
	}
	Lt struct {
		Key, Value string
	}
	Lte struct {
		Key, Value string
	}
	Gt struct {
		Key, Value string
	}
	Gte struct {
		Key, Value string
	}
	Eq struct {
		Key, Value string
	}
	Ne struct {
		Key, Value string
	}
	True  struct{}
	False struct{}
)

func (t True) Eval(p map[string]string) bool {
	return true
}

func (f False) Eval(p map[string]string) bool {
	return false
}

func (e And) Eval(p map[string]string) bool {
	return e.Lft.Eval(p) && e.Rgt.Eval(p)
}

func (e Or) Eval(p map[string]string) bool {
	return e.Lft.Eval(p) || e.Rgt.Eval(p)
}

func (e Lt) Eval(p map[string]string) bool {
	if val, found := p[e.Key]; found {
		v1, _ := strconv.Atoi(val)
		v2, _ := strconv.Atoi(e.Value)
		return v1 < v2
	}
	return false
}

func (e Lte) Eval(p map[string]string) bool {
	if val, found := p[e.Key]; found {
		v1, _ := strconv.Atoi(val)
		v2, _ := strconv.Atoi(e.Value)
		return v1 <= v2
	}
	return false
}

func (e Gt) Eval(p map[string]string) bool {
	if val, found := p[e.Key]; found {
		v1, _ := strconv.Atoi(val)
		v2, _ := strconv.Atoi(e.Value)
		return v1 > v2
	}
	return false
}

func (e Gte) Eval(p map[string]string) bool {
	if val, found := p[e.Key]; found {
		v1, _ := strconv.Atoi(val)
		v2, _ := strconv.Atoi(e.Value)
		return v1 >= v2
	}
	return false
}

func (e Ne) Eval(p map[string]string) bool {
	if val, found := p[e.Key]; found {
		if val2, found2 := p[e.Value]; found2 {
			return val != val2
		}
		return val != e.Value
	}
	return false
}

func (e Eq) Eval(p map[string]string) bool {
	if val, found := p[e.Key]; found {
		return val == e.Value
	}
	return false
}

func parse(s string) Node {
	var node Node

	boolPair := _glueExpr.FindStringSubmatch(s)
	if len(boolPair) == 0 {
		return parseArithm(s)
	}
	boolOp := boolPair[2]
	if boolOp == _boolOpAnd {
		node = And{Lft: parseArithm(boolPair[1]), Rgt: parseArithm(boolPair[3])}
	}
	if boolOp == _boolOpOr {
		node = Or{Lft: parseArithm(boolPair[1]), Rgt: parseArithm(boolPair[3])}
	}

	return node
}

func parseArithm(s string) Node {
	tokens := _arithExpr.FindStringSubmatch(s)[1:]

	switch tokens[1] {
	case _neOp:
		return Ne{Key: tokens[0], Value: tokens[2]}
	case _eqOp:
		return Eq{Key: tokens[0], Value: tokens[2]}
	case _ltOp:
		return Lt{Key: tokens[0], Value: tokens[2]}
	case _lteOp:
		return Lte{Key: tokens[0], Value: tokens[2]}
	case _gtOp:
		return Gt{Key: tokens[0], Value: tokens[2]}
	case _gteOp:
		return Gte{Key: tokens[0], Value: tokens[2]}
	}
	return False{}
}

func Solution(ops []string) int {
	chargeM := make(map[string]string)
	var allowTree Node
	var blockTree Node

	allowTree = True{}
	blockTree = False{}

	for _, op := range ops {
		parts := strings.Split(strings.Replace(op, " ", "", -1), ":")
		opKw := strings.TrimSpace(parts[0])
		opStr := strings.TrimSpace(parts[1])

		switch opKw {
		case _chargeOp:
			pairs := strings.Split(opStr, "&")
			for _, pairS := range pairs {
				pair := strings.Split(pairS, "=")
				chargeM[pair[0]] = pair[1]
			}
			break
		case _allowOp:
			allowTree = parse(opStr)
			break
		case _blockOp:
			blockTree = parse(opStr)
			break
		}
	}

	if allowTree.Eval(chargeM) && !blockTree.Eval(chargeM) {
		return 1
	}
	return 0
}
