// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package session

import (
	"bytes"
	"testing"

	"github.com/hanchuanchuan/goInception/ast"
	"github.com/hanchuanchuan/goInception/parser"
	"github.com/hanchuanchuan/goInception/sessionctx/variable"
)

func Test_checkDDLInstantMySQL(t *testing.T) {
	parser := parser.New()
	sessionVars := variable.NewSessionVars()
	parser.SetSQLMode(sessionVars.SQLMode)

	table := &TableInfo{Name: "t1",
		Fields: []FieldInfo{
			{Field: "c1", Extra: ""},
			{Field: "c2", Extra: "VIRTUAL"},
			{Field: "c3", Extra: "STORED"}}}

	tests := []struct {
		sql     string
		mysql57 bool
		mysql80 bool
	}{
		{"alter table t1 add column c1 int;", false, true},
		{"ALTER TABLE t1 ADD COLUMN c2 INT GENERATED ALWAYS AS (c1 + 1) VIRTUAL;", true, true},
		{"ALTER TABLE t1 ADD COLUMN c4 INT AFTER c1;", false, false},
		{"ALTER TABLE t1 ADD COLUMN c3 INT GENERATED ALWAYS AS (c1 + 1) STORED;", false, false},
	}

	for _, test := range tests {
		stmts, _, err := parser.Parse(test.sql, "", "")
		if err != nil {
			t.Error(err)
		}
		for _, stmtNode := range stmts {
			switch node := stmtNode.(type) {
			case *ast.AlterTableStmt:
				canInstant := checkDDLInstantMySQL57(node)
				if canInstant != test.mysql57 {
					t.Errorf("canInstant is %v, but excepted %v: sql: %v", canInstant, test.mysql57, test.sql)
				}

				canInstant = checkDDLInstantMySQL80(node, table, 80000)
				if canInstant != test.mysql80 {
					t.Errorf("canInstant is %v, but excepted %v: sql: %v", canInstant, test.mysql57, test.sql)
				}
			}
		}
	}
}

func newAuditTestSession(dbType int, lang string, offlineCheck bool, incLevel map[string]uint8) *session {
	dbVersion := 80036
	if dbType == DBTypeOceanBase {
		dbVersion = 4
	}
	s := &session{
		dbType:      dbType,
		dbVersion:   dbVersion,
		dbName:      "demo_db",
		opt:         &SourceOptions{},
		myRecord:    &Record{Buf: new(bytes.Buffer)},
		recordSets:  &MyRecordSets{},
		incLevel:    incLevel,
		dbCacheList: map[string]*DBInfo{"demo_db": {Name: "demo_db"}},
		tableCacheList: map[string]*TableInfo{
			"demo_db.demo_table": {Schema: "demo_db", Name: "demo_table", IsDeleted: true},
		},
	}
	s.inc.EnableNullable = true
	s.inc.Lang = lang
	s.inc.CheckOfflineDDL = offlineCheck
	return s
}

func TestColumnTypeChangeMessage(t *testing.T) {
	cases := []struct {
		name         string
		sql          string
		oldType      string
		dbType       int
		offlineCheck bool
		lang         string
		incLevel     map[string]uint8
		want         string
		wantLevel    uint8
	}{
		{
			name:      "mysql_zh_bigint_to_mediumtext",
			sql:       "ALTER TABLE demo_table MODIFY COLUMN demo_column MEDIUMTEXT NULL COMMENT '示例字段'",
			oldType:   "bigint",
			dbType:    DBTypeMysql,
			lang:      "zh_cn",
			want:      "字段 'demo_table.demo_column' 类型变更: bigint -> mediumtext.",
			wantLevel: 1,
		},
		{
			name:      "mysql_en_bigint_to_mediumtext",
			sql:       "ALTER TABLE demo_table MODIFY COLUMN demo_column MEDIUMTEXT NULL COMMENT '示例字段'",
			oldType:   "bigint",
			dbType:    DBTypeMysql,
			lang:      "en_us",
			want:      "Type conversion warning for column 'demo_table.demo_column': bigint -> mediumtext.",
			wantLevel: 1,
		},
		{
			name:      "tidb_zh_bigint_to_mediumtext",
			sql:       "ALTER TABLE demo_table MODIFY COLUMN demo_column MEDIUMTEXT NULL COMMENT '示例字段'",
			oldType:   "bigint",
			dbType:    DBTypeTiDB,
			lang:      "zh_cn",
			want:      "字段 'demo_table.demo_column' 类型变更: bigint -> mediumtext.",
			wantLevel: 1,
		},
		{
			name:      "mysql_zh_varchar_shrink",
			sql:       "ALTER TABLE demo_table MODIFY COLUMN demo_column VARCHAR(20) NOT NULL COMMENT '示例字段'",
			oldType:   "varchar(100)",
			dbType:    DBTypeMysql,
			lang:      "zh_cn",
			want:      "字段 'demo_table.demo_column' 类型变更: varchar(100) -> varchar(20).",
			wantLevel: 1,
		},
		{
			name:      "mysql_zh_int_narrow_unsigned",
			sql:       "ALTER TABLE demo_table MODIFY COLUMN demo_column INT UNSIGNED NOT NULL COMMENT '示例字段'",
			oldType:   "bigint",
			dbType:    DBTypeMysql,
			lang:      "zh_cn",
			incLevel:  map[string]uint8{ER_CHANGE_COLUMN_TYPE.String(): 2},
			want:      "字段 'demo_table.demo_column' 类型变更: bigint -> int(11) unsigned.",
			wantLevel: 2,
		},
		{
			name:         "oceanbase_zh_offline_enabled",
			sql:          "ALTER TABLE demo_table MODIFY COLUMN demo_column MEDIUMTEXT NULL COMMENT '示例字段'",
			oldType:      "bigint",
			dbType:       DBTypeOceanBase,
			offlineCheck: true,
			lang:         "zh_cn",
			want:         oceanBaseOfflineDDLMessage,
			wantLevel:    1,
		},
		{
			name:         "oceanbase_en_offline_enabled",
			sql:          "ALTER TABLE demo_table MODIFY COLUMN demo_column MEDIUMTEXT NULL COMMENT '示例字段'",
			oldType:      "bigint",
			dbType:       DBTypeOceanBase,
			offlineCheck: true,
			lang:         "en_us",
			want:         oceanBaseOfflineDDLMessage,
			wantLevel:    1,
		},
		{
			name:         "oceanbase_zh_offline_disabled_keeps_rule_level",
			sql:          "ALTER TABLE demo_table MODIFY COLUMN demo_column MEDIUMTEXT NULL COMMENT '示例字段'",
			oldType:      "bigint",
			dbType:       DBTypeOceanBase,
			offlineCheck: false,
			lang:         "zh_cn",
			incLevel:     map[string]uint8{ER_CHANGE_COLUMN_TYPE.String(): 2},
			want:         oceanBaseOfflineDDLMessage,
			wantLevel:    2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmt, err := parser.New().ParseOneStmt(tc.sql, "", "")
			if err != nil {
				t.Fatal(err)
			}
			s := newAuditTestSession(tc.dbType, tc.lang, tc.offlineCheck, tc.incLevel)
			s.inc.CheckColumnTypeChange = true
			table := &TableInfo{
				Schema:    "demo_db",
				Name:      "demo_table",
				Collation: "utf8mb4_general_ci",
				Fields:    []FieldInfo{{Field: "demo_column", Type: tc.oldType}},
			}
			s.checkModifyColumn(table, stmt.(*ast.AlterTableStmt).Specs[0])
			if got := s.myRecord.Buf.String(); got != tc.want+"\n" {
				t.Errorf("message: got %q, want %q", got, tc.want+"\n")
			}
			if s.myRecord.ErrLevel != tc.wantLevel {
				t.Errorf("ErrLevel: got %d, want %d", s.myRecord.ErrLevel, tc.wantLevel)
			}
		})
	}
}

func TestSharedDDLAuditMessages(t *testing.T) {
	type targetCase struct {
		name   string
		dbType int
		lang   string
	}
	targets := []targetCase{
		{name: "mysql_zh", dbType: DBTypeMysql, lang: "zh_cn"},
		{name: "mysql_en", dbType: DBTypeMysql, lang: "en_us"},
		{name: "tidb_zh", dbType: DBTypeTiDB, lang: "zh_cn"},
		{name: "tidb_en", dbType: DBTypeTiDB, lang: "en_us"},
		{name: "oceanbase_zh", dbType: DBTypeOceanBase, lang: "zh_cn"},
		{name: "oceanbase_en", dbType: DBTypeOceanBase, lang: "en_us"},
	}

	tests := []struct {
		name      string
		code      ErrorCode
		customLvl uint8
		zh        string
		en        string
		run       func(t *testing.T, s *session)
		obLevel   uint8 // 0 means same as customLvl
	}{
		{
			name:      "drop_table",
			code:      ER_CANT_DROP_TABLE,
			customLvl: 2,
			zh:        "禁止删除表 'demo_table'。",
			en:        "DROP TABLE is not allowed for table 'demo_table'.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("DROP TABLE demo_table", "", "")
				if err != nil {
					t.Fatal(err)
				}
				s.checkDropTable(stmt.(*ast.DropTableStmt), "")
			},
		},
		{
			name:      "truncate_table",
			code:      ER_CANT_TRUNCATE_TABLE,
			customLvl: 2,
			zh:        "禁止清空表 'demo_table'。",
			en:        "TRUNCATE TABLE is not allowed for table 'demo_table'.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("TRUNCATE TABLE demo_table", "", "")
				if err != nil {
					t.Fatal(err)
				}
				s.checkTruncateTable(stmt.(*ast.TruncateTableStmt), "")
			},
		},
		{
			name:      "create_table_charset",
			code:      ER_TABLE_CHARSET_MUST_NULL,
			customLvl: 2,
			zh:        "禁止为表 'demo_table' 指定字符集。",
			en:        "Setting the character set for table 'demo_table' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt(
					"CREATE TABLE demo_table (id INT NOT NULL PRIMARY KEY COMMENT 'id') CHARACTER SET utf8mb4 COMMENT 'demo'",
					"", "",
				)
				if err != nil {
					t.Fatal(err)
				}
				s.checkCreateTable(stmt.(*ast.CreateTableStmt), "")
			},
		},
		{
			name:      "create_table_collation",
			code:      ErrTableCollationNotSupport,
			customLvl: 2,
			zh:        "禁止为表 'demo_table' 指定排序规则。",
			en:        "Setting the collation for table 'demo_table' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt(
					"CREATE TABLE demo_table (id INT NOT NULL PRIMARY KEY COMMENT 'id') COLLATE utf8mb4_general_ci COMMENT 'demo'",
					"", "",
				)
				if err != nil {
					t.Fatal(err)
				}
				s.checkCreateTable(stmt.(*ast.CreateTableStmt), "")
			},
		},
		{
			name:      "alter_table_charset",
			code:      ER_TABLE_CHARSET_MUST_NULL,
			customLvl: 2,
			zh:        "禁止为表 'demo_table' 指定字符集。",
			en:        "Setting the character set for table 'demo_table' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("ALTER TABLE demo_table CHARACTER SET utf8mb4", "", "")
				if err != nil {
					t.Fatal(err)
				}
				s.checkTableOptions(stmt.(*ast.AlterTableStmt).Specs[0].Options, "demo_table", false)
			},
		},
		{
			name:      "alter_table_collation",
			code:      ErrTableCollationNotSupport,
			customLvl: 2,
			zh:        "禁止为表 'demo_table' 指定排序规则。",
			en:        "Setting the collation for table 'demo_table' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("ALTER TABLE demo_table COLLATE utf8mb4_general_ci", "", "")
				if err != nil {
					t.Fatal(err)
				}
				s.checkTableOptions(stmt.(*ast.AlterTableStmt).Specs[0].Options, "demo_table", false)
			},
		},
		{
			name:      "modify_column_position",
			code:      ErCantChangeColumnPosition,
			customLvl: 2,
			zh:        "禁止指定列 'demo_table.col_a' 的位置。",
			en:        "Specifying the position of column 'demo_table.col_a' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("ALTER TABLE demo_table MODIFY COLUMN col_a INT NOT NULL AFTER col_b", "", "")
				if err != nil {
					t.Fatal(err)
				}
				table := &TableInfo{
					Schema: "demo_db",
					Name:   "demo_table",
					Fields: []FieldInfo{
						{Field: "col_a", Type: "int(11)"},
						{Field: "col_b", Type: "int(11)"},
					},
				}
				s.checkModifyColumn(table, stmt.(*ast.AlterTableStmt).Specs[0])
			},
		},
		{
			name:      "change_column_position",
			code:      ErCantChangeColumnPosition,
			customLvl: 2,
			zh:        "禁止指定列 'demo_table.col_c' 的位置。",
			en:        "Specifying the position of column 'demo_table.col_c' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("ALTER TABLE demo_table CHANGE COLUMN col_a col_c INT NOT NULL AFTER col_b", "", "")
				if err != nil {
					t.Fatal(err)
				}
				table := &TableInfo{
					Schema: "demo_db",
					Name:   "demo_table",
					Fields: []FieldInfo{
						{Field: "col_a", Type: "int(11)"},
						{Field: "col_b", Type: "int(11)"},
					},
				}
				s.checkModifyColumn(table, stmt.(*ast.AlterTableStmt).Specs[0])
			},
		},
		{
			name:      "add_column_position",
			code:      ErCantChangeColumnPosition,
			customLvl: 2,
			obLevel:   1, // OceanBase add-column-with-position reports ER_CANT_ADD_STORED_GENERATED_COLUMN (default level 1)
			zh:        "禁止指定列 'demo_table.col_c' 的位置。",
			en:        "Specifying the position of column 'demo_table.col_c' is not allowed.",
			run: func(t *testing.T, s *session) {
				stmt, err := parser.New().ParseOneStmt("ALTER TABLE demo_table ADD COLUMN col_c INT NOT NULL AFTER col_a", "", "")
				if err != nil {
					t.Fatal(err)
				}
				table := &TableInfo{
					Schema: "demo_db",
					Name:   "demo_table",
					Fields: []FieldInfo{
						{Field: "col_a", Type: "int(11)"},
						{Field: "col_b", Type: "int(11)"},
					},
				}
				s.checkAddColumn(table, stmt.(*ast.AlterTableStmt).Specs[0])
			},
		},
	}

	for _, tc := range tests {
		for _, target := range targets {
			t.Run(tc.name+"/"+target.name, func(t *testing.T) {
				incLevel := map[string]uint8{tc.code.String(): tc.customLvl}
				s := newAuditTestSession(target.dbType, target.lang, target.dbType == DBTypeOceanBase, incLevel)
				s.inc.CheckColumnPositionChange = true
				tc.run(t, s)

				wantMsg := tc.zh
				if target.lang == "en_us" {
					wantMsg = tc.en
				}
				wantLevel := tc.customLvl
				if target.dbType == DBTypeOceanBase {
					wantMsg = oceanBaseOfflineDDLMessage
					if tc.obLevel > 0 {
						wantLevel = tc.obLevel
					}
				}
				if got := s.myRecord.Buf.String(); got != wantMsg+"\n" {
					t.Errorf("message: got %q, want %q", got, wantMsg+"\n")
				}
				if s.myRecord.ErrLevel != wantLevel {
					t.Errorf("ErrLevel: got %d, want %d", s.myRecord.ErrLevel, wantLevel)
				}
			})
		}
	}
}
