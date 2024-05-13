package clienterrorsif

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const Doc = `report problematic usage of "error" in WorkerClientError details`

var Analyzer = &analysis.Analyzer{
	Name:     "jsonerrif",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return // not a static call
		}
		if fn.FullName() != "github.com/osbuild/osbuild-composer/internal/worker/clienterrors.WorkerClientError" {
			return
		}
		detailsArg := call.Args[2]
		if tv, ok := pass.TypesInfo.Types[detailsArg]; ok {
			switch tv.Type.String() {
			case "error":
				pass.Reportf(n.Pos(), "do not pass 'error' to WorkerClientError() details, use error.Error() instead")
			case "[]error":
				pass.Reportf(n.Pos(), "do not pass '[]error' to WorkerClientError() details, use []string instead")
			}
		}
	})
	return nil, nil
}
