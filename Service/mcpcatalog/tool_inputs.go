package mcpcatalog

// BatchRowIDsIn 主列表 rowId 批量入参（rowId = strconv(theology)）。
type BatchRowIDsIn struct {
	RowID  string   `json:"rowId,omitempty" jsonschema:"decimal string of theology, e.g. 66 for theology 66; only when exactly one row"`
	RowIDs []string `json:"rowIds,omitempty" jsonschema:"multiple theology decimal strings in one call, e.g. [\"12\",\"66\"]"`
	IDs    []string `json:"ids,omitempty" jsonschema:"alias of rowIds"`
}

// BatchRowIDsNoteIn 批量写主列表注释（main_row_note_set）。
type BatchRowIDsNoteIn struct {
	BatchRowIDsIn
	Note string `json:"note" jsonschema:"comment for all listed rows; empty string clears note"`
}

// BatchTheologyRowIn theology、listIndex、rowId（=strconv(theology)）可混用的批量入参。
type BatchTheologyRowIn struct {
	BatchRowIDsIn
	ListIndex   int   `json:"listIndex,omitempty" jsonschema:"main-list display row number from top (1-based), same as UI 序号 column; use this when user says 第N条"`
	ListIndexes []int `json:"listIndexes,omitempty" jsonschema:"multiple display row numbers in one call"`
	Theology    int   `json:"theology,omitempty" jsonschema:"internal session id; NOT the UI 序号 column"`
	Theologies  []int `json:"theologies,omitempty" jsonschema:"multiple internal session ids in one call"`
}

// BreakSyncRequestIn 断点同步改请求（须拦截上行；不可改 Method）。
type BreakSyncRequestIn struct {
	BatchTheologyRowIn
	RequestURL  string `json:"requestURL"`
	HeadersJSON string `json:"headersJSON"`
	BodyB64     string `json:"bodyB64"`
	// method/requestMethod 不可用于改 Method
	Continue    *bool  `json:"continue,omitempty" jsonschema:"true to release after sync; default false"`
	Release     *bool  `json:"release,omitempty"`
}

// BreakSyncResponseIn 断点同步改响应（默认不放行）。
type BreakSyncResponseIn struct {
	BatchTheologyRowIn
	StatusCode  int    `json:"statusCode"`
	HeadersJSON string `json:"headersJSON"`
	BodyB64     string `json:"bodyB64"`
	Continue    *bool  `json:"continue,omitempty"`
	Release     *bool  `json:"release,omitempty"`
}

// HTTPReplayIn HTTP 重放。
type HTTPReplayIn struct {
	BatchTheologyRowIn
	InterceptMode int `json:"interceptMode,omitempty" jsonschema:"0 normal 1 request break 2 response break"`
	RepeatCount   int `json:"repeatCount,omitempty"`
}

// HTTPGetPartIn HTTP 正文分块读取。
type HTTPGetPartIn struct {
	BatchTheologyRowIn
	Part   string `json:"part" jsonschema:"requestBody responseBody rawRequest rawResponse"`
	Type   string `json:"type,omitempty" jsonschema:"auto hex base64 str"`
	Offset int    `json:"offset,omitempty"`
	MaxLen int    `json:"maxLen,omitempty"`
}

// StreamGetPartIn 流消息正文分块读取。
type StreamGetPartIn struct {
	BatchTheologyRowIn
	MessageID int    `json:"messageId"`
	Type      string `json:"type,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	MaxLen    int    `json:"maxLen,omitempty"`
}

// MainDeleteIn 按 rowId 批量删除。
type MainDeleteIn struct {
	IDs []string `json:"ids" jsonschema:"main-list # column ids to delete"`
}

// MainDeleteExceptIn 仅保留指定 id。
type MainDeleteExceptIn struct {
	KeepIDs []string `json:"keepIds" jsonschema:"row ids to keep; delete all others"`
}

// MainApplyRowMarkIn 批量行标记。
type MainApplyRowMarkIn struct {
	IDs  []string `json:"ids"`
	Mark string   `json:"mark"`
}

// MainSliceIn 主列表分页。
type MainSliceIn struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// RecordsImportIn 导入 .sy4 记录文件。
type RecordsImportIn struct {
	FilePath string `json:"filePath" jsonschema:"absolute path to .sy4 capture record file"`
}

// RecordsExportIn 导出 .sy4 记录文件。
type RecordsExportIn struct {
	BatchTheologyRowIn
	FilePath string `json:"filePath" jsonschema:"absolute path to save .sy4 file"`
}

// SessionPackExportIn 导出单会话流消息 .bin 包。
type SessionPackExportIn struct {
	BatchTheologyRowIn
	FilePath string `json:"filePath" jsonschema:"absolute path to save .bin session pack"`
}

// StreamSliceIn 流表分页。
type StreamSliceIn struct {
	StreamKey string `json:"streamKey" jsonschema:"usually theology as decimal string"`
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
}

// StreamCountIn 流表行数。
type StreamCountIn struct {
	StreamKey string `json:"streamKey"`
}

// EngineStartIn 启动引擎。
type EngineStartIn struct {
	Port            int  `json:"port,omitempty"`
	UseSystemProxy  bool `json:"useSystemProxy,omitempty"`
}

// GenerateBuiltinCodeIn 生成内置代码（单条）。
type GenerateBuiltinCodeIn struct {
	Theology int    `json:"theology,omitempty"`
	RowID    string `json:"rowId,omitempty"`
	Language string `json:"language"`
	Module   string `json:"module"`
}

// EngineApplyAdvancedIn 下发高级设置。
type EngineApplyAdvancedIn struct {
	Payload string `json:"payload"`
	Manual  bool   `json:"manual,omitempty"`
}

// ConfigSetRulesIn 全量写规则 JSON。
type ConfigSetRulesIn struct {
	RulesJSON  string `json:"rulesJSON"`
	ReplaceAll bool `json:"replaceAll,omitempty"`
}

// ConfigRuleSetStateIn 单条规则启用/禁用。
type ConfigRuleSetStateIn struct {
	ID      int    `json:"id"`
	State   string `json:"state,omitempty" jsonschema:"已启用 or 已禁用"`
	Enabled *bool  `json:"enabled,omitempty"`
}

// DeviceLoadIn 加载驱动。
type DeviceLoadIn struct {
	Mode int `json:"mode"`
}

// DeviceProcessNameIn 按进程名。
type DeviceProcessNameIn struct {
	Name string `json:"name"`
}

// DeviceProcessPIDIn 按 PID。
type DeviceProcessPIDIn struct {
	PID int `json:"pid"`
}

// UIThemeSetIn 切换主题。
type UIThemeSetIn struct {
	Mode string `json:"mode" jsonschema:"dark light toggle"`
}

// PbToJSONIn Protobuf 转 JSON。
type PbToJSONIn struct {
	DataB64         string `json:"dataB64"`
	SkipFirstBytes  int    `json:"skipFirstBytes,omitempty"`
}

// HostRuleAddIn 新增 Host 映射。
type HostRuleAddIn struct {
	Lod  string `json:"lod"`
	New  string `json:"new"`
	Note string `json:"note,omitempty"`
	ID   int    `json:"id,omitempty"`
}

// HostRuleDeleteIn 删除 Host 映射。
type HostRuleDeleteIn struct {
	ID  int    `json:"id,omitempty"`
	Lod string `json:"lod,omitempty"`
}

// HostRuleUpdateIn 更新 Host 映射。
type HostRuleUpdateIn struct {
	ID   int    `json:"id"`
	Lod  string `json:"lod,omitempty"`
	New  string `json:"new,omitempty"`
	Note string `json:"note,omitempty"`
}

// ProxyDnsSetIn DNS 解析方式。
type ProxyDnsSetIn struct {
	Mode         string `json:"mode" jsonschema:"local remote remotes"`
	RemoteServer string `json:"remoteServer,omitempty"`
}

// ProxyWayAddIn 新增上游代理（仅 http/https/socks5）。
type ProxyWayAddIn struct {
	URL      string `json:"url,omitempty"`
	Scheme   string `json:"scheme,omitempty" jsonschema:"http https socks5"`
	Type     string `json:"type,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	State    string `json:"state,omitempty"`
	Note     string `json:"note,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

// ProxyWayIDIn 上游代理 id。
type ProxyWayIDIn struct {
	ID int `json:"id"`
}

// ProxyWayUpdateNoteIn 仅改注释。
type ProxyWayUpdateNoteIn struct {
	ID   int    `json:"id"`
	Note string `json:"note"`
}

// ProxyWaySetStateIn 启用/禁用上游代理。
type ProxyWaySetStateIn struct {
	ID      int    `json:"id"`
	State   string `json:"state,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

// ProxyWayUpdateIn 更新上游代理。
type ProxyWayUpdateIn struct {
	ID       int    `json:"id"`
	URL      string `json:"url,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	State    string `json:"state,omitempty"`
	Note     string `json:"note,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

// ProxyRolesSetIn 不使用上游代理规则。
type ProxyRolesSetIn struct {
	Roles string `json:"roles"`
}

// MustTcpSetIn 强制 TCP。
type MustTcpSetIn struct {
	TypeUI string `json:"typeUI,omitempty"`
	Type   int    `json:"type,omitempty"`
	Roles  string `json:"roles,omitempty"`
}

// BoolToggleIn 布尔开关。
type BoolToggleIn struct {
	Disabled *bool `json:"disabled,omitempty"`
	Enabled  *bool `json:"enabled,omitempty"`
}

// LimitRequestSizeIn 限制请求大小。
type LimitRequestSizeIn struct {
	LimitRequestSize int `json:"limitRequestSize"`
	Size             int `json:"size,omitempty"`
}

// HTTPSProtocolSetIn HTTPS 协议（含可选指纹/JA3，保存后自动应用）。
type HTTPSProtocolSetIn struct {
	SendIsHTTP1      *bool  `json:"sendIsHTTP1,omitempty"`
	Protocol         string `json:"protocol,omitempty"`
	HTTP2Fingerprint string `json:"http2Fingerprint,omitempty"`
	Template         string `json:"template,omitempty"`
	RandomJa3        *bool  `json:"randomJa3,omitempty"`
}

// RandomJa3SetIn 随机 JA3。
type RandomJa3SetIn struct {
	RandomJa3 *bool `json:"randomJa3,omitempty"`
	Enabled   *bool `json:"enabled,omitempty"`
}

// HTTP2FingerprintSetIn 自定义 HTTP2 指纹。
type HTTP2FingerprintSetIn struct {
	HTTP2Fingerprint string `json:"http2Fingerprint,omitempty"`
	Fingerprint      string `json:"fingerprint,omitempty"`
	Template         string `json:"template,omitempty"`
}

// HTTP2TemplateApplyIn 应用内置模板。
type HTTP2TemplateApplyIn struct {
	Name     string `json:"name,omitempty"`
	Template string `json:"template,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// HTTP2TemplateGetIn 模板名。
type HTTP2TemplateGetIn struct {
	Name string `json:"name"`
}

// RequestCertAddIn 新增请求证书。
type RequestCertAddIn struct {
	CertPath string `json:"certPath"`
	Domain   string `json:"domain"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty" jsonschema:"解析及发送 仅解析 仅发送"`
	Note     string `json:"note,omitempty"`
}

// RequestCertIDIn 按 id 操作请求证书。
type RequestCertIDIn struct {
	ID int `json:"id"`
}

// RequestCertUpdateIn 更新请求证书。
type RequestCertUpdateIn struct {
	ID       int    `json:"id"`
	Domain   string `json:"domain,omitempty"`
	Note     string `json:"note,omitempty"`
	Role     string `json:"role,omitempty" jsonschema:"解析及发送 仅解析 仅发送"`
	CertPath string `json:"certPath,omitempty"`
	Password string `json:"password,omitempty"`
}

// MainSearchIn 内存搜索（唯一命中时返回 listIndex）。
type MainSearchIn struct {
	QueryText       string `json:"queryText,omitempty"`
	Needle          string `json:"needle,omitempty"`
	SearchType      string `json:"searchType,omitempty"`
	CaseInsensitive bool   `json:"caseInsensitive,omitempty"`
	StripSpaces     bool   `json:"stripSpaces,omitempty"`
	SearchRange     string `json:"searchRange,omitempty"`
}
