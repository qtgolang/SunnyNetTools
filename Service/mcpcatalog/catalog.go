// Package mcpcatalog holds MCP bridge op catalog JSON and Streamable HTTP tool metadata.
package mcpcatalog

import (
	"encoding/json"
	"strings"
)

// BridgeOpCapability 单条桥接能力说明，供 list_supported_ops / MCPListOpsJSON 返回给第三方客户端。
type BridgeOpCapability struct {
	Op          string `json:"op"`
	Description string `json:"description"`
	Args        string `json:"args,omitempty"`
	Returns     string `json:"returns,omitempty"`
}

// bridgeOpCatalog 与 mcpbridge.Host.Invoke 支持的 op 一一对应；顺序即 SupportedBridgeOps。
var bridgeOpCatalog = []BridgeOpCapability{
	{Op: "ping", Description: "【场景】脚本或 MCP 客户端自检：主程序 TCP 桥已监听且能调度到业务 Host。【效果】无业务副作用，仅证明链路通。【勿用】不可替代 get_status（不含引擎状态）。", Args: "无", Returns: "{\"pong\":true} 表示桥接可读且 Host 可调度。"},
	{Op: "get_status", Description: "【场景】任何自动化前先读一眼：引擎是否在跑、监听端口、上次错误、系统代理是否指向本机、捕获区是否显示、主列表行数。【效果】只读快照；breakMode 为当前全局拦截状态只读（不可经 MCP 修改，请在 UI 或规则中配置）。", Args: "无", Returns: "{\"sunnyRunning\",\"sunnyPort\",\"sunnyLastError\",\"systemProxyEnabled\",\"captureVisible\",\"breakMode\",\"mainCount\"}；breakMode：0 关 / 1 上行 / 2 下行（只读）。"},
	{Op: "row_theology", Description: "【场景】校验或读取 rowId 与 theology 的对应关系。【约定】rowId 即 theology 的十进制字符串（rowId=\"66\" ↔ theology=66），无 http-* 等其它格式。【效果】解析 rowId/ids 得到 theology；无效 rowId 时 theology 为 0。", Args: "rowId | rowIds | ids | theology | theologies", Returns: "单条：{\"theology\":number,\"rowId\":string}；多条：{\"items\":[{\"rowId\",\"theology\"},…]}。"},
	{Op: "engine_start", Description: "【场景】要开始抓包、重放或脚本依赖 Sunny 监听端口时调用。【效果】在指定端口启动 Sunny 引擎；可选是否顺带设置系统代理。【前提】端口未被占用；驱动/设备策略仍按 UI 与 device_* 系列单独处理。", Args: "port? useSystemProxy?", Returns: "{\"ok\":true,\"port\":number}；失败抛错（错误信息在异常文本中）。"},
	{Op: "engine_stop", Description: "【场景】抓包结束、释放端口、或关闭应用前希望干净停止。【效果】停止引擎；若引擎内部返回提示则 ok 可能为 false 并带 message。", Args: "无", Returns: "{\"ok\":true} 或 {\"ok\":false,\"message\":string}。"},
	{Op: "capture_hide", Description: "【场景】演示、截图、或自动化时希望主界面隐藏中间抓包列表区域（与 UI「隐藏捕获」同源）。【效果】仅影响主窗口捕获区可见性，不停止引擎。", Args: "无", Returns: "{\"ok\":true}。"},
	{Op: "capture_show", Description: "【场景】在 capture_hide 之后恢复抓包区显示。", Args: "无", Returns: "{\"ok\":true}。"},
	{Op: "system_proxy_enable", Description: "【场景】希望浏览器等走系统代理的流量经过 Sunny（指向当前 Sunny 端口）。【效果】写入系统代理设置。【注意】需用户环境允许改代理；失败会抛错。", Args: "无", Returns: "{\"ok\":true}；失败抛错。"},
	{Op: "system_proxy_disable", Description: "【场景】抓包结束或异常退出后，恢复系统不再指向 Sunny。【效果】取消由本工具设置的系统代理路径。", Args: "无", Returns: "{\"ok\":true}；失败抛错。"},
	{Op: "break_continue", Description: "【批量】多会话放行请一次传 rowIds 或 theologies，禁止逐条调用。【场景】HTTP 断点等待中放行（与 UI「继续」同源）。", Args: "theology | theologies | rowId | rowIds | ids", Returns: "单条：{\"ok\":true}；多条：{\"ok\":true,\"count\":n}。"},
	{Op: "break_continue_all", Description: "【场景】列表里多条请求同时卡在断点，一键全部放行。", Args: "无", Returns: "{\"ok\":true}。"},
	{Op: "break_skip_to_response", Description: "【场景】当前卡在请求阶段断点，希望先放行请求，等响应回来再在响应阶段进入断点（与 UI「跳到响应」类操作同源）。", Args: "theology | theologies | rowId | rowIds | ids", Returns: "单条：{\"ok\":true}；多条：{\"ok\":true,\"count\":n}。"},
	{Op: "break_sync_request", Description: "【前置】须拦截上行（IsWait+请求阶段）；非拦截直接报错不修改。【限制】不可改 Method。【放行】continue/release:true。", Args: "theology|rowIds requestURL? headersJSON? bodyB64? continue? release?", Returns: "{\"ok\",continued,listUpdated}。"},
	{Op: "break_sync_response", Description: "【前置】须拦截下行（IsWait+响应阶段）；非拦截直接报错。【放行】continue/release:true。", Args: "theology|rowIds statusCode? headersJSON? bodyB64? continue? release?", Returns: "同 break_sync_request。"},
	{Op: "http_replay", Description: "【场景】把主列表里已抓到的一条会话按快照再发一遍（HTTP/WS/TCP 等由引擎与快照决定）；用于复现请求、压测单条、或调试脚本。【参数】interceptMode：0 普通重放；1 重放时走请求断点；2 重放时走响应断点（仅单条 HTTP 非 WS 场景与引擎一致）。repeatCount：重复次数（>1 时引擎有并发上限 10）。", Args: "theology|theologies|rowId|rowIds|ids interceptMode? repeatCount?", Returns: "单条：{\"ok\":true}；多条：{\"ok\":true,\"count\":n}。"},
	{Op: "generate_builtin_code", Description: "【场景】已有一条抓到的 HTTP/TCP/WebSocket 会话，要在「Go / C# / Python / cURL / 火山 / 易语言」里生成可粘贴的客户端示例代码（与主界面「生成代码」同一套模板）。UDP 会话不支持，会失败。【语言 language】GoLang、C#、Python、cURL、火山、易语言（须与代码生成器 switch 完全一致）。【模块 module】该语言下的子模板名，例如 Go 侧 net/http、Python 侧 requests；cURL 侧为 Linux Terminal（bash 单引号）/ Windows CMD（cmd 双引号，二进制为 powershell -EncodedCommand）/ Windows PowerShell（curl.exe + 单引号）等，须与 UI 一致。【重要】成功时主程序会把代码写入系统剪贴板，HTTP 返回的 text 为空字符串；text 非空表示错误说明而非代码正文。", Args: "theology|rowId language module", Returns: "{\"text\":string}；成功时 text 常为空（代码在剪贴板）；失败时 text 为错误文案。"},
	{Op: "engine_apply_advanced", Description: "【场景】不打开设置窗口就要改「高级设置」：强制 TCP 规则脚本、上游代理列表、上游 DNS 模式、出口路由、Socks5 认证、最大请求体长度、是否禁用 TCP/UDP 等（与设置页高级选项、SunnyApplyEngineAdvanced 使用同一份 JSON 结构）。【效果】合并进当前运行中的 Sunny 实例。【参数】payload 为整份 engineAdvanced JSON 字符串；manual 与 UI「手动下发」勾选含义一致（布尔，影响是否记为手动快照等内部逻辑）。", Args: "payload（与 UI engineAdvanced 同结构的 JSON 字符串） manual?（bool）", Returns: "{\"ok\":true}；失败抛错。"},
	{Op: "main_count", Description: "【场景】分页拉 main_slice 前需要总条数；或监控列表是否为空。【效果】只读当前过滤后的主列表行数。【协定】故意返回 JSON 对象而非裸整数：许多 MCP 宿主对 tools/call 的 result 做 JSON 解码，裸 number 会校验失败。调用方必须把 result 当对象解析后读字段 total；不要假设 result 的类型是 number，也不要为「拿个数」改用 main_slice（除非同时需要列或样本行）。", Args: "无", Returns: "恒为对象 {\"total\": 整数}。与 main_slice / get_status.mainCount 的 total 含义一致。"},
	{Op: "main_slice", Description: "【场景】按页遍历主列表、导出表格、或把列展示给大模型分析。【效果】返回一页行数据；断点等待中的行可正常读取（状态列显示拦截上行/下行）。【顺序】与界面一致（过滤后、theology 升序）。", Args: "offset limit", Returns: "JSON：columns、rows[{listIndex,theology,rowId,cells}]、total、offset、limit。"},
	{Op: "session_get_json", Description: "【场景】需要单条会话的完整结构化信息（方法、URL、状态、时间、元数据等），体积可能较大。【定位】用户说「第 N 条」时必须传 listIndex:N；指定会话主键用 theology 或 rowId（rowId=strconv(theology)）。【勿用】勿把 listIndex 填进 rowId/theology。【勿用】超大 body 请配合 http_get_part。", Args: "listIndex | listIndexes | theology | theologies | rowId | rowIds | ids", Returns: "单条：会话详情大 JSON（含 theology、rowId）；多条：{\"items\":[{\"theology\",\"session\"},…]}。"},
	{Op: "main_cells", Description: "【场景】已知 rowId，只要这一行在界面上的各列字符串，不要整页 JSON。", Args: "rowId | rowIds | ids", Returns: "单条：该行 cells JSON；多条：{\"items\":[…]}。"},
	{Op: "main_clear", Description: "【场景】清空主列表与关联流表展示（与 UI 清空列表同源钩子；会联动会话清理逻辑）。", Args: "无", Returns: "{\"ok\":true}。"},
	{Op: "records_import", Description: "【场景】从 SunnyNetV4 记录文件（.sy4）还原主列表会话，与 UI「打开记录文件」同源。【效果】写入内存并刷新主界面列表。", Args: "filePath（.sy4 绝对路径）", Returns: "{\"ok\":true,\"count\":n,\"path\":string}；失败抛错。"},
	{Op: "records_export", Description: "【场景】将主列表会话导出为 .sy4 记录文件，与 UI「保存记录」同源。【范围】未指定 theology/rowId/listIndex 时导出当前内存中全部会话。", Args: "filePath（.sy4 绝对路径） theology? theologies? rowId? rowIds? listIndex? listIndexes?", Returns: "{\"ok\":true,\"count\":n,\"path\":string}。"},
	{Op: "session_pack_export", Description: "【场景】导出单条 WebSocket/TCP 会话的流消息包（.bin），供代码生成器等使用，与 ExportMessage 同源。", Args: "filePath（.bin 绝对路径） theology | rowId | listIndex", Returns: "{\"ok\":true,\"path\":string,\"theology\":n,\"rowId\":string}；HTTP 会话不支持。"},
	{Op: "main_delete", Description: "【场景】按多条主键删除行（ids 为主列表 # 列 id 字符串数组）。", Args: "ids: string[]", Returns: "{\"ok\":true}。"},
	{Op: "main_delete_except", Description: "【场景】「只保留这些 id」，删除其余所有行（keepIds 为保留列表）。", Args: "keepIds: string[]", Returns: "{\"ok\":true}。"},
	{Op: "main_apply_row_mark", Description: "【场景】给指定行打颜色标记（仅当对应标记单元格仍为空时写入，避免覆盖已有标记）。", Args: "ids: string[] mark", Returns: "{\"ok\":true}。"},
	{Op: "main_search", Description: "【场景】在内存中的全部会话里搜索，要求唯一命中。【效果】只读搜索，不改 UI 筛选状态；命中唯一时返回 theology、listIndex（界面序号）、rowId。【勿用】多条命中时请缩小 needle 或改用 main_slice；不要用序号当 theology。", Args: "needle | queryText searchType? searchRange? caseInsensitive? stripSpaces?", Returns: "唯一：{\"ok\":true,\"unique\":true,\"theology\",listIndex,rowId}；多条：{\"ok\":false,\"unique\":false,\"count\",theologies}；0 条抛错。"},
	{Op: "main_row_note_get", Description: "【批量】≥2 行必须一次传 rowIds/ids，禁止逐条调用。【场景】读取主列表「注释」列（与 UI 同源）。", Args: "rowId | rowIds | ids", Returns: "单条：{\"rowId\",\"note\",\"found\"}；多条：{\"items\":[…]}。"},
	{Op: "main_row_note_set", Description: "【批量必选】≥2 行必须一次传 rowIds/ids + note，禁止对每行单独调用本工具。【场景】写入或清空注释（与 UI、MainSetNote 同源）；note 空字符串表示清空。", Args: "rowIds|ids（多行推荐） rowId（仅单行） note", Returns: "单条：{\"ok\":true,\"rowId\"}；多条：{\"ok\":true,\"rowIds\":[],\"count\":n}；任一 rowId 不存在则抛错。"},
	{Op: "main_row_break_get", Description: "【批量】读取主列表行的断点/拦截状态（与 UI「断点模式」列、请求详情放行按钮同源）。【字段】breakMode：0 未拦截；1 请求阶段断点；2 响应阶段断点。isWaiting：是否正卡在 WaitGroup。interceptState：非拦截/上行/下行/拦截。", Args: "theology | theologies | rowId | rowIds | ids | listIndex | listIndexes", Returns: "单条：{theology,rowId,found,breakMode,isWaiting,interceptState,state}；多条：{\"items\":[…]}。"},
	{Op: "stream_count", Description: "【场景】某 TCP/UDP/WebSocket 会话展开后，子表「流消息」在过滤条件下的总行数。【协定】与 main_count 相同：返回对象 {\"total\"}，不返回裸整数。", Args: "streamKey（一般为 theology 的十进制字符串，与 UI 流表一致）", Returns: "恒为对象 {\"total\": 整数}；解析后读 .total。与 stream_slice 内 total 同义。"},
	{Op: "stream_slice", Description: "【场景】分页读取某会话的子流消息表（无 main_slice 的全局 columns，列语义随流类型变化）。", Args: "streamKey offset limit", Returns: "流表分页 JSON（含 rows、total 等）。"},
	{Op: "http_get_part", Description: "【场景】HTTP 请求体/响应体或 rawRequest/rawResponse 分块读取。【part】requestBody | responseBody | rawRequest | rawResponse。【type】auto（默认，含 C0 控制字节则 base64 否则 str）| hex | base64 | str。【返回】{\"ok\", \"total\", \"type\":实际类型, \"data\"}；HTTP 无 frameType。maxLen=0 时从 offset 起最多 4MB；仅要 total 可先 offset=0 maxLen=0 读 total 字段。", Args: "theology|theologies|rowId|rowIds|ids part type? offset? maxLen?", Returns: "单条：分块 JSON；多条：{\"items\":[…]}。"},
	{Op: "stream_get_part", Description: "【场景】TCP/UDP/WebSocket 单条流消息正文。【type】同 http_get_part。【返回】WebSocket 额外含 frameType（Text/Binary，来自流表类型列）；TCP/UDP 无 frameType。", Args: "theology|theologies|rowId|rowIds|ids messageId type? offset? maxLen?", Returns: "单条：分块 JSON；多条：{\"items\":[…]}。"},
	{Op: "stream_get_hex", Description: "【场景】流消息按十六进制读取；等价 stream_get_part 且 type 默认为 hex。", Args: "theology|theologies|rowId|rowIds|ids messageId type? offset? maxLen?", Returns: "同 stream_get_part。"},
	{Op: "config_get_replace", Description: "【场景】读取「请求拦截/数据替换」页全部规则（与 ReplaceBody UI 同源）。分类：type=「拦截请求」为拦截(ruleKind=intercept)，其余为替换(ruleKind=replace)。", Args: "无", Returns: "{\"rules\":[…],\"total\":n,\"convention\":{…}}；rules 含 ruleKind、lodHint、newHint；type=替换类型，source=查找范围，lod=旧数据/匹配数据，new=新数据。"},
	{Op: "config_set_replace", Description: "【场景】全量覆盖全部规则（拦截+替换）；rulesJSON 须为完整数组。", Args: "rulesJSON", Returns: "{\"ok\":true,\"total\":n}。"},
	{Op: "config_get_rewrite", Description: "【场景】仅读取替换规则（type≠拦截请求：字符串 UTF8/GBK、Base64、十六进制）。source=查找范围：旧数据出现则替换为新数据。", Args: "无", Returns: "同 config_get_replace，但 rules 仅含 ruleKind=replace；convention.replace 列出可用 source。"},
	{Op: "config_set_rewrite", Description: "【场景】按 id 新增/更新替换规则；type 不得为「拦截请求」。默认 upsert。replaceAll:true 可整类替换。", Args: "rulesJSON replaceAll?", Returns: "{\"ok\":true,\"total\":n}。"},
	{Op: "config_get_intercept", Description: "【场景】仅读取拦截规则（type=「拦截请求」）。source=查找范围：lod 匹配则添加 HTTP 断点。", Args: "无", Returns: "同 config_get_replace，但 rules 仅含 ruleKind=intercept；convention.intercept 列出可用 source。"},
	{Op: "config_set_intercept", Description: "【场景】按 id 新增/更新拦截规则；type 须为「拦截请求」。默认 upsert。【禁用】config_rule_set_state。【删除】replaceAll:true 且 rules 不含该 id。", Args: "rulesJSON replaceAll?", Returns: "{\"ok\":true,\"total\":n}。"},
	{Op: "config_rule_set_state", Description: "【场景】启用/禁用「数据替换/拦截」单条规则（ReplaceRoles）。【勿用于 Host】Host 删除用 config_host_delete。", Args: "id state（已启用|已禁用）或 id enabled（bool）", Returns: "{\"ok\":true,\"id\",state}。"},
	{Op: "config_get_block", Description: "【场景】读取「屏蔽规则」全文。", Args: "无", Returns: "屏蔽规则 JSON。"},
	{Op: "config_set_block", Description: "【场景】全量覆盖屏蔽规则。", Args: "rulesJSON", Returns: "{\"ok\":true}。"},
	{Op: "config_get_host", Description: "【场景】读取 Host 映射规则列表。", Args: "无", Returns: "{\"rules\":[{\"id\",lod,new,note}],\"total\":n}。"},
	{Op: "config_set_host", Description: "【场景】全量覆盖 Host 规则；rulesJSON 为数组 [{id,lod,new,note},...]。", Args: "rulesJSON", Returns: "{\"ok\":true,\"total\":n}。"},
	{Op: "config_host_add", Description: "【场景】新增一条 Host 映射并解析生效。", Args: "lod new note? id?", Returns: "{\"ok\":true,\"id\",lod,new,note}。"},
	{Op: "config_host_delete", Description: "【场景】删除一条 Host 映射（按 id 或旧 Host lod）。", Args: "id | lod", Returns: "{\"ok\":true,\"id\":n}。"},
	{Op: "config_host_update", Description: "【场景】修改指定 id 的 Host 映射。", Args: "id lod? new? note?", Returns: "{\"ok\":true,\"id\",lod,new,note}。"},
	{Op: "config_get_proxy_dns", Description: "【场景】读取 DNS 解析方式（设置-上游代理/二级代理）。【模式】local/remote/remotes。", Args: "无", Returns: "{\"mode\",\"remoteServer\",\"raw\",convention}。"},
	{Op: "config_set_proxy_dns", Description: "【场景】设置 DNS 并立即 SetDnsServer 生效；remotes 须 remoteServer（:853）。", Args: "mode remoteServer?", Returns: "{\"ok\",applied,engine,dns,mode}。"},
	{Op: "config_reapply_engine", Description: "【场景】将 Config 全部引擎项重新同步到 Sunny（TCP/UDP/DNS/代理/MustTcp/JA3/请求大小等）。", Args: "无", Returns: "{\"ok\",applied,engine}。"},
	{Op: "config_get_proxy_way", Description: "【场景】读取上游代理列表（含 type：http/https/socks5）。", Args: "无", Returns: "{\"proxies\":[{\"id\",\"url\",\"type\",\"state\",\"note\",\"enabled\"}],\"total\":n,\"convention\":{...}}。"},
	{Op: "config_proxy_way_add", Description: "【场景】新增上游代理（仅 http/https/socks5）。【必填】url 或 scheme+host+port。【认证】无账号 scheme://host:port；有账号 scheme://user:pass@host:port 或 username+password。", Args: "url | scheme host port username? password? state? note? enabled?", Returns: "{\"ok\",applied,engine,proxy:{id,url,type,state,note}}。"},
	{Op: "config_proxy_way_update", Description: "【场景】修改上游代理 url/状态/注释（url 须为 http/https/socks5|socket）。", Args: "id url? scheme? host? port? username? password? state? note? enabled?", Returns: "{\"ok\",applied,engine,proxy}。"},
	{Op: "config_proxy_way_update_note", Description: "【场景】仅修改指定上游代理的注释。", Args: "id note", Returns: "{\"ok\",applied,engine,note,proxy}。"},
	{Op: "config_proxy_way_delete", Description: "【场景】删除指定上游代理。", Args: "id", Returns: "{\"ok\",applied,engine,id}。"},
	{Op: "config_proxy_way_set_state", Description: "【场景】启用/禁用指定上游代理（同时仅一条启用）。", Args: "id state? | id enabled?", Returns: "{\"ok\",applied,engine,state}。"},
	{Op: "config_get_proxy_roles", Description: "【场景】读取「不使用上游代理」规则文本（; 或换行分割）。", Args: "无", Returns: "{\"roles\":string,convention}。"},
	{Op: "config_set_proxy_roles", Description: "【场景】设置不使用上游代理规则。", Args: "roles", Returns: "{\"ok\":true}。"},
	{Op: "config_get_must_tcp", Description: "【场景】读取强制 TCP 模式与规则。", Args: "无", Returns: "{\"type\",\"typeUI\",\"roles\",convention}。"},
	{Op: "config_set_must_tcp", Description: "【场景】设置强制 TCP。", Args: "typeUI|type roles?", Returns: "{\"ok\":true,\"type\",\"typeUI\"}。"},
	{Op: "config_get_engine_toggles", Description: "【场景】一次读取禁用 TCP/UDP、浏览器缓存及限制请求大小。", Args: "无", Returns: "{\"disableTCP\",\"disableUDP\",\"disableCache\",\"limitRequestSize\",\"items\":[…]}。"},
	{Op: "config_get_disable_tcp", Description: "【场景】读取「禁用TCP」开关状态（与设置-基础设置一致）。", Args: "无", Returns: "{\"disableTCP\":bool,\"disabled\":bool,\"enabled\":bool,\"label\":\"禁用TCP\"}。"},
	{Op: "config_get_disable_udp", Description: "【场景】读取「禁用UDP」开关状态。", Args: "无", Returns: "{\"disableUDP\":bool,\"disabled\":bool,\"enabled\":bool,\"label\":\"禁用UDP\"}。"},
	{Op: "config_get_disable_cache", Description: "【场景】读取「禁用浏览器缓存」开关状态。", Args: "无", Returns: "{\"disableCache\":bool,\"disabled\":bool,\"enabled\":bool,\"label\":\"禁用浏览器缓存\"}。"},
	{Op: "config_set_disable_tcp", Description: "【场景】禁用/启用 TCP，立即 App.DisableTCP。", Args: "disableTCP|disabled", Returns: "{\"ok\",applied,engine,disableTCP}。"},
	{Op: "config_set_disable_udp", Description: "【场景】禁用/启用 UDP，立即 App.DisableUDP。", Args: "disableUDP|disabled", Returns: "{\"ok\",applied,engine,disableUDP}。"},
	{Op: "config_set_disable_cache", Description: "【场景】禁用/启用浏览器缓存（后续 HTTP 响应处理生效）。", Args: "disableCache|disabled", Returns: "{\"ok\",applied,engine,disableCache}。"},
	{Op: "config_get_limit_request_size", Description: "【场景】读取限制请求提交大小（默认 1024000）。", Args: "无", Returns: "{\"limitRequestSize\":n}。"},
	{Op: "config_set_limit_request_size", Description: "【场景】设置限制请求提交大小。", Args: "limitRequestSize|size", Returns: "{\"ok\":true,\"limitRequestSize\":n}。"},
	{Op: "config_get_https_protocol", Description: "【场景】读取 HTTPS 发送协议、JA3、HTTP2 指纹。", Args: "无", Returns: "{\"sendIsHTTP1\",\"protocol\",\"randomJa3\",\"http2Fingerprint\"}。"},
	{Op: "config_set_https_protocol", Description: "【场景】设置协议并自动应用：保存配置+JA3立即生效；H2 时写入指纹后后续请求自动 SetHTTP2Config。【参数】protocol/sendIsHTTP1；可选 http2Fingerprint、template（内置模板名）、randomJa3。", Args: "protocol? sendIsHTTP1? http2Fingerprint? template? randomJa3?", Returns: "{\"ok\",applied,sendIsHTTP1,protocol,http2Fingerprint,http2WillApply,ja3Applied}。"},
	{Op: "config_get_random_ja3", Description: "【场景】读取是否随机 JA3。", Args: "无", Returns: "{\"randomJa3\":bool}。"},
	{Op: "config_set_random_ja3", Description: "【场景】开启/关闭随机 JA3，立即 SetRandomTLS。", Args: "randomJa3|enabled", Returns: "{\"ok\",applied,engine,randomJa3}。"},
	{Op: "config_get_http2_fingerprint", Description: "【场景】读取自定义 HTTP2 指纹 JSON（编辑框内容）。", Args: "无", Returns: "{\"http2Fingerprint\":string}。"},
	{Op: "config_set_http2_fingerprint", Description: "【场景】写入 HTTP2 指纹并自动应用（默认切 H2 优先）。", Args: "http2Fingerprint|fingerprint|template", Returns: "同 config_set_https_protocol 的 applied 字段。"},
	{Op: "config_apply_http2_template", Description: "【场景】载入内置 HTTP2 模板并自动应用。", Args: "name|template protocol?", Returns: "同 config_set_https_protocol。"},
	{Op: "config_list_http2_templates", Description: "【场景】列出内置 HTTP2 指纹模板名与配置。", Args: "无", Returns: "{\"templates\":[…],\"names\":[],\"total\":n}。"},
	{Op: "config_get_http2_template", Description: "【场景】按名称获取单个 HTTP2 模板。", Args: "name", Returns: "{\"ok\":true,\"template\":{name,label,config}}。"},
	{Op: "request_cert_list", Description: "【场景】列出请求证书页全部条目（与 Cert UI 同源）。", Args: "无", Returns: "{\"certs\":[…],\"total\":n,\"convention\":{roles,certTypes}}。"},
	{Op: "request_cert_add", Description: "【场景】新增并载入一条请求证书。【必填】certPath、domain。【P12】password 必填。【可选】role（默认解析及发送）、note。", Args: "certPath domain password? role? note?", Returns: "{\"ok\":true,\"id\",\"status\",\"cert\":{…}}。"},
	{Op: "request_cert_delete", Description: "【场景】删除指定 id 的请求证书并从引擎卸载。", Args: "id", Returns: "{\"ok\":true,\"id\":n}。"},
	{Op: "request_cert_update", Description: "【场景】修改指定证书的 domain、note、role，可选更换 certPath/password 并重新载入。", Args: "id domain? note? role? certPath? password?", Returns: "{\"ok\":true,\"id\",\"status\",\"cert\":{…}}。"},
	{Op: "device_status", Description: "【场景】在 device_load 前判断当前 OS 是否支持、驱动是否已装载。", Args: "无", Returns: "{\"isWindows\":bool,\"deviceLoaded\":bool}。"},
	{Op: "device_load", Description: "【场景】切换底层捕获驱动（与 SunnyNet.OpenDrive 一致）。【mode】0=Proxifier 类；1=NFAPI；2=Tun（具体以 SunnyNet 文档为准）。", Args: "mode（整数）", Returns: "{\"ok\":bool}。"},
	{Op: "device_process_add_name", Description: "【场景】进程捕获模式：按可执行文件名加入过滤列表，只抓指定进程。", Args: "name", Returns: "{\"ok\":true}。"},
	{Op: "device_process_del_name", Description: "【场景】从进程捕获列表移除指定名称。", Args: "name", Returns: "{\"ok\":true}。"},
	{Op: "device_process_add_pid", Description: "【场景】按 PID 精确加入进程捕获目标。", Args: "pid", Returns: "{\"ok\":true}。"},
	{Op: "device_process_del_pid", Description: "【场景】按 PID 移除进程捕获目标。", Args: "pid", Returns: "{\"ok\":true}。"},
	{Op: "device_process_cancel_all", Description: "【场景】清空所有进程名/PID 捕获目标，恢复为不限制进程。", Args: "无", Returns: "{\"ok\":true}。"},
	{Op: "ui_theme_get", Description: "【场景】自动化或 MCP 想知道当前明暗主题、以及已持久化到配置里的界面偏好（字体、行颜色等镜像键）。【效果】只读 SunnyNetV5.json 中的 themeLocal。", Args: "无", Returns: "{\"themeLocal\",\"displayMode\",\"dataTheme\",\"colorHints\"}；displayMode 为 dark/light/unknown；精确色值多在 listRowColors 等键。"},
	{Op: "ui_theme_set", Description: "【场景】远程切换主窗口明暗或反转主题，并写回配置文件；主窗口订阅 theme-local-changed 后立即应用。【mode】dark / light / toggle。", Args: "mode: dark|light|toggle", Returns: "{\"ok\":true,\"themeLocal\":object,\"displayMode\":string}。"},
	{Op: "pb_to_json", Description: "【场景】有一段原始 Protobuf 二进制（如从 stream_get_part 导出），要在文本里查看字段树或交给只懂 JSON 的工具。【参数】dataB64 为整段 PB 的 Base64；skipFirstBytes 跳过前 N 字节（与 JsonView「忽略前 N 字节」一致，用于去掉自定义包头）。", Args: "dataB64 skipFirstBytes?", Returns: "{\"text\":格式化 JSON 字符串,\"inputBytesAfterSkip\":number,\"charCount\":number}。"},
	{Op: "list_supported_ops", Description: "【场景】客户端想拿到机器可读的全量能力表（含本字段 description/args/returns），而不是依赖 tools/list 的拼接长描述。【效果】返回 version + ops 数组 + capabilities 数组。", Args: "无", Returns: "{version,ops,capabilities[]}；capabilities 每项含 op、description、args、returns。"},
}

// SupportedBridgeOps MCP TCP 桥支持的 op 名称（与 Host.Invoke 一致），由 bridgeOpCatalog 派生。
var SupportedBridgeOps []string

func init() {
	SupportedBridgeOps = make([]string, len(bridgeOpCatalog))
	for i := range bridgeOpCatalog {
		SupportedBridgeOps[i] = bridgeOpCatalog[i].Op
	}
}

type opsListEnvelope struct {
	Version      int                  `json:"version"`
	Ops          []string             `json:"ops"`
	Capabilities []BridgeOpCapability `json:"capabilities"`
}

// SupportedOpsJSON 返回能力目录：ops 为名称列表，capabilities 含 description、args、returns。
func SupportedOpsJSON() string {
	env := opsListEnvelope{
		Version:      1,
		Ops:          append([]string(nil), SupportedBridgeOps...),
		Capabilities: append([]BridgeOpCapability(nil), bridgeOpCatalog...),
	}
	b, _ := json.Marshal(env)
	return string(b)
}

// BridgeMCPHTTPPrefix 为 Streamable HTTP 注册的 MCP 工具名前缀，与裸桥接 op 区分命名空间。
const BridgeMCPHTTPPrefix = "sunnynet_"

// BridgeStdioMCPPrefix 已废弃，请使用 BridgeMCPHTTPPrefix。
const BridgeStdioMCPPrefix = BridgeMCPHTTPPrefix

// BridgeMCPTool 描述单条桥接能力在 MCP tools/list 中的工具项。
type BridgeMCPTool struct {
	MCPName     string
	Op          string
	Description string
}

// BridgeStdioMCPTool 已废弃，请使用 BridgeMCPTool。
type BridgeStdioMCPTool = BridgeMCPTool

// BridgeMCPTools 供主程序 Streamable HTTP Handler 注册全部 MCP Tool。
func BridgeMCPTools() []BridgeMCPTool {
	out := make([]BridgeMCPTool, 0, len(bridgeOpCatalog))
	for _, c := range bridgeOpCatalog {
		desc := strings.TrimSpace(c.Description)
		if opUsesBatchIDs(c.Op) && !strings.Contains(desc, "【批量") {
			desc = "【批量】多行/多会话请一次传 rowIds 或 theologies，禁止循环单条调用。 " + desc
		}
		if a := strings.TrimSpace(c.Args); a != "" && a != "无" {
			desc += " 参数：" + a
		}
		if r := strings.TrimSpace(c.Returns); r != "" {
			desc += " 返回值：" + r
		}
		desc += "（需 SunnyNetTools 主程序已启用 MCP 桥；listIndex=界面序号；theology=会话主键整数；rowId=strconv(theology)。）"
		out = append(out, BridgeMCPTool{
			MCPName:     BridgeMCPHTTPPrefix + c.Op,
			Op:          c.Op,
			Description: desc,
		})
	}
	return out
}

// BridgeStdioMCPTools 已废弃，请使用 BridgeMCPTools。
func BridgeStdioMCPTools() []BridgeMCPTool { return BridgeMCPTools() }
