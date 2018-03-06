package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bigflood/calc"
)

type ctx struct {
	calc.Context
}

func (p *ctx) HasVar(name string) bool {
	switch name {
	case "x":
		return true
	case "sum":
		return true
	default:
		return false
	}
}

func (p *ctx) Var(name string) func() interface{} {
	switch name {
	case "x":
		return func() interface{} { return 123.0 }
	default:
		return p.Context.Var(name)
	}
}

func (p *ctx) Call(name string) func(args ...interface{}) interface{} {
	switch name {
	case "sum":
		return func(args ...interface{}) interface{} {
			sum := 0.0
			for _, v := range args {
				f, _ := v.(float64)
				sum += f
			}
			return sum
		}
	default:
		return p.Context.Call(name)
	}
}

func main() {
	for _, s := range os.Args[1:] {
		p1 := time.Now()
		f, err := calc.New(&ctx{Context: calc.DefaultContext}, s)
		p2 := time.Now()
		if err != nil {
			fmt.Println(s, " -> ", err)
			continue
		}

		fmt.Println("parse time:", p2.Sub(p1))

		p3 := time.Now()
		v := f.Eval()
		p4 := time.Now()

		fmt.Println("eval time:", p4.Sub(p3))
		fmt.Printf("'%s' -> '%v'\n", s, v)
	}
}
