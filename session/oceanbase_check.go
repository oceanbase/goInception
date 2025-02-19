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
			} else if s.dbVersion == 4 {
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
	segments1 := splitVersion(v1)
	segments2 := splitVersion(v2)
	maxLen := max(len(segments1), len(segments2))

	for i := 0; i < maxLen; i++ {
		s1 := getSegment(segments1, i)
		s2 := getSegment(segments2, i)

		if res := compareSegments(s1, s2); res != 0 {
			return res
		}
	}
	return 0
}

// 分割版本号为分段（按点分割）
func splitVersion(v string) []string {
	return strings.Split(v, ".")
}

// 获取分段，超出长度返回空字符串
func getSegment(segments []string, index int) string {
	if index < len(segments) {
		return segments[index]
	}
	return ""
}

// 比较两个版本段
func compareSegments(s1, s2 string) int {
	// 尝试将段转换为数字
	num1, isNum1 := parseSegment(s1)
	num2, isNum2 := parseSegment(s2)

	// 情况1：两者均为数字 → 数值比较
	if isNum1 && isNum2 {
		return compareInt(num1, num2)
	}

	// 情况2：数字优先级高于非数字
	if isNum1 {
		return 1
	}
	if isNum2 {
		return -1
	}

	// 情况3：均为非数字 → 字符串比较
	return strings.Compare(s1, s2)
}

// 尝试将段解析为数字
func parseSegment(s string) (int, bool) {
	num, err := strconv.Atoi(s)
	if err == nil {
		return num, true
	}
	// 尝试提取前缀数字（如 "3bp10" → 3）
	for i := range s {
		if !isDigit(rune(s[i])) {
			if i > 0 {
				num, err := strconv.Atoi(s[:i])
				if err == nil {
					return num, true
				}
			}
			break
		}
	}
	return 0, false
}

// 判断字符是否为数字
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// 比较两个整数
func compareInt(a, b int) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

// 辅助函数：返回最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
