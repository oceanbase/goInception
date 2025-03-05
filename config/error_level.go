// Code generated by make level. DO NOT EDIT.

package config

type IncLevel struct {
	ErAlterTableOnce                  uint8 `toml:"er_alter_table_once"`
	ErAutoIncrIdWarning               uint8 `toml:"er_auto_incr_id_warning"`
	ErAutoincUnsigned                 uint8 `toml:"er_autoinc_unsigned"`
	ErBadFieldError                   uint8 `toml:"er_bad_field_error"`
	ErBadNullError                    uint8 `toml:"er_bad_null_error"`
	ErBinlogCorrupted                 uint8 `toml:"er_binlog_corrupted"`
	ErBinlogFormatStatement           uint8 `toml:"er_binlog_format_statement"`
	ErBlobCantHaveDefault             uint8 `toml:"er_blob_cant_have_default"`
	ErBlobUsedAsKey                   uint8 `toml:"er_blob_used_as_key"`
	ErCantChangeColumn                uint8 `toml:"er_cant_change_column"`
	ErCantChangeColumnPosition        uint8 `toml:"er_cant_change_column_position"`
	ErCantDropDatabase                uint8 `toml:"er_cant_drop_database"`
	ErCantDropFieldOrKey              uint8 `toml:"er_cant_drop_field_or_key"`
	ErCantDropTable                   uint8 `toml:"er_cant_drop_table"`
	ErCantRemoveAllFields             uint8 `toml:"er_cant_remove_all_fields"`
	ErCantSetCharset                  uint8 `toml:"er_cant_set_charset"`
	ErCantSetCollation                uint8 `toml:"er_cant_set_collation"`
	ErCantSetEngine                   uint8 `toml:"er_cant_set_engine"`
	ErChangeColumnType                uint8 `toml:"er_change_column_type"`
	ErChangeTooMuchRows               uint8 `toml:"er_change_too_much_rows"`
	ErCharToVarcharLen                uint8 `toml:"er_char_to_varchar_len"`
	ErCharsetNotSupport               uint8 `toml:"er_charset_not_support"`
	ErCharsetOnColumn                 uint8 `toml:"er_charset_on_column"`
	ErCollationCharsetMismatch        uint8 `toml:"er_collation_charset_mismatch"`
	ErCollationNotSupport             uint8 `toml:"er_collation_not_support"`
	ErColumnExisted                   uint8 `toml:"er_column_existed"`
	ErColumnHaveNoComment             uint8 `toml:"er_column_have_no_comment"`
	ErColumnNotExisted                uint8 `toml:"er_column_not_existed"`
	ErColumnsMustHaveIndex            uint8 `toml:"er_columns_must_have_index"`
	ErColumnsMustHaveIndexTypeErr     uint8 `toml:"er_columns_must_have_index_type_err"`
	ErConCountError                   uint8 `toml:"er_con_count_error"`
	ErConflictingDeclarations         uint8 `toml:"er_conflicting_declarations"`
	ErDataTooLong                     uint8 `toml:"er_data_too_long"`
	ErDbNotExistedError               uint8 `toml:"er_db_not_existed_error"`
	ErDdlDmlCoexist                   uint8 `toml:"er_ddl_dml_coexist"`
	ErDupFieldname                    uint8 `toml:"er_dup_fieldname"`
	ErDupIndex                        uint8 `toml:"er_dup_index"`
	ErDupKeyname                      uint8 `toml:"er_dup_keyname"`
	ErEndWithCommit                   uint8 `toml:"er_end_with_commit"`
	ErEndWithSemicolon                uint8 `toml:"er_end_with_semicolon"`
	ErEngineNotSupport                uint8 `toml:"er_engine_not_support"`
	ErErrorExistBefore                uint8 `toml:"er_error_exist_before"`
	ErFieldNotInGroupBy               uint8 `toml:"er_field_not_in_group_by"`
	ErFieldSpecifiedTwice             uint8 `toml:"er_field_specified_twice"`
	ErFloatDoubleToDecimal            uint8 `toml:"er_float_double_to_decimal"`
	ErForcingClose                    uint8 `toml:"er_forcing_close"`
	ErForeignKey                      uint8 `toml:"er_foreign_key"`
	ErHaveBegin                       uint8 `toml:"er_have_begin"`
	ErHostname                        uint8 `toml:"er_hostname"`
	ErIdIsUper                        uint8 `toml:"er_id_is_uper"`
	ErIdentUseCustomKeyword           uint8 `toml:"er_ident_use_custom_keyword"`
	ErIdentUseKeyword                 uint8 `toml:"er_ident_use_keyword"`
	ErIdentifierLower                 uint8 `toml:"er_identifier_lower"`
	ErIdentifierUpper                 uint8 `toml:"er_identifier_upper"`
	ErImplicitTypeConversion          uint8 `toml:"er_implicit_type_conversion"`
	ErIncInitErr                      uint8 `toml:"er_inc_init_err"`
	ErInceptionEmptyQuery             uint8 `toml:"er_inception_empty_query"`
	ErIncorrectDatetimeValue          uint8 `toml:"er_incorrect_datetime_value"`
	ErIncorrectGlobalLocalVar         uint8 `toml:"er_incorrect_global_local_var"`
	ErIndexColumnRepeat               uint8 `toml:"er_index_column_repeat"`
	ErIndexNameIdxPrefix              uint8 `toml:"er_index_name_idx_prefix"`
	ErIndexNameUniqPrefix             uint8 `toml:"er_index_name_uniq_prefix"`
	ErIndexNotExisted                 uint8 `toml:"er_index_not_existed"`
	ErIndexUseAlterTable              uint8 `toml:"er_index_use_alter_table"`
	ErInsertTooMuchRows               uint8 `toml:"er_insert_too_much_rows"`
	ErInvalidBackupHostInfo           uint8 `toml:"er_invalid_backup_host_info"`
	ErInvalidCommand                  uint8 `toml:"er_invalid_command"`
	ErInvalidDataType                 uint8 `toml:"er_invalid_data_type"`
	ErInvalidDefault                  uint8 `toml:"er_invalid_default"`
	ErInvalidGroupFuncUse             uint8 `toml:"er_invalid_group_func_use"`
	ErInvalidIdent                    uint8 `toml:"er_invalid_ident"`
	ErInvalidOnUpdate                 uint8 `toml:"er_invalid_on_update"`
	ErJoinNoOnCondition               uint8 `toml:"er_join_no_on_condition"`
	ErJsonTypeSupport                 uint8 `toml:"er_json_type_support"`
	ErKeyColumnDoesNotExits           uint8 `toml:"er_key_column_does_not_exits"`
	ErMariadbRollbackWarn             uint8 `toml:"er_mariadb_rollback_warn"`
	ErMaxColumnCount                  uint8 `toml:"er_max_column_count"`
	ErMaxVarcharLength                uint8 `toml:"er_max_varchar_length"`
	ErMixOfGroupFuncAndFields         uint8 `toml:"er_mix_of_group_func_and_fields"`
	ErMultiplePriKey                  uint8 `toml:"er_multiple_pri_key"`
	ErMustAtLeastOneColumn            uint8 `toml:"er_must_at_least_one_column"`
	ErMustHaveColumns                 uint8 `toml:"er_must_have_columns"`
	ErNetPacketsOutOfOrder            uint8 `toml:"er_net_packets_out_of_order"`
	ErNetReadError                    uint8 `toml:"er_net_read_error"`
	ErNetReadInterrupted              uint8 `toml:"er_net_read_interrupted"`
	ErNetworkReadEventChecksumFailure uint8 `toml:"er_network_read_event_checksum_failure"`
	ErNoDbError                       uint8 `toml:"er_no_db_error"`
	ErNoWhereCondition                uint8 `toml:"er_no_where_condition"`
	ErNonUniqError                    uint8 `toml:"er_non_uniq_error"`
	ErNonUniqTable                    uint8 `toml:"er_non_uniq_table"`
	ErNormalShutdown                  uint8 `toml:"er_normal_shutdown"`
	ErNotAllowedNullable              uint8 `toml:"er_not_allowed_nullable"`
	ErNotFoundMasterStatus            uint8 `toml:"er_not_found_master_status"`
	ErNotFoundTableInfo               uint8 `toml:"er_not_found_table_info"`
	ErNotSupportedAlterOption         uint8 `toml:"er_not_supported_alter_option"`
	ErNotSupportedItemType            uint8 `toml:"er_not_supported_item_type"`
	ErNotSupportedKeyType             uint8 `toml:"er_not_supported_key_type"`
	ErNotSupportedYet                 uint8 `toml:"er_not_supported_yet"`
	ErNotValidPassword                uint8 `toml:"er_not_valid_password"`
	ErOrderyByRand                    uint8 `toml:"er_ordery_by_rand"`
	ErOscKillFailed                   uint8 `toml:"er_osc_kill_failed"`
	ErOutofmemory                     uint8 `toml:"er_outofmemory"`
	ErParseError                      uint8 `toml:"er_parse_error"`
	ErPartitionNotAllowed             uint8 `toml:"er_partition_not_allowed"`
	ErPartitionNotExisted             uint8 `toml:"er_partition_not_existed"`
	ErPkColsNotInt                    uint8 `toml:"er_pk_cols_not_int"`
	ErPkTooManyParts                  uint8 `toml:"er_pk_too_many_parts"`
	ErPrimaryCantHaveNull             uint8 `toml:"er_primary_cant_have_null"`
	ErRemoteExeError                  uint8 `toml:"er_remote_exe_error"`
	ErRemovedSpaces                   uint8 `toml:"er_removed_spaces"`
	ErRepeatConstDefinition           uint8 `toml:"er_repeat_const_definition"`
	ErSameNamePartition               uint8 `toml:"er_same_name_partition"`
	ErSelectOnlyStar                  uint8 `toml:"er_select_only_star"`
	ErSetDataTypeIntBigint            uint8 `toml:"er_set_data_type_int_bigint"`
	ErShutdownComplete                uint8 `toml:"er_shutdown_complete"`
	ErSlaveCorruptEvent               uint8 `toml:"er_slave_corrupt_event"`
	ErSlaveRelayLogWriteFailure       uint8 `toml:"er_slave_relay_log_write_failure"`
	ErSqlInvalidOpType                uint8 `toml:"er_sql_invalid_op_type"`
	ErSqlInvalidSource                uint8 `toml:"er_sql_invalid_source"`
	ErSqlNoOpType                     uint8 `toml:"er_sql_no_op_type"`
	ErSqlNoSource                     uint8 `toml:"er_sql_no_source"`
	ErStartAsBegin                    uint8 `toml:"er_start_as_begin"`
	ErSyntaxError                     uint8 `toml:"er_syntax_error"`
	ErTableCantHandleAutoIncrement    uint8 `toml:"er_table_cant_handle_auto_increment"`
	ErTableCharsetMustNull            uint8 `toml:"er_table_charset_must_null"`
	ErTableCharsetMustUtf8            uint8 `toml:"er_table_charset_must_utf8"`
	ErTableCollationNotSupport        uint8 `toml:"er_table_collation_not_support"`
	ErTableExistsError                uint8 `toml:"er_table_exists_error"`
	ErTableGroupExistsError           uint8 `toml:"er_table_group_exists_error"`
	ErTableGroupNotExistedError       uint8 `toml:"er_table_group_not_existed_error"`
	ErTableMustHaveComment            uint8 `toml:"er_table_must_have_comment"`
	ErTableMustHavePk                 uint8 `toml:"er_table_must_have_pk"`
	ErTableNotExistedError            uint8 `toml:"er_table_not_existed_error"`
	ErTablePrefix                     uint8 `toml:"er_table_prefix"`
	ErTempTableTmpPrefix              uint8 `toml:"er_temp_table_tmp_prefix"`
	ErTextNotNullableError            uint8 `toml:"er_text_not_nullable_error"`
	ErTimestampDefault                uint8 `toml:"er_timestamp_default"`
	ErTooLongBakdbName                uint8 `toml:"er_too_long_bakdb_name"`
	ErTooLongIdent                    uint8 `toml:"er_too_long_ident"`
	ErTooLongIndexComment             uint8 `toml:"er_too_long_index_comment"`
	ErTooLongKey                      uint8 `toml:"er_too_long_key"`
	ErTooManyKeyParts                 uint8 `toml:"er_too_many_key_parts"`
	ErTooManyKeys                     uint8 `toml:"er_too_many_keys"`
	ErTooMuchAutoTimestampCols        uint8 `toml:"er_too_much_auto_timestamp_cols"`
	ErTruncatedWrongValue             uint8 `toml:"er_truncated_wrong_value"`
	ErUdpateTooMuchRows               uint8 `toml:"er_udpate_too_much_rows"`
	ErUnknownCharacterSet             uint8 `toml:"er_unknown_character_set"`
	ErUnknownCharset                  uint8 `toml:"er_unknown_charset"`
	ErUnknownCollation                uint8 `toml:"er_unknown_collation"`
	ErUnknownSystemVariable           uint8 `toml:"er_unknown_system_variable"`
	ErUnknownTable                    uint8 `toml:"er_unknown_table"`
	ErUseEnum                         uint8 `toml:"er_use_enum"`
	ErUseIndexVisibility              uint8 `toml:"er_use_index_visibility"`
	ErUseTextOrBlob                   uint8 `toml:"er_use_text_or_blob"`
	ErUseValueExpr                    uint8 `toml:"er_use_value_expr"`
	ErUsername                        uint8 `toml:"er_username"`
	ErVarcharToTextLen                uint8 `toml:"er_varchar_to_text_len"`
	ErViewSelectClause                uint8 `toml:"er_view_select_clause"`
	ErViewSupport                     uint8 `toml:"er_view_support"`
	ErWithDefaultAddColumn            uint8 `toml:"er_with_default_add_column"`
	ErWithInsertField                 uint8 `toml:"er_with_insert_field"`
	ErWithInsertValues                uint8 `toml:"er_with_insert_values"`
	ErWithLimitCondition              uint8 `toml:"er_with_limit_condition"`
	ErWithOrderbyCondition            uint8 `toml:"er_with_orderby_condition"`
	ErWrongAndExpr                    uint8 `toml:"er_wrong_and_expr"`
	ErWrongArguments                  uint8 `toml:"er_wrong_arguments"`
	ErWrongAutoKey                    uint8 `toml:"er_wrong_auto_key"`
	ErWrongColumnName                 uint8 `toml:"er_wrong_column_name"`
	ErWrongDbName                     uint8 `toml:"er_wrong_db_name"`
	ErWrongKeyColumn                  uint8 `toml:"er_wrong_key_column"`
	ErWrongNameForIndex               uint8 `toml:"er_wrong_name_for_index"`
	ErWrongStringLength               uint8 `toml:"er_wrong_string_length"`
	ErWrongSubKey                     uint8 `toml:"er_wrong_sub_key"`
	ErWrongTableName                  uint8 `toml:"er_wrong_table_name"`
	ErWrongUsage                      uint8 `toml:"er_wrong_usage"`
	ErWrongValueCountOnRow            uint8 `toml:"er_wrong_value_count_on_row"`
	ErWrongValueForVar                uint8 `toml:"er_wrong_value_for_var"`
}

var defaultLevel = IncLevel{
	ErAlterTableOnce:                  2,
	ErAutoIncrIdWarning:               2,
	ErAutoincUnsigned:                 2,
	ErBadFieldError:                   2,
	ErBadNullError:                    2,
	ErBinlogCorrupted:                 2,
	ErBinlogFormatStatement:           2,
	ErBlobCantHaveDefault:             2,
	ErBlobUsedAsKey:                   2,
	ErCantChangeColumn:                2,
	ErCantChangeColumnPosition:        2,
	ErCantDropDatabase:                2,
	ErCantDropFieldOrKey:              2,
	ErCantDropTable:                   1,
	ErCantRemoveAllFields:             2,
	ErCantSetCharset:                  2,
	ErCantSetCollation:                2,
	ErCantSetEngine:                   2,
	ErChangeColumnType:                1,
	ErChangeTooMuchRows:               2,
	ErCharToVarcharLen:                2,
	ErCharsetNotSupport:               2,
	ErCharsetOnColumn:                 2,
	ErCollationCharsetMismatch:        2,
	ErCollationNotSupport:             2,
	ErColumnExisted:                   2,
	ErColumnHaveNoComment:             2,
	ErColumnNotExisted:                2,
	ErColumnsMustHaveIndex:            2,
	ErColumnsMustHaveIndexTypeErr:     2,
	ErConCountError:                   2,
	ErConflictingDeclarations:         2,
	ErDataTooLong:                     2,
	ErDbNotExistedError:               2,
	ErDdlDmlCoexist:                   2,
	ErDupFieldname:                    2,
	ErDupIndex:                        2,
	ErDupKeyname:                      2,
	ErEndWithCommit:                   2,
	ErEndWithSemicolon:                2,
	ErEngineNotSupport:                2,
	ErErrorExistBefore:                2,
	ErFieldNotInGroupBy:               2,
	ErFieldSpecifiedTwice:             2,
	ErFloatDoubleToDecimal:            2,
	ErForcingClose:                    2,
	ErForeignKey:                      2,
	ErHaveBegin:                       2,
	ErHostname:                        2,
	ErIdIsUper:                        2,
	ErIdentUseCustomKeyword:           2,
	ErIdentUseKeyword:                 2,
	ErIdentifierLower:                 2,
	ErIdentifierUpper:                 2,
	ErImplicitTypeConversion:          2,
	ErIncInitErr:                      2,
	ErInceptionEmptyQuery:             2,
	ErIncorrectDatetimeValue:          2,
	ErIncorrectGlobalLocalVar:         2,
	ErIndexColumnRepeat:               2,
	ErIndexNameIdxPrefix:              2,
	ErIndexNameUniqPrefix:             2,
	ErIndexNotExisted:                 2,
	ErIndexUseAlterTable:              2,
	ErInsertTooMuchRows:               2,
	ErInvalidBackupHostInfo:           2,
	ErInvalidCommand:                  2,
	ErInvalidDataType:                 2,
	ErInvalidDefault:                  2,
	ErInvalidGroupFuncUse:             2,
	ErInvalidIdent:                    2,
	ErInvalidOnUpdate:                 2,
	ErJoinNoOnCondition:               2,
	ErJsonTypeSupport:                 2,
	ErKeyColumnDoesNotExits:           2,
	ErMariadbRollbackWarn:             2,
	ErMaxColumnCount:                  2,
	ErMaxVarcharLength:                2,
	ErMixOfGroupFuncAndFields:         2,
	ErMultiplePriKey:                  2,
	ErMustAtLeastOneColumn:            2,
	ErMustHaveColumns:                 2,
	ErNetPacketsOutOfOrder:            2,
	ErNetReadError:                    2,
	ErNetReadInterrupted:              2,
	ErNetworkReadEventChecksumFailure: 2,
	ErNoDbError:                       2,
	ErNoWhereCondition:                2,
	ErNonUniqError:                    2,
	ErNonUniqTable:                    2,
	ErNormalShutdown:                  2,
	ErNotAllowedNullable:              2,
	ErNotFoundMasterStatus:            2,
	ErNotFoundTableInfo:               2,
	ErNotSupportedAlterOption:         2,
	ErNotSupportedItemType:            2,
	ErNotSupportedKeyType:             2,
	ErNotSupportedYet:                 2,
	ErNotValidPassword:                2,
	ErOrderyByRand:                    2,
	ErOscKillFailed:                   2,
	ErOutofmemory:                     2,
	ErParseError:                      2,
	ErPartitionNotAllowed:             2,
	ErPartitionNotExisted:             2,
	ErPkColsNotInt:                    2,
	ErPkTooManyParts:                  2,
	ErPrimaryCantHaveNull:             2,
	ErRemoteExeError:                  2,
	ErRemovedSpaces:                   2,
	ErRepeatConstDefinition:           2,
	ErSameNamePartition:               2,
	ErSelectOnlyStar:                  2,
	ErSetDataTypeIntBigint:            2,
	ErShutdownComplete:                2,
	ErSlaveCorruptEvent:               2,
	ErSlaveRelayLogWriteFailure:       2,
	ErSqlInvalidOpType:                2,
	ErSqlInvalidSource:                2,
	ErSqlNoOpType:                     2,
	ErSqlNoSource:                     2,
	ErStartAsBegin:                    2,
	ErSyntaxError:                     2,
	ErTableCantHandleAutoIncrement:    2,
	ErTableCharsetMustNull:            1,
	ErTableCharsetMustUtf8:            2,
	ErTableCollationNotSupport:        1,
	ErTableExistsError:                2,
	ErTableGroupExistsError:           2,
	ErTableGroupNotExistedError:       2,
	ErTableMustHaveComment:            2,
	ErTableMustHavePk:                 2,
	ErTableNotExistedError:            2,
	ErTablePrefix:                     2,
	ErTempTableTmpPrefix:              2,
	ErTextNotNullableError:            2,
	ErTimestampDefault:                2,
	ErTooLongBakdbName:                2,
	ErTooLongIdent:                    2,
	ErTooLongIndexComment:             2,
	ErTooLongKey:                      2,
	ErTooManyKeyParts:                 2,
	ErTooManyKeys:                     2,
	ErTooMuchAutoTimestampCols:        2,
	ErTruncatedWrongValue:             2,
	ErUdpateTooMuchRows:               2,
	ErUnknownCharacterSet:             2,
	ErUnknownCharset:                  2,
	ErUnknownCollation:                2,
	ErUnknownSystemVariable:           2,
	ErUnknownTable:                    2,
	ErUseEnum:                         2,
	ErUseIndexVisibility:              2,
	ErUseTextOrBlob:                   2,
	ErUseValueExpr:                    2,
	ErUsername:                        2,
	ErVarcharToTextLen:                2,
	ErViewSelectClause:                2,
	ErViewSupport:                     2,
	ErWithDefaultAddColumn:            2,
	ErWithInsertField:                 2,
	ErWithInsertValues:                2,
	ErWithLimitCondition:              2,
	ErWithOrderbyCondition:            2,
	ErWrongAndExpr:                    2,
	ErWrongArguments:                  2,
	ErWrongAutoKey:                    2,
	ErWrongColumnName:                 2,
	ErWrongDbName:                     2,
	ErWrongKeyColumn:                  2,
	ErWrongNameForIndex:               2,
	ErWrongStringLength:               2,
	ErWrongSubKey:                     2,
	ErWrongTableName:                  2,
	ErWrongUsage:                      2,
	ErWrongValueCountOnRow:            2,
	ErWrongValueForVar:                2,
}
