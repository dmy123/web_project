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

type value struct {
	val any
}

func (v value) expr() {
	//TODO implement me
	panic("implement me")
}
