package mcpcatalog

import (
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

var (
	toolInputSchemas map[string]any
	schemaOnce       sync.Once
)

func loadToolInputSchemas() {
	toolInputSchemas = map[string]any{
		"row_theology":           mustInputSchema[BatchRowIDsIn](),
		"main_cells":             mustInputSchema[BatchRowIDsIn](),
		"main_row_note_get":      mustInputSchema[BatchRowIDsIn](),
		"main_row_note_set":      mustInputSchema[BatchRowIDsNoteIn](),
		"main_row_break_get":     mustInputSchema[BatchTheologyRowIn](),
		"break_continue":         mustInputSchema[BatchTheologyRowIn](),
		"break_skip_to_response": mustInputSchema[BatchTheologyRowIn](),
		"break_sync_request":     mustInputSchema[BreakSyncRequestIn](),
		"break_sync_response":    mustInputSchema[BreakSyncResponseIn](),
		"http_replay":            mustInputSchema[HTTPReplayIn](),
		"session_get_json":       mustInputSchema[BatchTheologyRowIn](),
		"http_get_part":          mustInputSchema[HTTPGetPartIn](),
		"stream_get_part":        mustInputSchema[StreamGetPartIn](),
		"stream_get_hex":         mustInputSchema[StreamGetPartIn](),
		"main_delete":            mustInputSchema[MainDeleteIn](),
		"main_delete_except":     mustInputSchema[MainDeleteExceptIn](),
		"main_apply_row_mark":    mustInputSchema[MainApplyRowMarkIn](),
		"main_slice":             mustInputSchema[MainSliceIn](),
		"records_import":         mustInputSchema[RecordsImportIn](),
		"records_export":         mustInputSchema[RecordsExportIn](),
		"session_pack_export":    mustInputSchema[SessionPackExportIn](),
		"stream_count":           mustInputSchema[StreamCountIn](),
		"stream_slice":           mustInputSchema[StreamSliceIn](),
		"engine_start":           mustInputSchema[EngineStartIn](),
		"generate_builtin_code":  mustInputSchema[GenerateBuiltinCodeIn](),
		"engine_apply_advanced":  mustInputSchema[EngineApplyAdvancedIn](),
		"config_set_replace":     mustInputSchema[ConfigSetRulesIn](),
		"config_set_rewrite":     mustInputSchema[ConfigSetRulesIn](),
		"config_set_intercept":   mustInputSchema[ConfigSetRulesIn](),
		"config_set_block":       mustInputSchema[ConfigSetRulesIn](),
		"config_set_host":        mustInputSchema[ConfigSetRulesIn](),
		"config_host_add":        mustInputSchema[HostRuleAddIn](),
		"config_host_delete":     mustInputSchema[HostRuleDeleteIn](),
		"config_host_update":     mustInputSchema[HostRuleUpdateIn](),
		"config_rule_set_state":  mustInputSchema[ConfigRuleSetStateIn](),
		"device_load":            mustInputSchema[DeviceLoadIn](),
		"device_process_add_name": mustInputSchema[DeviceProcessNameIn](),
		"device_process_del_name": mustInputSchema[DeviceProcessNameIn](),
		"device_process_add_pid":  mustInputSchema[DeviceProcessPIDIn](),
		"device_process_del_pid":  mustInputSchema[DeviceProcessPIDIn](),
		"ui_theme_set":           mustInputSchema[UIThemeSetIn](),
		"pb_to_json":             mustInputSchema[PbToJSONIn](),
		"main_search":            mustInputSchema[MainSearchIn](),
		"request_cert_add":       mustInputSchema[RequestCertAddIn](),
		"request_cert_delete":    mustInputSchema[RequestCertIDIn](),
		"request_cert_update":       mustInputSchema[RequestCertUpdateIn](),
		"config_set_proxy_dns":      mustInputSchema[ProxyDnsSetIn](),
		"config_proxy_way_add":        mustInputSchema[ProxyWayAddIn](),
		"config_proxy_way_update_note": mustInputSchema[ProxyWayUpdateNoteIn](),
		"config_proxy_way_delete":   mustInputSchema[ProxyWayIDIn](),
		"config_proxy_way_set_state": mustInputSchema[ProxyWaySetStateIn](),
		"config_proxy_way_update":   mustInputSchema[ProxyWayUpdateIn](),
		"config_set_proxy_roles":    mustInputSchema[ProxyRolesSetIn](),
		"config_set_must_tcp":       mustInputSchema[MustTcpSetIn](),
		"config_set_disable_tcp":    mustInputSchema[BoolToggleIn](),
		"config_set_disable_udp":    mustInputSchema[BoolToggleIn](),
		"config_set_disable_cache":  mustInputSchema[BoolToggleIn](),
		"config_set_limit_request_size": mustInputSchema[LimitRequestSizeIn](),
		"config_set_https_protocol": mustInputSchema[HTTPSProtocolSetIn](),
		"config_set_random_ja3":     mustInputSchema[RandomJa3SetIn](),
		"config_set_http2_fingerprint": mustInputSchema[HTTP2FingerprintSetIn](),
		"config_apply_http2_template":  mustInputSchema[HTTP2TemplateApplyIn](),
		"config_get_http2_template": mustInputSchema[HTTP2TemplateGetIn](),
	}
}

func mustInputSchema[T any]() any {
	s, err := jsonschema.For[T](nil)
	if err != nil {
		panic("mcpcatalog: input schema for " + err.Error())
	}
	return s
}

// ToolInputSchema 返回指定 op 的 JSON Schema；无参数 op 返回 nil（由 SDK 推断为空 object）。
func ToolInputSchema(op string) any {
	schemaOnce.Do(loadToolInputSchemas)
	if s, ok := toolInputSchemas[op]; ok {
		return s
	}
	return nil
}
