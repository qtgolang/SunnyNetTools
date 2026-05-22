# SunnyNetTools





### 基于 SunnyNet 的跨平台网络抓包与调试桌面工具

类似 Fiddler / Charles 的图形化中间人分析工具，内置规则引擎、断点改包、脚本扩展与 MCP 自动化桥接



---

## 项目简介

**SunnyNetTools** 是 [SunnyNet](https://github.com/qtgolang/SunnyNet) 网络中间件 SDK 的**官方图形化配套工具**。它将 SunnyNet 的抓包、改包、重放、驱动捕获等能力封装为开箱即用的桌面应用，面向安全研究、接口调试、自动化测试与 AI 辅助分析等场景。

与仅提供 SDK 的 SunnyNet 不同，本仓库侧重：

- 可视化会话列表与详情编辑（AG Grid + Monaco 编辑器）
- 完整的设置中心（规则、证书、代理、驱动、主题等）
- **MCP（Model Context Protocol）HTTP 桥**，供 Cursor、Claude 等客户端以 `op` 方式驱动主程序

> SunnyNet 负责底层网络栈；SunnyNetTools 负责交互、配置持久化与 MCP 对外接口。

---

## 开源协议与二次开发

本仓库及所依赖的 **SunnyNet** 网络核心均采用 **[MIT License](LICENSE)**（与 [SunnyNet 仓库](https://github.com/qtgolang/SunnyNet) 一致）。

**凡基于 SunnyNet / SunnyNetTools 进行二次开发、集成、分发或商用，请务必：**

1. **保留** MIT 协议全文及版权声明（见仓库根目录 [LICENSE](LICENSE)）。
2. **在衍生作品**（修改版、插件、内嵌 SDK、打包发行物等）的源码或发行包中**同样附带**上述许可声明，不得移除或篡改。
3. 遵守 MIT 条款中的免责声明；违法用途与作者无关（应用内亦提示：禁止一切违法用途）。

未经遵守开源协议的分发、闭源篡改后商用等行为，可能构成对版权的侵犯。二次开发前请完整阅读 [LICENSE](LICENSE)。

---

## 主要特性

### 抓包与会话分析

- 多协议：**HTTP / HTTPS / WebSocket / TCP / UDP**
- 主列表：过滤、排序、注释、颜色标记、内存搜索、批量删除
- 会话详情：请求/响应头与正文（文本 / Hex / JSON / Protobuf）
- 流消息表：WS / TCP / UDP 子消息分页查看、Hex 展示、主动发包
- 压缩解码：gzip、deflate、br、zstd、zlib 等
- **断点模式**：请求阶段 / 响应阶段拦截，在线改 URL、头、体后放行或继续
- **重放**：单条或批量 HTTP 重放，可带断点模式
- **记录文件**：导入 / 导出 `.sy4`（SunnyNetV4 记录格式）；WS/TCP 可导出 `.bin` 会话包
- 支持拖拽 `.sy4` 到主窗口还原记录

### 规则与流量控制

- **请求拦截 / 数据替换**：按 URL、协议头、提交数据、响应等范围匹配；支持 UTF-8、GBK、Base64、十六进制
- **Host 映射**：域名重定向
- **上游代理**：HTTP / HTTPS / SOCKS5 链式代理，账号密码、启用状态
- **DNS**：本地 / 远程 / DoH 等模式
- **强制 TCP（MustTcp）**：按规则强制走 TCP
- **请求证书**：多证书管理，解析/发送角色，P12 支持
- **Go 脚本**：SunnyNet 脚本回调，可自定义处理逻辑（脚本日志窗口）

### 引擎与系统

- 监听端口、启动/停止 Sunny 引擎
- 一键设置 / 取消**系统代理**（Windows）
- 高级选项：禁用 TCP/UDP、限制 POST 体大小、禁用浏览器缓存、出口路由、Socks5 认证等
- **TLS / HTTP2**：JA3、HTTP/2 指纹模板、协议优先级
- **Windows 驱动捕获**：Proxifier / NFAPI / Tun(WinDivert) 模式
- **进程过滤**：按进程名或 PID 指定抓包范围

### 内置工具箱


| 工具          | 说明                               |
| ----------- | -------------------------------- |
| HTTP 调试     | 类 Postman 发请求，可选中抓包会话填充          |
| 编码转换        | 常见编码互转                           |
| 加密解密        | AES/DES/3DES/RC4、RSA、SM2、带盐算法等   |
| 导出证书        | 导出 SunnyNet 根证书 `.cer`           |
| 证书安装        | Windows / macOS / Android 等安装指引  |
| 代码生成        | Go / C# / Python / cURL / 易语言等模板 |
| 文本对比        | 双栏 diff                          |
| 主题调色 / 主题设计 | 明暗主题与列表行颜色                       |
| 自定义工具       | 挂载外部可执行程序（Windows）               |


### MCP 自动化桥（可选）

默认**不自动启动**，可在底部状态栏 **MCP** 入口手动启用（默认端口 `6987`）。

启用后提供：


| 端点                                         | 用途                          |
| ------------------------------------------ | --------------------------- |
| `GET /doc`                                 | 内置 MCP 文档页                  |
| `GET /sunnynet/tools/health`               | 健康检查                        |
| `GET /sunnynet/tools/supported-ops`        | 能力目录 JSON                   |
| `POST /sunnynet/tools/invoke`              | REST 调用业务 `op`              |
| `GET /sunnynet/tools/events`               | 列表变更 SSE                    |
| `http://127.0.0.1:6987/sunnynet/tools/mcp` | Cursor 等 **Streamable MCP** |


**Cursor `mcp.json` 示例：**

```json
{
  "mcpServers": {
    "SunnyNetTools": {
      "url": "http://127.0.0.1:6987/sunnynet/tools/mcp"
    }
  }
}
```

典型 `op` 能力包括：`get_status`、`engine_start` / `engine_stop`、`main_slice`、`main_search`、`break_sync_request` / `break_sync_response`、`http_get_part`、`stream_send`、`config_*` 规则与证书、`records_import` / `records_export`、`session_pack_export` 等。完整列表见启用后的文档页或 `list_supported_ops`。

---

## 技术栈


| 层级   | 技术                                                                            |
| ---- | ----------------------------------------------------------------------------- |
| 桌面壳  | [Wails v3](https://wails.io/)                                                 |
| 前端   | Vue 3、Element Plus、AG Grid、Monaco Editor                                      |
| 后端   | Go 1.25+                                                                      |
| 网络核心 | [github.com/qtgolang/SunnyNet](https://github.com/qtgolang/SunnyNet) v1.4.x   |
| MCP  | [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk) |


---

## 系统要求

- **Windows** 10/11（推荐；驱动捕获、系统代理完整支持）
- **Linux** / **macOS**：可编译运行，部分驱动相关功能仅 Windows 可用
- 管理员权限：使用 NFAPI / Tun 等驱动模式时通常需要

---

## 快速开始

### 环境准备

1. 安装 [Go](https://go.dev/dl/)（≥ 1.25，与 `go.mod` 一致）
2. 安装 [Node.js](https://nodejs.org/)（建议 LTS）
3. 安装 [Wails CLI v3](https://wails.io/docs/gettingstarted/installation)

### 开发模式

```bash
git clone <本仓库地址>
cd SunnyNetTools
task dev
# 或: wails3 dev -config ./build/config.yml
```

### 生产构建

```bash
task build
# 产物位于 bin/ 目录（因平台任务而异）
```

各平台详细打包见 `build/windows`、`build/linux`、`build/darwin` 下的 Taskfile。

---

## 目录结构（简要）

```
SunnyNetTools/
├── main.go                 # 程序入口
├── Service/                # Go 业务：抓包回调、配置、MCP 桥、工具
│   ├── mcp/                # MCP HTTP 服务与文档
│   ├── mcpbridge/        # invoke / SSE / Streamable MCP
│   ├── mcpcatalog/       # op 能力与 JSON Schema
│   └── Session/          # 会话存储、导入导出、代码生成
├── frontend/               # Vue 3 前端
└── build/                  # Wails 构建配置与平台任务
```

---

## 与 SunnyNet SDK 的关系


| 项目                                               | 定位                              |
| ------------------------------------------------ | ------------------------------- |
| [SunnyNet](https://github.com/qtgolang/SunnyNet) | 跨平台网络中间件 **SDK**，供二次开发          |
| **SunnyNetTools**（本仓库）                           | 基于 SDK 的 **桌面抓包工具** + **MCP 桥** |


SDK 的 API 说明、多驱动对比、各语言示例请参阅 SunnyNet 仓库：

- [Go 使用示例](https://github.com/qtgolang/SunnyNet/blob/main/README_go.md)
- [API 参考](https://github.com/qtgolang/SunnyNet/blob/main/README_api.md)

---

## 快捷键（部分）

可在「快捷键设置」中自定义，内置示例：


| 快捷键            | 功能          |
| -------------- | ----------- |
| F12            | 设置 / 取消系统代理 |
| Ctrl + Z       | 放行当前断点请求    |
| Shift + Z      | 全部放行        |
| Ctrl + R       | 批量重发        |
| Ctrl + Alt + X | 清空全部记录      |
| Alt + Q        | 老板键（隐藏窗口）   |


---

## 注意事项与免责声明

1. **合法用途**：本软件仅供学习、调试与授权测试使用。禁止用于任何未授权的网络窃听、攻击或违法活动；使用者须自行承担法律责任。
2. **开源协议**：二次开发、再分发须遵守 [LICENSE](LICENSE)（MIT）；SunnyNet SDK 许可要求同样适用，详见 [SunnyNet LICENSE](https://github.com/qtgolang/SunnyNet/blob/main/LICENSE)。
3. **HTTPS 抓包**：需安装并信任 SunnyNet 根证书，详见工具箱「证书安装」。
4. **Win7**：若需支持 Windows 7，请使用 Go 1.21 以下版本编译 SunnyNet 相关依赖（参见 SunnyNet 文档）。
5. **MCP 安全**：MCP 服务默认监听 `127.0.0.1`，请勿在不受信任的网络环境中暴露端口。

---

## 反馈与交流

- 项目网站：[https://esunny.vip/](https://esunny.vip/)
- SunnyNet 相关 QQ 群等信息见 [SunnyNet README](https://github.com/qtgolang/SunnyNet/blob/main/README.md)

---



**如果本项目对你有帮助，欢迎 Star 与 Issue**

