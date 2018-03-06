package calc

import "math"

var DefaultContext Context = defaultContext{}

type defaultContext struct {
}

func (ctx defaultContext) Var(name string) func() interface{} {
	return nil
}

func (ctx defaultContext) Call(name string) func(args ...interface{}) interface{} {
	return nil
}

func (ctx defaultContext) Add(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 + f2
}

func (ctx defaultContext) Sub(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 - f2
}

func (ctx defaultContext) Mul(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 * f2
}

func (ctx defaultContext) Div(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 / f2
}

func (ctx defaultContext) Rem(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return math.Mod(f1, f2)
}

func (ctx defaultContext) Equal(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 == f2
}

func (ctx defaultContext) NotEqual(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 != f2
}

func (ctx defaultContext) Less(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 < f2
}

func (ctx defaultContext) Greater(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 > f2
}

func (ctx defaultContext) LessEqual(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 <= f2
}

func (ctx defaultContext) GreaterEqual(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return f1 >= f2

}

func (ctx defaultContext) LogicalOr(a, b interface{}) interface{} {
	b1, _ := a.(bool)
	b2, _ := b.(bool)
	return b1 || b2

}

func (ctx defaultContext) LogicalAnd(a, b interface{}) interface{} {
	b1, _ := a.(bool)
	b2, _ := b.(bool)
	return b1 && b2
}

func (ctx defaultContext) BitOr(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return float64(int64(f1) | int64(f2))
}

func (ctx defaultContext) BitAnd(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return float64(int64(f1) & int64(f2))
}

func (ctx defaultContext) BitXor(a, b interface{}) interface{} {
	f1, _ := a.(float64)
	f2, _ := b.(float64)
	return float64(int64(f1) ^ int64(f2))
}

func (ctx defaultContext) Sub1(a interface{}) interface{} {
	f1, _ := a.(float64)
	return -f1
}

func (ctx defaultContext) Not(a interface{}) interface{} {
	b1, _ := a.(bool)
	return !b1
}
