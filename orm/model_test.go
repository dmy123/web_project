package orm

import (
	"awesomeProject1/orm/internal/errs"
	"github.com/stretchr/testify/assert"
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
		want    *model
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "test model",
			args: args{
				entity: TestModel{},
			},
			wantErr: errs.ErrPointerOnly,
			//want: &model{
			//	tableName: "test_model",
			//	fields: map[string]*field{
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
			want: &model{
				tableName: "test_model",
				fields: map[string]*field{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := MustNewDB()
			got, err := db.r.parseModel(tt.args.entity)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("parseModel() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equalf(t, tt.want, got, "parseModel(%v)", tt.args.entity)
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
		models map[reflect.Type]*model
	}
	type args struct {
		val any
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *model
		wantErr   error
		cacheSize int
	}{
		// TODO: Add test cases.
		{
			name: "pointer",
			args: args{
				val: &TestModel{},
			},
			want: &model{
				tableName: "test_model",
				fields: map[string]*field{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mustNewRegistry()
			if tt.fields.models != nil {
				r.models = tt.fields.models
			}
			gotM, err := r.get(tt.args.val)
			if err != nil || tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, gotM, "get(%v)", tt.args.val)
			assert.Equal(t, tt.cacheSize, len(r.models))

			// 确认缓存下来
			typ := reflect.TypeOf(tt.args.val)
			m, ok := r.models[typ]
			assert.True(t, ok, true)
			assert.Equal(t, tt.want, m)

		})
	}
}
