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
		name         string
		dbType       int
		skipOSC      bool
		checkOffline bool
		oscOn        bool
		ghostOn      bool
		minTableMB   uint
		tableSizeMB  uint
		expectUse    bool
	}{
		{
			name:         "ob skip tools when osc enabled",
			dbType:       DBTypeOceanBase,
			skipOSC:      true,
			checkOffline: true,
			oscOn:        true,
			minTableMB:   16,
			tableSizeMB:  64,
			expectUse:    false,
		},
		{
			name:         "ob skip tools when ghost enabled",
			dbType:       DBTypeOceanBase,
			skipOSC:      true,
			checkOffline: false,
			ghostOn:      true,
			minTableMB:   0,
			tableSizeMB:  1,
			expectUse:    false,
		},
		{
			name:         "ob keep original logic when switch disabled",
			dbType:       DBTypeOceanBase,
			skipOSC:      false,
			checkOffline: true,
			oscOn:        true,
			minTableMB:   16,
			tableSizeMB:  32,
			expectUse:    true,
		},
		{
			name:         "mysql unaffected by ob switch",
			dbType:       DBTypeMysql,
			skipOSC:      true,
			checkOffline: true,
			oscOn:        true,
			minTableMB:   16,
			tableSizeMB:  32,
			expectUse:    true,
		},
		{
			name:        "table below threshold should not use osc",
			dbType:      DBTypeMysql,
			oscOn:       true,
			minTableMB:  16,
			tableSizeMB: 8,
			expectUse:   false,
		},
	}

	for _, tt := range tests {
		s := &session{
			myRecord: &Record{},
			dbType:   tt.dbType,
			inc: config.Inc{
				ObOnlineDDLSkipOsc: tt.skipOSC,
				CheckOfflineDDL:    tt.checkOffline,
			},
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
