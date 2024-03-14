package orm

type TableReference interface {
	table()
}

type Table struct {
	entity any
	alias  string
}

func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

func (t Table) C(name string) Column {
	return Column{
		name:  name,
		table: t,
	}
}

func (t Table) table() {
	//TODO implement me
	panic("implement me")
}
func (t Table) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: right,
		typ:   "JOIN",
	}
}
func (t Table) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: right,
		typ:   "LEFT JOIN",
	}
}
func (t Table) RightJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: right,
		typ:   "RIGHT JOIN",
	}
}

type Join struct {
	left  TableReference
	right TableReference
	typ   string
	on    []Predicate
	using []string
}

func (j Join) table() {
	//TODO implement me
	panic("implement me")
}

func (j Join) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: right,
		typ:   "JOIN",
	}
}
func (j Join) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: right,
		typ:   "LEFT JOIN",
	}
}
func (j Join) RightJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: right,
		typ:   "RIGHT JOIN",
	}
}

type JoinBuilder struct {
	left  TableReference
	right TableReference
	typ   string
}

// On t3:=t1.Join(t2).On(C("Id").Eq("RefId"))
func (j *JoinBuilder) On(ps ...Predicate) Join {
	return Join{
		left:  j.left,
		right: j.right,
		typ:   j.typ,
		on:    ps,
	}
}

// Using t3:=t1.Join(t2).Using(col1)
func (j *JoinBuilder) Using(cols ...string) Join {
	return Join{
		left:  j.left,
		right: j.right,
		typ:   j.typ,
		using: cols,
	}
}
