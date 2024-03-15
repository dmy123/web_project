package orm_gen

import (
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"text/template"
)

//go:embed tpl.gohtml
var genOrm string

func gen(writer io.Writer, srcFile string) error {
	// ast解析
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	s := &SingleFileEntryVisitor{}
	ast.Walk(s, f)
	file := s.Get()

	// 操作模板
	tpl := template.New("gen_orm")
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		return err
	}
	return tpl.Execute(writer, file)
}
