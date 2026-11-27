package session

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hanchuanchuan/goInception/ast"
	"github.com/hanchuanchuan/goInception/config"
	"github.com/hanchuanchuan/goInception/model"
	"github.com/hanchuanchuan/goInception/types"
)

func TestResolveInnodbLargePrefix(t *testing.T) {
	tests := []struct {
		name           string
		dbType         int
		dbVersion      int
		variableExists bool
		variableValue  bool
		want           bool
	}{
		{
			name:      "oceanbase without mysql compatibility variable",
			dbType:    DBTypeOceanBase,
			dbVersion: 4,
			want:      true,
		},
		{
			name:           "oceanbase ignores unsupported mysql limit variable",
			dbType:         DBTypeOceanBase,
			dbVersion:      4,
			variableExists: true,
			variableValue:  false,
			want:           true,
		},
		{
			name:      "mysql 5.6 without variable",
			dbType:    DBTypeMysql,
			dbVersion: 50600,
			want:      false,
		},
		{
			name:      "mysql 5.7.8 without removed variable",
			dbType:    DBTypeMysql,
			dbVersion: 50708,
			want:      true,
		},
		{
			name:           "mysql explicit disabled variable",
			dbType:         DBTypeMysql,
			dbVersion:      50708,
			variableExists: true,
			variableValue:  false,
			want:           false,
		},
		{
			name:           "mysql explicit enabled variable",
			dbType:         DBTypeMysql,
			dbVersion:      50600,
			variableExists: true,
			variableValue:  true,
			want:           true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := resolveInnodbLargePrefix(
				test.dbType,
				test.dbVersion,
				test.variableExists,
				test.variableValue,
			)
			if got != test.want {
				t.Fatalf("resolveInnodbLargePrefix() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestOceanBaseCompositeIndexUsesLargePrefixLimit(t *testing.T) {
	record := &Record{Buf: new(bytes.Buffer)}
	s := &session{
		dbType:              DBTypeOceanBase,
		dbVersion:           4,
		databaseCharset:     "utf8mb4",
		innodbLargePrefix:   resolveInnodbLargePrefix(DBTypeOceanBase, 4, false, false),
		myRecord:            record,
		recordSets:          NewRecordSets(),
		opt:                 &SourceOptions{},
		inc:                 config.Inc{Lang: "en-US"},
		incLevel:            make(map[string]uint8),
		stage:               StageCheck,
		lowerCaseTableNames: 1,
	}
	table := &TableInfo{
		Schema: "test",
		Name:   "hhhhh",
		IsNew:  true,
		Fields: []FieldInfo{
			{Field: "namespace", Type: "varchar(256)", Collation: "utf8mb4_general_ci", Null: "NO"},
			{Field: "release_id", Type: "varchar(36)", Collation: "utf8mb4_general_ci", Null: "NO"},
		},
	}
	columns := []*ast.IndexColName{
		{
			Column: &ast.ColumnName{Name: model.NewCIStr("namespace")},
			Length: types.UnspecifiedLength,
		},
		{
			Column: &ast.ColumnName{Name: model.NewCIStr("release_id")},
			Length: types.UnspecifiedLength,
		},
	}

	s.checkCreateIndex(
		&ast.TableName{Schema: model.NewCIStr("test"), Name: model.NewCIStr("hhhhh")},
		"idx_namespace_release_id",
		columns,
		nil,
		table,
		false,
		ast.ConstraintIndex,
	)

	if message := record.Buf.String(); strings.Contains(message, "max key length") {
		t.Fatalf("OceanBase index was checked against the legacy 767-byte limit: %s", message)
	}
}
