package orm

// Aggregate 聚合函数
// AVG("age"),  SUM("age"), COUNT("age"), MAX("age"), MIN("age")
type Aggregate struct {
	fn    string
	arg   string
	alias string
}

func (a Aggregate) selectable() {
	//TODO implement me
	panic("implement me")
}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
	}
}

func Avg(column string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: column,
	}
}

func Sum(column string) Aggregate {
	return Aggregate{
		fn:  "SUM",
		arg: column,
	}
}
func Count(column string) Aggregate {
	return Aggregate{
		fn:  "COUNT",
		arg: column,
	}
}
func Max(column string) Aggregate {
	return Aggregate{
		fn:  "MAX",
		arg: column,
	}
}
func Min(column string) Aggregate {
	return Aggregate{
		fn:  "MIN",
		arg: column,
	}
}
