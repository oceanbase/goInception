// Copyright 2013 The ql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSES/QL-LICENSE file.

package session

import (
	"fmt"
	"github.com/hanchuanchuan/goInception/ast"
	"strconv"
	"strings"

	"github.com/hanchuanchuan/goInception/model"
)

func (s *session) checkPartitionTruncate(t *TableInfo, parts []model.CIStr) {
	for _, part := range parts {
		found := false
		for _, oldPart := range t.Partitions {
			if strings.EqualFold(part.String(), oldPart.PartName) {
				found = true
				break
			}
		}
		if found && s.dbType == DBTypeOceanBase && s.inc.CheckOfflineDDL {
			if s.dbVersion == 3 {
				// Do Nothing, It's OnLine DDL under 3.x
			} else if s.dbVersion > 3 {
				s.appendErrorNo(ER_CANT_TRUNCATE_PARTITION)
				break
			}
		}

		if !found {
			s.appendErrorNo(ErrPartitionNotExisted, part.String())
		}
	}
}

func (s *session) checkAlterPartitionRule(t *TableInfo, opts *ast.PartitionOptions) {
	if opts == nil {
		return
	}

	if s.dbType == DBTypeOceanBase && s.inc.CheckOfflineDDL {
		s.appendErrorNo(ER_CANT_ALTER_PARTITION_RULE)
	}
}

func (s *session) checkCharsetChange(charset string, table string) bool {
	if s.inc.CheckOfflineDDL {
		s.appendWarningMessage(fmt.Sprintf("Can't change chartset of table '%s' to '%s'.", table, charset))
		return true
	}
	return false
}

func (s *session) checkAddPrimaryKey(t *TableInfo, c *ast.AlterTableSpec) {
	if s.inc.CheckOfflineDDL && s.dbType == DBTypeOceanBase {
		var primaryNames []string
		for _, key := range c.Constraint.Keys {
			primaryNames = append(primaryNames, key.Column.Name.O)
		}
		s.appendErrorNo(ER_CANT_ADD_PRIMARY_KEY)
	}
}

// compareOBVersion 比较两个版本号，返回：
// 1 表示 v1 > v2
// 0 表示 v1 == v2
// -1 表示 v1 < v2
func (s *session) compareOBVersion(v1, v2 string) int {
	segments1 := parseSegments(v1)
	segments2 := parseSegments(v2)

	// 统一长度，不足补0
	maxLen := max(len(segments1), len(segments2))
	segments1 = fillZeros(segments1, maxLen)
	segments2 = fillZeros(segments2, maxLen)

	// 逐个段比较
	for i := 0; i < maxLen; i++ {
		if segments1[i] > segments2[i] {
			return 1
		} else if segments1[i] < segments2[i] {
			return -1
		}
	}
	return 0
}

// 解析版本号为整数数组
func parseSegments(v string) []int {
	parts := strings.Split(v, ".")
	segments := make([]int, len(parts))
	for i, part := range parts {
		num, _ := strconv.Atoi(part) // 格式固定，默认无错误
		segments[i] = num
	}
	return segments
}

// 填充数组到指定长度，补0
func fillZeros(segments []int, targetLen int) []int {
	for len(segments) < targetLen {
		segments = append(segments, 0)
	}
	return segments
}

// 辅助函数：计算最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
