// Package S011 defines an Analyzer that checks for
// Schema with only Computed enabled and DiffSuppressFunc configured
package S011

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/bflad/tfproviderlint/helper/terraformtype"
	"github.com/bflad/tfproviderlint/passes/commentignore"
	"github.com/bflad/tfproviderlint/passes/schemaschema"
)

const Doc = `check for Schema with only Computed enabled and DiffSuppressFunc configured

The S011 analyzer reports cases of schemas which only enables Computed
and configures DiffSuppressFunc, which will fail provider schema validation.`

const analyzerName = "S011"

var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  Doc,
	Requires: []*analysis.Analyzer{
		schemaschema.Analyzer,
		commentignore.Analyzer,
	},
	Run: run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	ignorer := pass.ResultOf[commentignore.Analyzer].(*commentignore.Ignorer)
	schemas := pass.ResultOf[schemaschema.Analyzer].([]*terraformtype.HelperSchemaSchemaInfo)
	for _, schema := range schemas {
		if ignorer.ShouldIgnore(analyzerName, schema.AstCompositeLit) {
			continue
		}

		if !schema.Schema.Computed || schema.Schema.Optional || schema.Schema.Required {
			continue
		}

		if schema.Schema.DiffSuppressFunc == nil {
			continue
		}

		switch t := schema.AstCompositeLit.Type.(type) {
		default:
			pass.Reportf(schema.AstCompositeLit.Lbrace, "%s: schema should not only enable Computed and configure DiffSuppressFunc", analyzerName)
		case *ast.SelectorExpr:
			pass.Reportf(t.Sel.Pos(), "%s: schema should not only enable Computed and configure DiffSuppressFunc", analyzerName)
		}
	}

	return nil, nil
}
