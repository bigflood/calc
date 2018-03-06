package calc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

type Calc interface {
	Eval() interface{}
}

type Context interface {
	Var(name string) func() interface{}
	Call(name string) func(args ...interface{}) interface{}

	Add(a, b interface{}) interface{}
	Sub(a, b interface{}) interface{}
	Mul(a, b interface{}) interface{}
	Div(a, b interface{}) interface{}
	Rem(a, b interface{}) interface{}
	Equal(a, b interface{}) interface{}
	NotEqual(a, b interface{}) interface{}
	Less(a, b interface{}) interface{}
	Greater(a, b interface{}) interface{}
	LessEqual(a, b interface{}) interface{}
	GreaterEqual(a, b interface{}) interface{}
	LogicalOr(a, b interface{}) interface{}
	LogicalAnd(a, b interface{}) interface{}
	BitOr(a, b interface{}) interface{}
	BitAnd(a, b interface{}) interface{}
	BitXor(a, b interface{}) interface{}

	Sub1(a interface{}) interface{}
	Not(a interface{}) interface{}
}

func New(ctx Context, src string) (Calc, error) {
	expr, err := parser.ParseExpr(src)
	if err != nil {
		return nil, err
	}

	f, err := evalData{ctx: ctx, src: src}.evalExpr(expr)
	if err != nil {
		return nil, err
	}

	p := new(calc)
	p.ctx = ctx
	p.f = f
	return p, nil
}

type calc struct {
	ctx Context
	f   func() interface{}
}

func (p *calc) Eval() interface{} {
	return p.f()
}

var zeroFunc = func() interface{} { return 0 }

func ifThanElse(i bool, a, b interface{}) interface{} {
	if i {
		return a
	}
	return b
}

type evalData struct {
	ctx Context
	src string
}

func parseError(e ast.Expr, msg string) error {
	return fmt.Errorf("%d:%d: %s", e.Pos(), e.End(), msg)
}

func (data evalData) evalExpr(e ast.Expr) (func() interface{}, error) {
	switch v := e.(type) {
	case *ast.Ident:
		callFunc := data.ctx.Var(v.Name)
		if callFunc == nil {
			return zeroFunc, parseError(v, "unknown Ident")
		}
		return callFunc, nil
	case *ast.CallExpr:
		funcName, ok := v.Fun.(*ast.Ident)
		if !ok {
			return zeroFunc, parseError(v, "invalid function call")
		}

		callFunc := data.ctx.Call(funcName.Name)

		if callFunc == nil {
			return zeroFunc, parseError(funcName, "unknown Ident")
		}

		args := make([]func() interface{}, len(v.Args))
		for i, a := range v.Args {
			f, err := data.evalExpr(a)
			if err != nil {
				return zeroFunc, err
			}
			args[i] = f
		}

		argsBuf := make([]interface{}, len(args))
		return func() interface{} {
			for i, v := range args {
				argsBuf[i] = v()
			}
			return callFunc(argsBuf...)
		}, nil
	case *ast.BasicLit:
		s := strings.TrimSpace(v.Value)
		if strings.Contains(s, ".") {
			f, err := strconv.ParseFloat(s, 64)
			return func() interface{} { return f }, err
		}
		var i int64
		var err error

		if strings.HasPrefix(s, "0x") {
			i, err = strconv.ParseInt(s[2:], 16, 64)
		} else if len(s) > 1 && strings.HasPrefix(s, "0") {
			i, err = strconv.ParseInt(s[1:], 8, 64)
		} else {
			i, err = strconv.ParseInt(s, 10, 64)
		}
		f := float64(i)
		return func() interface{} { return f }, err
	case *ast.ParenExpr:
		return data.evalExpr(v.X)
	case *ast.UnaryExpr:
		switch v.Op {
		case token.ADD:
			return data.evalExpr(v.X)
		case token.SUB:
			f, err := data.evalExpr(v.X)
			return func() interface{} { return data.ctx.Sub1(f()) }, err
		case token.NOT:
			f, err := data.evalExpr(v.X)
			return func() interface{} { return ifThanElse(f() != 0, 1, 0) }, err
		default:
			return zeroFunc, parseError(v, "unknown Unary Op")
		}
	case *ast.BinaryExpr:
		a, err := data.evalExpr(v.X)
		if err != nil {
			return zeroFunc, err
		}
		b, err := data.evalExpr(v.Y)
		if err != nil {
			return zeroFunc, err
		}
		switch v.Op {
		case token.ADD:
			return func() interface{} { return data.ctx.Add(a(), b()) }, nil
		case token.SUB:
			return func() interface{} { return data.ctx.Sub(a(), b()) }, nil
		case token.MUL:
			return func() interface{} { return data.ctx.Mul(a(), b()) }, nil
		case token.QUO:
			return func() interface{} { return data.ctx.Div(a(), b()) }, nil
		case token.REM:
			return func() interface{} { return data.ctx.Rem(a(), b()) }, nil
		case token.EQL:
			return func() interface{} { return data.ctx.Equal(a(), b()) }, nil
		case token.LSS:
			return func() interface{} { return data.ctx.Less(a(), b()) }, nil
		case token.GTR:
			return func() interface{} { return data.ctx.Greater(a(), b()) }, nil
		case token.NEQ:
			return func() interface{} { return data.ctx.NotEqual(a(), b()) }, nil
		case token.LEQ:
			return func() interface{} { return data.ctx.LessEqual(a(), b()) }, nil
		case token.GEQ:
			return func() interface{} { return data.ctx.GreaterEqual(a(), b()) }, nil
		case token.LOR:
			return func() interface{} { return data.ctx.LogicalOr(a(), b()) }, nil
		case token.LAND:
			return func() interface{} { return data.ctx.LogicalAnd(a(), b()) }, nil
		case token.OR:
			return func() interface{} { return data.ctx.BitOr(a(), b()) }, nil
		case token.AND:
			return func() interface{} { return data.ctx.BitAnd(a(), b()) }, nil
		case token.XOR:
			return func() interface{} { return data.ctx.BitXor(a(), b()) }, nil
		default:
			return zeroFunc, parseError(v, "unknown Binary Op")
		}
	default:
		return zeroFunc, parseError(e, "invalid expression")
	}
}
