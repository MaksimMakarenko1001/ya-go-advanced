package exitcontrol

import (
	"go/ast"
	"go/types"

	"github.com/MaksimMakarenko1001/ya-go-advanced/pkg"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	fatalCalls = map[string]struct{}{
		"Fatal":   {},
		"Fatalf":  {},
		"Fatalln": {},
	}
	exitCalls = map[string]struct{}{
		"Exit": {},
	}

	Analyzer = &analysis.Analyzer{
		Name:     "exitcontrol",
		Doc:      `checks for built-in panic call and log.Fatal and os.Exit calls outside of main function in main package`,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run,
	}
)

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for function calls
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		if isPanicCall(pass, call) {
			pass.Reportf(call.Pos(), "panic call should not be used in production code")
			return
		}

		if isFatalCall(pass, call) {
			pass.Reportf(call.Pos(), "log.Fatal call should only be used in main function of main package")
			return
		}

		if isExitCall(pass, call) {
			pass.Reportf(call.Pos(), "os.Exit call should only be used in main function of main package")
		}

	})

	return nil, nil
}

func isPanicCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	fun, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.Uses[fun]
	if obj == nil {
		return false
	}

	_, ok = obj.(*types.Builtin)
	return ok && fun.Name == "panic"
}

func isFatalCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	se, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if _, ok := fatalCalls[se.Sel.Name]; !ok {
		return false
	}

	ident, ok := se.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	if pkg, ok := obj.(*types.PkgName); !ok || pkg.Imported().Path() != "log" {
		return false
	}

	return !isInEntryPoint(pass, call)
}

func isExitCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	se, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if _, ok := exitCalls[se.Sel.Name]; !ok {
		return false
	}

	ident, ok := se.X.(*ast.Ident)
	if !ok {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	if pkg, ok := obj.(*types.PkgName); !ok || pkg.Imported().Path() != "os" {
		return false
	}

	return !isInEntryPoint(pass, call)
}

// isInEntryPoint checks if the call goes inside the main function of the main package
func isInEntryPoint(pass *analysis.Pass, call *ast.CallExpr) bool {
	// Check if we're in the main package
	if pass.Pkg.Name() != "main" {
		return false
	}

	found := false
	files := pkg.SliceFilter(pass.Files, func(f *ast.File) bool {
		return call.Pos() >= f.FileStart && call.Pos() <= f.FileEnd
	})
	// Find the function declaration that contains this call
	for i := 0; i < len(files) && !found; i++ {
		ast.Inspect(files[i], func(node ast.Node) bool {
			fn, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
			// Check if the call is within this function's body
			if call.Pos() < fn.Pos() || call.Pos() > fn.End() {
				return true
			}
			// Check if this is the main function
			if fn.Name.Name != "main" {
				return false
			}
			if fn.Recv != nil {
				return false
			}

			found = true
			return false
		})
	}

	return found
}
