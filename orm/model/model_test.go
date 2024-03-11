package model

import (
	"awesomeProject1/orm/internal/errs"
	"database/sql"
	_ "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_parseModel(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name    string
		args    args
		want    *Model
		wantErr error
		opts    []Option
	}{
		// TODO: Add test cases.
		{
			name: "test Model",
			args: args{
				entity: TestModel{},
			},
			wantErr: errs.ErrPointerOnly,
			//want: &Model{
			//	TableName: "test_model",
			//	FieldMap: map[string]*Field{
			//		"Id": {
			//			ColName: "id",
			//		},
			//		"FirstName": {
			//			ColName: "first_name",
			//		},
			//		"LastName": {
			//			ColName: "last_name",
			//		},
			//		"Age": {
			//			ColName: "age",
			//		},
			//	},
			//},
		},
		{
			name: "pointer",
			args: args{
				entity: &TestModel{},
			},
			want: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						GoName:  "Id",
						Typ:     reflect.TypeOf(int64(0)),
						Offset:  0,
					},
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					"LastName": {
						ColName: "last_name",
						GoName:  "LastName",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
					"Age": {
						ColName: "age",
						GoName:  "Age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						GoName:  "Id",
						Typ:     reflect.TypeOf(int64(0)),
						Offset:  0,
					},
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					"last_name": {
						ColName: "last_name",
						GoName:  "LastName",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
					"age": {
						ColName: "age",
						GoName:  "Age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
				},
				Fields: []*Field{
					{
						ColName: "id",
						GoName:  "Id",
						Typ:     reflect.TypeOf(int64(0)),
						Offset:  0,
					},
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					{
						ColName: "age",
						GoName:  "Age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
					{
						ColName: "last_name",
						GoName:  "LastName",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
				},
			},
		},

		{
			name: "basic types",
			args: args{
				entity: 0,
			},
			wantErr: errs.ErrPointerOnly,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//db := orm.memoryDB(t)
			r := MustNewRegistry()
			got, err := r.Registry(tt.args.entity, tt.opts...)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("Registry() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equalf(t, tt.want, got, "Registry(%v)", tt.args.entity)
		})
	}
}

func Test_underscoreName(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		//  ID 不能转化为 id ,只能转化为 i_d
		{
			name:    "upper cases",
			srcStr:  "ID",
			wantStr: "i_d",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := underscoreName(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}

func Test_registry_get(t *testing.T) {
	type fields struct {
		models map[reflect.Type]*Model
	}
	type args struct {
		val any
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *Model
		wantErr   error
		cacheSize int
	}{
		// TODO: Add test cases.
		{
			name: "pointer",
			args: args{
				val: &TestModel{},
			},
			want: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						GoName:  "Id",
						Typ:     reflect.TypeOf(int64(0)),
					},
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					"LastName": {
						ColName: "last_name",
						GoName:  "LastName",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
					"Age": {
						ColName: "age",
						GoName:  "Age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						GoName:  "Id",
						Typ:     reflect.TypeOf(int64(0)),
					},
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					"last_name": {
						ColName: "last_name",
						GoName:  "LastName",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
					"age": {
						ColName: "age",
						GoName:  "Age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
				},
				Fields: []*Field{
					{
						ColName: "id",
						GoName:  "Id",
						Typ:     reflect.TypeOf(int64(0)),
						Offset:  0,
					},
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						Offset:  8,
					},
					{
						ColName: "age",
						GoName:  "Age",
						Typ:     reflect.TypeOf(int8(0)),
						Offset:  24,
					},
					{
						ColName: "last_name",
						GoName:  "LastName",
						Typ:     reflect.TypeOf(&sql.NullString{}),
						Offset:  32,
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "tag",
			args: args{
				func() any {
					type TagTable struct {
						FirstName string `orm:"column=first_name_t"`
					}
					return &TagTable{}
				}(),
			},
			want: &Model{
				TableName: "tag_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name_t",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"first_name_t": {
						ColName: "first_name_t",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name_t",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						//Offset:  8,
					},
				},
			},

			cacheSize: 1,
		},
		{
			name: "empty column",
			args: args{
				func() any {
					type TagTable struct {
						FirstName string `orm:"column="`
					}
					return &TagTable{}
				}(),
			},
			want: &Model{
				TableName: "tag_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						//Offset:  8,
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "column only",
			args: args{
				func() any {
					type TagTable struct {
						FirstName string `orm:"column"`
					}
					return &TagTable{}
				}(),
			},

			wantErr:   errs.NewErrInvalidTagContent("column"),
			cacheSize: 1,
		},
		{
			name: "ignore tag",
			args: args{
				func() any {
					type TagTable struct {
						FirstName string `orm:"abc=abc"`
					}
					return &TagTable{}
				}(),
			},
			want: &Model{
				TableName: "tag_table",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						//Offset:  8,
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "table name",
			args: args{
				val: &CustomTableName{},
			},
			want: &Model{
				TableName: "custom_table_name_t",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						//Offset:  8,
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "table name ptr",
			args: args{
				val: &CustomTableNamePtr{},
			},
			want: &Model{
				TableName: "custom_table_name_ptr_t",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						//Offset:  8,
					},
				},
			},
			cacheSize: 1,
		},
		{
			name: "empty table name ptr",
			args: args{
				val: &EmptyTableNamePtr{},
			},
			want: &Model{
				TableName: "empty_table_name_ptr",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name",
						GoName:  "FirstName",
						Typ:     reflect.TypeOf(""),
						//Offset:  8,
					},
				},
			},
			cacheSize: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MustNewRegistry()
			//if tt.FieldMap.models != nil {
			//	r.models = tt.FieldMap.models
			//}
			gotM, err := r.Get(tt.args.val)
			if err != nil || tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, gotM, "get(%v)", tt.args.val)
			//assert.Equal(t, tt.cacheSize, len(r.models.))

			// 确认缓存下来
			typ := reflect.TypeOf(tt.args.val)
			m, ok := r.models.Load(typ)
			assert.True(t, ok, true)
			assert.Equal(t, tt.want, m)

		})
	}
}

type CustomTableName struct {
	FirstName string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
	FirstName string
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_ptr_t"
}

type EmptyTableNamePtr struct {
	FirstName string
}

func (c *EmptyTableNamePtr) TableName() string {
	return ""
}

func TestModelWithTableName(t *testing.T) {
	r, _ := newRegistry()
	m, err := r.Registry(&TestModel{}, WithTableName("test_model_t"))
	assert.Equal(t, err, nil)
	assert.Equal(t, m.TableName, "test_model_t")
}

func TestModelWithColumnName(t *testing.T) {
	type args struct {
		field   string
		colName string
	}
	tests := []struct {
		name    string
		args    args
		wantCol string
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "column name",
			args: args{
				field:   "FirstName",
				colName: "first_name_ccc",
			},
			wantCol: "first_name_ccc",
		},
		{
			name: "invalid column name",
			args: args{
				field:   "XXX",
				colName: "first_name_ccc",
			},
			wantErr: errs.NewErrUnknownField("XXX"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := newRegistry()
			m, err := r.Registry(&TestModel{}, WithColumnName(tt.args.field, tt.args.colName))
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			fd, ok := m.FieldMap[tt.args.field]
			require.True(t, ok)
			assert.Equal(t, tt.wantCol, fd.ColName)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
