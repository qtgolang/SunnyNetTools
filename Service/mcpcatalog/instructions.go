package mcpcatalog

// MCPStreamableInstructions 在 MCP initialize 时下发给客户端，约束多行操作的调用方式。
const MCPStreamableInstructions = `SunnyNetTools MCP 批量约定（务必遵守）：
1. 界面「序号」列 = 从 1 开始的显示行号 listIndex（第 66 条 → listIndex:66 或 main_slice offset=65）。禁止把序号当成 theology/rowId 使用。
2. theology 为内部会话主键（整数）；rowId 与其一一对应，rowId = strconv(theology)（十进制字符串，如 theology 66 → rowId "66"）。与 listIndex 通常不相等。多行一次传入 rowIds、theologies 或 listIndexes，禁止循环单条调用。
3. 仅当目标确实只有 1 行时，才使用 theology 或 rowId（二者等价；勿把 listIndex 当 rowId）。
4. 写注释：sunnynet_main_row_note_set 多行时必须 rowIds + note 一次完成；读注释多行用 sunnynet_main_row_note_get 的 rowIds。
5. 读断点状态：main_row_break_get 返回行级 breakMode/isWaiting；get_status.breakMode 为全局拦截状态只读，不可经 MCP 修改（请在应用 UI 切换）。
6. 断点改包：break_sync_request 仅用于拦截上行（不可改 Method）；break_sync_response 仅用于拦截下行；非对应模式直接报错。放行用 continue:true 或 break_continue。
7. 在内存中定位唯一一行：main_search（返回 listIndex/theology）；多条命中时缩小条件。读 HTTP/流正文用 http_get_part、stream_get_part 分块。
8. 改规则后 UI 由各规则页监听 configreload 自动刷新。禁用规则用 config_rule_set_state（id + state:已禁用），禁止用空 rulesJSON 或省略规则行。
9. 按序号查行可先 main_slice 或 listindextotheology；断点等待中 main_slice/session_get_json 可正常调用，放行用 break_continue。`

// batchRowOps 支持 rowId | rowIds | ids 批量主键的 op。
var batchRowOps = map[string]bool{
	"row_theology": true, "main_cells": true, "main_row_note_get": true, "main_row_note_set": true,
}

// batchTheologyRowOps 支持 theology | theologies | rowId | rowIds | ids 的 op。
var batchTheologyRowOps = map[string]bool{
	"break_continue": true, "break_skip_to_response": true,
	"break_sync_request": true, "break_sync_response": true,
	"http_replay": true, "session_get_json": true, "main_row_break_get": true,
	"http_get_part": true, "stream_get_part": true, "stream_get_hex": true,
}

func opUsesBatchIDs(op string) bool {
	return batchRowOps[op] || batchTheologyRowOps[op]
}
