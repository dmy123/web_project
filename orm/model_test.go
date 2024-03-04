package orm

import (
	"awesomeProject1/orm/internal/errs"
	"github.com/stretchr/testify/assert"
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
			got, err := parseModel(tt.args.entity)
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
