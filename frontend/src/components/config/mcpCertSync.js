import {Events} from "@wailsio/runtime";
import {McpFuncRes} from "../../../bindings/changeme/Service/appmain.js";

/**
 * 请求证书页（Cert 独立窗口）监听 MCP requestcert 事件。
 * @param {() => void | Promise<void>} reloadFn 从 RequestList 全量刷新表格
 */
export function attachMcpRequestCertReload(reloadFn) {
  Events.On("mcp", async (evt) => {
    const mcp = evt?.data ?? {};
    const reply = (text) => {
      mcp.res = text;
      typeof McpFuncRes === "function" && McpFuncRes(mcp);
    };
    try {
      const page = String(mcp.page ?? "").toLowerCase();
      const tag = String(mcp.tag ?? "").toLowerCase();
      if (page !== "cert" || tag !== "requestcert") {
        return;
      }
      let body = {};
      try {
        body = JSON.parse(String(mcp.msg ?? "{}"));
      } catch (_) {
        return reply("invalid payload");
      }
      if (String(body.action ?? "reload").toLowerCase() !== "reload") {
        return reply("ok");
      }
      await reloadFn();
      return reply("ok");
    } catch (e) {
      try {
        reply("处理失败");
      } catch (_) {
      }
    }
  });
}
