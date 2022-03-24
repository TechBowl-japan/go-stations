package sta10_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"
)

func TestStation10(t *testing.T) {
	t.Parallel()

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "../../model/error.go", nil, 0)
	if err != nil {
		t.Error("パースに失敗しました", err)
		return
	}

	config := &types.Config{
		Importer: importer.Default(),
	}

	pkg, err := config.Check("model", fset, []*ast.File{f}, nil)
	if err != nil {
		t.Error("パッケージチェックに失敗しました", err)
		return
	}

	obj := pkg.Scope().Lookup("ErrNotFound")
	if obj == nil {
		t.Error("ErrNotFound がみつかりません")
		return
	}

	_, ok := obj.Type().(*types.Named)
	if !ok {
		t.Error("ErrNotFound 型が見つかりません")
		return
	}

	typ := obj.Type()
	errInterface := types.Universe.Lookup("error").Type().Underlying().(*types.Interface)
	if !types.Implements(typ, errInterface) && !types.Implements(types.NewPointer(typ), errInterface) {
		t.Error("ErrNotFound に error interface が実装されていません")
	}
}
