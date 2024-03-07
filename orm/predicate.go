package orm

// 衍生类型
type Op string

func (o Op) expr() {
	//TODO implement me
	panic("implement me")
}

// 别名
//type Op=string

const (
	opEq  Op = "="
	opLt  Op = "<"
	opNot Op = "NOT"
	opAnd Op = "AND"
	opOr  Op = "OR"
)

func (o Op) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    Op
	right Expression
}

//func Eq(column string, arg any) Predicate {
//	return Predicate{
//		Column:column,
//		Op: "=",
//		Arg: arg,
//	}
//}

type Column struct {
	name string
}

func C(name string) Column {
	return Column{name: name}
}

// C("id").Eq(12)
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opEq,
		right: value{
			val: arg,
		},
	}
}

func (c Column) Lt(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opLt,
		right: value{
			val: arg,
		},
	}
}

func (c Column) expr() {
	//TODO implement me
	panic("implement me")
}

func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

// C(“id”).Eq(12).And(C(“name”).Eq("Tom"))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

// C(“id”).Eq(12).Or(C(“name”).Eq("Tom"))
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

func (left Predicate) expr() {
	//TODO implement me
	panic("implement me")
}

// Expression 标记接口，代表表达式
type Expression interface {
	expr()
}

type value struct {
	val any
}

func (v value) expr() {
	//TODO implement me
	panic("implement me")
}
