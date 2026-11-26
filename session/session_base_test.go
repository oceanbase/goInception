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

func TestColumnTypeChangeMessage(t *testing.T) {
	const offline = "[OceanBase Offline DDL Check] 离线DDL，需要重整表数据，执行过程会阻塞DML即数据写入风险，请谨慎操作."
	for _, tc := range []struct {
		dbType       int
		offlineCheck bool
		lang         string
		want         string
	}{
		{DBTypeMysql, false, "zh_cn", "字段 'order.actual_conversion' 类型变更: bigint -> mediumtext."},
		{DBTypeMysql, false, "en_us", "Type conversion warning for column 'order.actual_conversion': bigint -> mediumtext."},
		{DBTypeOceanBase, true, "zh_cn", offline},
		{DBTypeOceanBase, true, "en_us", offline},
		{DBTypeOceanBase, false, "zh_cn", offline},
	} {
		stmt, err := parser.New().ParseOneStmt("ALTER TABLE `order` MODIFY COLUMN actual_conversion MEDIUMTEXT NULL COMMENT '实际转化类型ID列表'", "", "")
		if err != nil {
			t.Fatal(err)
		}
		s := &session{
			dbType: tc.dbType, dbVersion: 4, opt: &SourceOptions{},
			myRecord: &Record{Buf: new(bytes.Buffer)}, recordSets: &MyRecordSets{},
		}
		s.inc.EnableNullable = true
		s.inc.CheckColumnTypeChange = true
		s.inc.Lang = tc.lang
		s.inc.CheckOfflineDDL = tc.offlineCheck
		table := &TableInfo{Name: "order", Fields: []FieldInfo{{Field: "actual_conversion", Type: "bigint"}}}
		s.checkModifyColumn(table, stmt.(*ast.AlterTableStmt).Specs[0])
		if got := s.myRecord.Buf.String(); got != tc.want+"\n" {
			t.Errorf("dbType=%d lang=%s: got %q, want %q", tc.dbType, tc.lang, got, tc.want)
		}
	}
}

func TestSharedDDLAuditMessages(t *testing.T) {
	const offline = "[OceanBase Offline DDL Check] 离线DDL，需要重整表数据，执行过程会阻塞DML即数据写入风险，请谨慎操作."
	for _, tc := range []struct {
		code           ErrorCode
		object, zh, en string
	}{
		{ER_CANT_DROP_TABLE, "orders", "禁止删除表 'orders'。", "DROP TABLE is not allowed for table 'orders'."},
		{ER_CANT_TRUNCATE_TABLE, "orders", "禁止清空表 'orders'。", "TRUNCATE TABLE is not allowed for table 'orders'."},
		{ER_TABLE_CHARSET_MUST_NULL, "orders", "禁止为表 'orders' 指定字符集。", "Setting the character set for table 'orders' is not allowed."},
		{ErrTableCollationNotSupport, "orders", "禁止为表 'orders' 指定排序规则。", "Setting the collation for table 'orders' is not allowed."},
		{ErCantChangeColumnPosition, "orders.id", "禁止指定列 'orders.id' 的位置。", "Specifying the position of column 'orders.id' is not allowed."},
	} {
		for _, target := range []struct {
			dbType     int
			lang, want string
		}{
			{DBTypeMysql, "zh_cn", tc.zh}, {DBTypeMysql, "en_us", tc.en},
			{DBTypeOceanBase, "zh_cn", offline}, {DBTypeOceanBase, "en_us", offline},
		} {
			s := &session{dbType: target.dbType, myRecord: &Record{Buf: new(bytes.Buffer)}, recordSets: &MyRecordSets{}}
			s.inc.Lang = target.lang
			s.inc.CheckColumnPositionChange = true
			s.appendErrorNo(tc.code, tc.object)
			if got := s.myRecord.Buf.String(); got != target.want+"\n" {
				t.Errorf("%v dbType=%d lang=%s: got %q, want %q", tc.code, target.dbType, target.lang, got, target.want)
			}
		}
	}
}
