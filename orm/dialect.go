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
	buildUpsert(b *builder, upsert *Upsert) error
}

type standardSQL struct {
}

func (s standardSQL) quoter() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) buildUpsert(b *builder, upsert *Upsert) error {
	//TODO implement me
	panic("implement me")
}

type mysqlDialect struct {
	standardSQL
}

func (s mysqlDialect) quoter() byte {
	return '`'
}

func (m mysqlDialect) buildUpsert(b *builder, upsert *Upsert) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for j, assign := range upsert.assigns {
		switch exp := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldMap[exp.col]
			if !ok {
				return errs.NewErrUnknownField(exp.col)
			}
			if j > 0 {
				b.sb.WriteString(", ")
			}
			b.buildColumn(C(fd.GoName))
			//b.sb.WriteByte('`')
			//b.sb.WriteString(fd.ColName)
			//b.sb.WriteByte('`')
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

func (s sqliteDialect) buildUpsert(b *builder, upsert *Upsert) error {
	b.sb.WriteString(" ON CONFLICT(")
	for i, col := range upsert.conflictColumns {
		if i > 0 {
			b.sb.WriteString(", ")
		}
		err := b.buildColumn(C(col))
		if err != nil {
			return err
		}
	}
	b.sb.WriteString(") DO UPDATE SET ")
	for j, assign := range upsert.assigns {
		switch exp := assign.(type) {
		case Assignment:
			fd, ok := b.model.FieldMap[exp.col]
			if !ok {
				return errs.NewErrUnknownField(exp.col)
			}
			if j > 0 {
				b.sb.WriteString(", ")
			}
			b.buildColumn(C(fd.GoName))
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
			b.sb.WriteString("=excluded.")
			//b.sb.WriteByte('(')
			b.buildColumn(C(fd.GoName))
			//b.sb.WriteByte(')')
		default:
			return errs.NewErrUnsupportedAssignable(exp)
		}
	}
	return nil
}

func (s sqliteDialect) quoter() byte {
	return '`'
}

type postgreDialect struct {
	standardSQL
}
