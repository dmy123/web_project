package orm

import (
	"awesomeProject1/orm/internal/valuer"
	"awesomeProject1/orm/model"
)

type core struct {
	model *model.Model

	dialect Dialect
	creator valuer.Creator
	r       *model.Registry
}
