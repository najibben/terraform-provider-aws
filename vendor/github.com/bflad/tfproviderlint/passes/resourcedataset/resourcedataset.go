package resourcedataset

import (
	"go/ast"
	"reflect"

	"github.com/bflad/tfproviderlint/helper/terraformtype"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "resourcedataset",
	Doc:  "find github.com/hashicorp/terraform-plugin-sdk/helper/schema.ResourceData.Set() calls for later passes",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run:        run,
	ResultType: reflect.TypeOf([]*ast.CallExpr{}),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	var result []*ast.CallExpr

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		x := n.(*ast.CallExpr)

		if !isResourceDataSet(pass, x) {
			return
		}

		result = append(result, x)
	})

	return result, nil
}

func isResourceDataSet(pass *analysis.Pass, ce *ast.CallExpr) bool {
	switch f := ce.Fun.(type) {
	default:
		return false
	case *ast.SelectorExpr:
		if f.Sel.Name != "Set" {
			return false
		}

		switch x := f.X.(type) {
		default:
			return false
		case *ast.Ident:
			if x.Obj == nil {
				return false
			}

			switch decl := x.Obj.Decl.(type) {
			default:
				return false
			case *ast.Field:
				switch t := decl.Type.(type) {
				default:
					return false
				case *ast.StarExpr:
					return terraformtype.IsHelperSchemaTypeResourceData(pass.TypesInfo.TypeOf(t.X))
				case *ast.SelectorExpr:
					return terraformtype.IsHelperSchemaTypeResourceData(pass.TypesInfo.TypeOf(t))
				}
			case *ast.ValueSpec:
				return terraformtype.IsHelperSchemaTypeResourceData(pass.TypesInfo.TypeOf(decl.Type))
			}
		}
	}
}
