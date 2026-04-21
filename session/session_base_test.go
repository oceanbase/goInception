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
	"testing"

	"github.com/hanchuanchuan/goInception/ast"
	"github.com/hanchuanchuan/goInception/config"
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

func Test_checkAlterUseOsc(t *testing.T) {
	tests := []struct {
		name        string
		oscOn       bool
		ghostOn     bool
		minTableMB  uint
		tableSizeMB uint
		expectUse   bool
	}{
		{
			name:        "osc enabled and over threshold",
			oscOn:       true,
			minTableMB:  16,
			tableSizeMB: 64,
			expectUse:   true,
		},
		{
			name:        "ghost enabled and over threshold",
			ghostOn:     true,
			minTableMB:  8,
			tableSizeMB: 16,
			expectUse:   true,
		},
		{
			name:        "both tools disabled",
			minTableMB:  16,
			tableSizeMB: 64,
			expectUse:   false,
		},
		{
			name:        "table below threshold should not use osc",
			oscOn:       true,
			minTableMB:  16,
			tableSizeMB: 8,
			expectUse:   false,
		},
		{
			name:        "ghost also respects threshold",
			ghostOn:     true,
			minTableMB:  16,
			tableSizeMB: 8,
			expectUse:   false,
		},
	}

	for _, tt := range tests {
		s := &session{
			myRecord: &Record{},
			osc: config.Osc{
				OscOn:           tt.oscOn,
				OscMinTableSize: tt.minTableMB,
			},
			ghost: config.Ghost{
				GhostOn: tt.ghostOn,
			},
		}
		s.checkAlterUseOsc(&TableInfo{TableSize: tt.tableSizeMB})
		if s.myRecord.useOsc != tt.expectUse {
			t.Fatalf("%s: useOsc=%v, expected=%v", tt.name, s.myRecord.useOsc, tt.expectUse)
		}
	}
}

func TestOB4AlterOnlineClassifier_DocOnlineCases(t *testing.T) {
	table := newOBClassifierTable()
	s := &session{}

	tests := []string{
		"alter table t1 add index idx_c1(c1)",
		"alter table t1 drop index idx_c1",
		"alter table t1 rename index idx_old to idx_new",
		"alter table t1 add column c4 int null",
		"alter table t1 add column c5 int generated always as (id + 1) virtual",
		"alter table t1 modify column c1 varchar(32) null",
		"alter table t1 alter column c1 set default 'abc'",
		"alter table t1 alter column c1 drop default",
		"alter table t1 rename column c1 to c1_new",
		"alter table t1 modify column c2 decimal(12,2) default null",
		"alter table t1 rename to t2",
		"alter table t1 row_format = dynamic",
		"alter table t1 auto_increment = 100",
		"alter table t1 add partition (partition p202601 values less than (202602))",
	}

	for _, sql := range tests {
		stmt := parseAlterStmtForTest(t, sql)
		if !s.isOBAlterOnlineByWhitelist(stmt, table) {
			t.Fatalf("expect online but got offline: %s", sql)
		}
	}
}

func TestOB4AlterOnlineClassifier_DocOfflineCases(t *testing.T) {
	table := newOBClassifierTable()
	s := &session{}

	tests := []string{
		"alter table t1 add column c6 int first",
		"alter table t1 add column c7 int auto_increment",
		"alter table t1 modify column c3 int auto_increment",
		"alter table t1 modify column c1 varchar(10)",
		"alter table t1 modify column c1 int",
		"alter table t1 modify column c1 varchar(32) first",
		"alter table t1 add primary key(id)",
		"alter table t1 drop primary key",
		"alter table t1 add column c8 int generated always as (id + 1) stored",
		"alter table t1 drop column c1",
		"alter table t1 convert to character set utf8mb4 collate utf8mb4_bin",
		"alter table t1 drop partition p202601",
		"alter table t1 truncate partition p202601",
		"alter table t1 partition by key(id) partitions 8",
		"alter table t1 add constraint fk_t1 foreign key (c3) references t2(id)",
	}

	for _, sql := range tests {
		stmt := parseAlterStmtForTest(t, sql)
		if s.isOBAlterOnlineByWhitelist(stmt, table) {
			t.Fatalf("expect offline but got online: %s", sql)
		}
	}
}

func TestOB4AlterOnlineClassifier_ConditionalCases(t *testing.T) {
	table := newOBClassifierTable()
	s := &session{}

	tests := []struct {
		sql    string
		expect bool
	}{
		{
			sql:    "alter table t1 add index idx_c1(c1), rename column c1 to c1_new",
			expect: true,
		},
		{
			sql:    "alter table t1 add index idx_c1(c1), drop column c1",
			expect: false,
		},
	}

	for _, tt := range tests {
		stmt := parseAlterStmtForTest(t, tt.sql)
		got := s.isOBAlterOnlineByWhitelist(stmt, table)
		if got != tt.expect {
			t.Fatalf("sql=%s, got=%v, expect=%v", tt.sql, got, tt.expect)
		}
	}
}

func TestOB4AlterOnlineClassifier_UnknownAsOffline(t *testing.T) {
	table := newOBClassifierTable()
	s := &session{}

	tests := []string{
		"alter table t1 force",
		"alter table t1 alter index idx1 invisible",
	}

	for _, sql := range tests {
		stmt := parseAlterStmtForTest(t, sql)
		if s.isOBAlterOnlineByWhitelist(stmt, table) {
			t.Fatalf("unknown spec should fallback offline: %s", sql)
		}
	}
}

func TestOB4AlterOnlineClassifier_AddConstraintConservative(t *testing.T) {
	if isOBOnlineIndexConstraint(&ast.Constraint{Tp: ast.ConstraintForeignKey}) {
		t.Fatalf("foreign key must not be treated as online index constraint")
	}
	if isOBOnlineIndexConstraint(&ast.Constraint{Tp: ast.ConstraintCheck}) {
		t.Fatalf("check constraint must not be treated as online index constraint")
	}
	if !isOBOnlineIndexConstraint(&ast.Constraint{Tp: ast.ConstraintIndex}) {
		t.Fatalf("plain index constraint should still be online")
	}
}

func TestCheckAlterUseOsc_OBSkipOnlyWhenOnline(t *testing.T) {
	table := newOBClassifierTable()
	table.TableSize = 64

	s := &session{
		myRecord: &Record{},
		dbType:   DBTypeOceanBase,
		inc: config.Inc{
			ObOnlineDDLSkipOsc: true,
		},
		osc: config.Osc{
			OscOn:           true,
			OscMinTableSize: 16,
		},
	}

	onlineStmt := parseAlterStmtForTest(t, "alter table t1 add column c4 int null")
	s.checkAlterUseOsc(table)
	if !s.myRecord.useOsc {
		t.Fatalf("base osc decision should be true before online skip")
	}
	s.maybeSkipOscForOBOnlineDDL(onlineStmt, table)
	if s.myRecord.useOsc {
		t.Fatalf("online alter should skip osc")
	}

	offlineStmt := parseAlterStmtForTest(t, "alter table t1 drop column c1")
	s.checkAlterUseOsc(table)
	if !s.myRecord.useOsc {
		t.Fatalf("base osc decision should be true before offline classification")
	}
	s.maybeSkipOscForOBOnlineDDL(offlineStmt, table)
	if !s.myRecord.useOsc {
		t.Fatalf("offline alter should keep osc")
	}
}

func TestCheckAlterUseOsc_OBSkipIndependentFromCheckOfflineDDL(t *testing.T) {
	table := newOBClassifierTable()
	table.TableSize = 64
	stmt := parseAlterStmtForTest(t, "alter table t1 add column c4 int null")

	for _, checkOffline := range []bool{false, true} {
		s := &session{
			myRecord: &Record{},
			dbType:   DBTypeOceanBase,
			inc: config.Inc{
				ObOnlineDDLSkipOsc: true,
				CheckOfflineDDL:    checkOffline,
			},
			osc: config.Osc{
				OscOn:           true,
				OscMinTableSize: 16,
			},
		}

		s.checkAlterUseOsc(table)
		s.maybeSkipOscForOBOnlineDDL(stmt, table)
		if s.myRecord.useOsc {
			t.Fatalf("check_offline_ddl=%v should not affect online skip result", checkOffline)
		}
	}
}

func TestCheckAlterUseOsc_OBSkipNoSideEffectOnErrorLevel(t *testing.T) {
	table := newOBClassifierTable()
	table.TableSize = 64

	record := &Record{
		ErrLevel:     1,
		ErrorMessage: "warn",
	}

	s := &session{
		myRecord: record,
		dbType:   DBTypeOceanBase,
		inc: config.Inc{
			ObOnlineDDLSkipOsc: true,
		},
		osc: config.Osc{
			OscOn:           true,
			OscMinTableSize: 16,
		},
	}

	stmt := parseAlterStmtForTest(t, "alter table t1 add column c4 int null")
	s.checkAlterUseOsc(table)
	beforeLevel := s.myRecord.ErrLevel
	beforeMessage := s.myRecord.ErrorMessage

	s.maybeSkipOscForOBOnlineDDL(stmt, table)

	if s.myRecord.ErrLevel != beforeLevel {
		t.Fatalf("err level changed unexpectedly: before=%d after=%d", beforeLevel, s.myRecord.ErrLevel)
	}
	if s.myRecord.ErrorMessage != beforeMessage {
		t.Fatalf("error message changed unexpectedly: before=%q after=%q", beforeMessage, s.myRecord.ErrorMessage)
	}
}

func parseAlterStmtForTest(t *testing.T, sql string) *ast.AlterTableStmt {
	t.Helper()

	p := parser.New()
	sessionVars := variable.NewSessionVars()
	p.SetSQLMode(sessionVars.SQLMode)

	stmts, _, err := p.Parse(sql, "", "")
	if err != nil {
		t.Fatalf("parse sql failed: %s, err=%v", sql, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expect one stmt, got %d: %s", len(stmts), sql)
	}
	stmt, ok := stmts[0].(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("expect alter table stmt: %s", sql)
	}
	return stmt
}

func newOBClassifierTable() *TableInfo {
	return &TableInfo{
		Name: "t1",
		Fields: []FieldInfo{
			{Field: "id", Type: "int(11)", Null: "NO"},
			{Field: "c1", Type: "varchar(20)", Null: "YES"},
			{Field: "c2", Type: "decimal(10,1)", Null: "YES"},
			{Field: "c3", Type: "int(11)", Null: "YES"},
		},
	}
}
