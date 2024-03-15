package orm

type Column struct {
	table TableReference
	name  string
	alias string
}

func C(name string) Column {
	return Column{name: name}
}

func (c Column) assign() {
	return
}

func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
		table: c.table,
	}
}

// C("id").Eq(12)
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valueOf(arg),
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: valueOf(arg),
	}
}

func valueOf(arg any) Expression {
	switch exp := arg.(type) {
	case Expression:
		return exp
	default:
		return value{
			val: arg,
		}
	}
}

func (c Column) expr() {
	//TODO implement me

}

func (c Column) selectable() {
	//TODO implement me

}
