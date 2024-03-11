package orm

import (
	"awesomeProject1/orm/internal/errs"
)

var (
	DialectMySQL      Dialect = mysqlDialect{}
	DialectSQLite     Dialect = sqliteDialect{}
	DialectPostgreSQL Dialect = postgreDialect{}
)

type Dialect interface {
	quoter() byte
	buildOnDuplicateKey(b *builder, odk *OnDuplicateKey) error
}

type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildOnDuplicateKey(b *builder, odk *OnDuplicateKey) error {
	//TODO implement me
	panic("implement me")
}

type mysqlDialect struct {
	standardSQL
}

func (s mysqlDialect) quoter() byte {
	return '`'
}

func (m mysqlDialect) buildOnDuplicateKey(b *builder, odk *OnDuplicateKey) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for j, assign := range odk.assigns {
		switch exp := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldMap[exp.col]
			if !ok {
				return errs.NewErrUnknownField(exp.col)
			}
			if j > 0 {
				b.sb.WriteString(", ")
			}
			b.sb.WriteByte('`')
			b.sb.WriteString(fd.ColName)
			b.sb.WriteByte('`')
			b.sb.WriteString("=?")
			b.addArg(exp.val...)
			//b.args = append(b.args, exp.val...)
		case Column:
			fd, ok := b.model.FieldMap[exp.name]
			if !ok {
				return errs.NewErrUnknownField(exp.name)
			}
			if j > 0 {
				b.sb.WriteString(", ")
			}
			b.buildColumn(C(fd.GoName))
			b.sb.WriteString("=")
			b.sb.WriteString("VALUES")
			b.sb.WriteByte('(')
			b.buildColumn(C(fd.GoName))
			b.sb.WriteByte(')')
		default:
			return errs.NewErrUnsupportedAssignable(exp)
		}
	}
	return nil
}

type sqliteDialect struct {
	standardSQL
}

type postgreDialect struct {
	standardSQL
}
