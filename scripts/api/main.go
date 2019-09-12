package main

import (
	"fmt"
	"unicode"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedFiles | packages.NeedName | packages.NeedImports | packages.NeedDeps,
	}
	ps, err := packages.Load(cfg, "./...")
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		for _, n := range p.Types.Scope().Names() {
			if unicode.IsUpper([]rune(n)[0]) {
				fmt.Printf("%v.%v\n", p, n)
			}
		}
	}
}
