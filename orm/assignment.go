package orm

type Assignment struct {
	col string
	val []any
}

func (a Assignment) assign() {
	//TODO implement me
	panic("implement me")
}

func Assign(col string, val ...any) Assignment {
	return Assignment{
		col: col,
		val: val,
	}
}
