package session

import (
	"strings"

	"github.com/hanchuanchuan/goInception/ast"
	"github.com/hanchuanchuan/goInception/model"
)

// maybeSkipOscForOBOnlineDDL 仅在已开启工具链路时单向关闭 useOsc。
// 该逻辑与 check_offline_ddl 解耦,只依赖 OceanBase 白名单在线判定结果。
func (s *session) maybeSkipOscForOBOnlineDDL(node *ast.AlterTableStmt, t *TableInfo) {
	if !s.myRecord.useOsc || s.dbType != DBTypeOceanBase || !s.inc.ObOnlineDDLSkipOsc {
		return
	}
	if s.isOBAlterOnlineByWhitelist(node, t) {
		s.myRecord.useOsc = false
	}
}

// isOBAlterOnlineByWhitelist 按保守白名单判定 ALTER TABLE 是否为在线DDL。
// 仅当全部子句都命中在线白名单时返回 true；未知或无法稳定判定一律返回 false。
func (s *session) isOBAlterOnlineByWhitelist(node *ast.AlterTableStmt, t *TableInfo) bool {
	if node == nil || len(node.Specs) == 0 {
		return false
	}

	for _, spec := range node.Specs {
		if !isOB4AlterSpecOnlineByWhitelist(spec, t) {
			return false
		}
	}
	return true
}

func isOB4AlterSpecOnlineByWhitelist(spec *ast.AlterTableSpec, t *TableInfo) bool {
	if spec == nil {
		return false
	}

	switch spec.Tp {
	case ast.AlterTableAddConstraint:
		return isOBOnlineIndexConstraint(spec.Constraint)
	case ast.AlterTableDropIndex, ast.AlterTableRenameIndex:
		return true
	case ast.AlterTableAddColumns:
		return isOBAddColumnsOnline(spec)
	case ast.AlterTableModifyColumn:
		return isOBModifyColumnOnline(spec, t, false)
	case ast.AlterTableChangeColumn:
		return isOBModifyColumnOnline(spec, t, true)
	case ast.AlterTableRenameColumn, ast.AlterTableAlterColumn, ast.AlterTableRenameTable:
		return true
	case ast.AlterTableOption:
		return isOBAlterTableOptionsOnline(spec.Options)
	case ast.AlterTableAddPartitions:
		return true
	case ast.AlterTablePartition:
		return isOBAlterAutoPartitionAttrsOnly(spec.Partition)
	default:
		return false
	}
}

func isOBOnlineIndexConstraint(c *ast.Constraint) bool {
	if c == nil {
		return false
	}

	switch c.Tp {
	case ast.ConstraintKey, ast.ConstraintIndex,
		ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
		return true
	default:
		return false
	}
}

func isOBAddColumnsOnline(spec *ast.AlterTableSpec) bool {
	if spec == nil || len(spec.NewColumns) == 0 {
		return false
	}
	if spec.Position != nil && spec.Position.Tp != ast.ColumnPositionNone {
		return false
	}

	for _, col := range spec.NewColumns {
		if col == nil {
			return false
		}

		for _, opt := range col.Options {
			switch opt.Tp {
			case ast.ColumnOptionNoOption,
				ast.ColumnOptionNotNull,
				ast.ColumnOptionNull,
				ast.ColumnOptionDefaultValue,
				ast.ColumnOptionOnUpdate,
				ast.ColumnOptionComment,
				ast.ColumnOptionCollate:
				continue
			case ast.ColumnOptionGenerated:
				if opt.Stored {
					return false
				}
			default:
				return false
			}
		}
	}
	return true
}

func isOBModifyColumnOnline(spec *ast.AlterTableSpec, t *TableInfo, useOldColumnName bool) bool {
	if spec == nil || len(spec.NewColumns) != 1 {
		return false
	}
	if spec.Position != nil && spec.Position.Tp != ast.ColumnPositionNone {
		return false
	}

	newCol := spec.NewColumns[0]
	if newCol == nil {
		return false
	}

	oldName := newCol.Name.Name.O
	if useOldColumnName {
		if spec.OldColumnName == nil {
			return false
		}
		oldName = spec.OldColumnName.Name.O
	}

	oldField := findFieldByName(t, oldName)
	if oldField == nil {
		return false
	}

	for _, opt := range newCol.Options {
		switch opt.Tp {
		case ast.ColumnOptionNoOption,
			ast.ColumnOptionNotNull,
			ast.ColumnOptionNull,
			ast.ColumnOptionDefaultValue,
			ast.ColumnOptionOnUpdate,
			ast.ColumnOptionComment,
			ast.ColumnOptionCollate:
			continue
		default:
			return false
		}
	}

	return isOBTypeLengthOrPrecisionGrowOnly(oldField.Type, newCol.Tp.String())
}

func isOBAlterTableOptionsOnline(options []*ast.TableOption) bool {
	if len(options) == 0 {
		return false
	}

	for _, opt := range options {
		if opt == nil {
			return false
		}

		switch opt.Tp {
		case ast.TableOptionAutoIncrement,
			ast.TableOptionRowFormat,
			ast.TableOptionCompression,
			ast.TableOptionBlockSize:
			continue
		default:
			return false
		}
	}
	return true
}

// 自动分区属性调整仅放行“未携带显式分区定义”的场景。
func isOBAlterAutoPartitionAttrsOnly(partition *ast.PartitionOptions) bool {
	if partition == nil {
		return false
	}
	return partition.Tp == model.PartitionTypeRange &&
		partition.Expr == nil &&
		len(partition.ColumnNames) == 0 &&
		len(partition.Definitions) == 0 &&
		partition.Sub == nil
}

func findFieldByName(t *TableInfo, name string) *FieldInfo {
	if t == nil || name == "" {
		return nil
	}
	for i := range t.Fields {
		if !t.Fields[i].IsDeleted && strings.EqualFold(t.Fields[i].Field, name) {
			return &t.Fields[i]
		}
	}
	return nil
}

func isOBTypeLengthOrPrecisionGrowOnly(oldType string, newType string) bool {
	oldType = strings.ToLower(strings.TrimSpace(oldType))
	newType = strings.ToLower(strings.TrimSpace(newType))
	if oldType == "" || newType == "" {
		return false
	}

	if GetDataTypeBase(oldType) != GetDataTypeBase(newType) {
		return false
	}

	oldLen := GetDataTypeLength(oldType)
	newLen := GetDataTypeLength(newType)

	if len(oldLen) != len(newLen) {
		return false
	}
	if len(oldLen) == 0 {
		return false
	}
	if oldLen[0] == -1 {
		return newLen[0] == -1
	}

	for i := range oldLen {
		if newLen[i] < oldLen[i] {
			return false
		}
	}
	return true
}
