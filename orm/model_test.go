package orm

import (
	"awesomeProject1/orm/internal/errs"
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
		opts    []ModelOption
	}{
		// TODO: Add test cases.
		{
			name: "test Model",
			args: args{
				entity: TestModel{},
			},
			wantErr: errs.ErrPointerOnly,
			//want: &Model{
			//	tableName: "test_model",
			//	fields: map[string]*Field{
			//		"Id": {
			//			colName: "id",
			//		},
			//		"FirstName": {
			//			colName: "first_name",
			//		},
			//		"LastName": {
			//			colName: "last_name",
			//		},
			//		"Age": {
			//			colName: "age",
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
				tableName: "test_model",
				fields: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
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
			db := MustNewDB()
			got, err := db.r.Registry(tt.args.entity, tt.opts...)
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
				tableName: "test_model",
				fields: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
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
				tableName: "tag_table",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name_t",
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
				tableName: "tag_table",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
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
				tableName: "tag_table",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
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
				tableName: "custom_table_name_t",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
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
				tableName: "custom_table_name_ptr_t",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
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
				tableName: "empty_table_name_ptr",
				fields: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
			cacheSize: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mustNewRegistry()
			//if tt.fields.models != nil {
			//	r.models = tt.fields.models
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
	m, err := r.Registry(&TestModel{}, ModelWithTableName("test_model_t"))
	assert.Equal(t, err, nil)
	assert.Equal(t, m.tableName, "test_model_t")
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
			m, err := r.Registry(&TestModel{}, ModelWithColumnName(tt.args.field, tt.args.colName))
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			fd, ok := m.fields[tt.args.field]
			require.True(t, ok)
			assert.Equal(t, tt.wantCol, fd.colName)
		})
	}
}
