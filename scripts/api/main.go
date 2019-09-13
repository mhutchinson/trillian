package main

import (
	"fmt"
	"go/types"
	"strings"
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
		s := p.Types.Scope()
		for _, n := range s.Names() {
			if unicode.IsUpper([]rune(n)[0]) {
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("%v.%v: ", p, n))
				switch t := s.Lookup(n).(type) {
				case *types.Func:
					var sig *types.Signature = t.Type().(*types.Signature)
					sb.WriteString(fmt.Sprintf("\tfunction (%v params, %v returns)", sig.Params().Len(), sig.Results().Len()))
				case *types.TypeName:
					numMethods, numFields := scrapeType(t)
					sb.WriteString(fmt.Sprintf("\tclass (%v methods, %v fields)", numMethods, numFields))
				default:
					sb.WriteString(fmt.Sprintf("\t%T", t))
				}
				fmt.Printf("%v\n", sb.String())
			}
		}
	}
}

// Returns the number of methods and fields in the type
func scrapeType(t *types.TypeName) (int, int) {
	var named *types.Named = t.Type().(*types.Named)

	var expFields int
	switch u := named.Underlying().(type) {
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			if u.Field(i).Exported() {
				expFields++
			}
		}
	default:
		// No fields to count I guess
	}

	var expMethods int
	for i := 0; i < named.NumMethods(); i++ {
		if named.Method(i).Exported() {
			expMethods++
		}
	}

	return expMethods, expFields
}
