package orm

// Expression 标记接口，代表表达式
type Expression interface {
	expr()
}

type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) selectable() {
	//TODO implement me
	panic("implement me")
}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{raw: expr, args: args}
}

func (r RawExpr) expr() {
	//TODO implement me
	panic("implement me")
}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}
