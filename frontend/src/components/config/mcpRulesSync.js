import {Events} from "@wailsio/runtime";
import {McpFuncRes} from "../../../bindings/changeme/Service/appmain.js";

/**
 * 规则页监听 MCP configreload（page=main, tag=configreload）。
 * @param {string} ruleKey replace | host | intercept | rewrite | block
 * @param {(action: string, body: object) => void | Promise<void>} reloadFn 重新从后端拉列表并刷新表格
 */
export function attachMcpConfigReload(ruleKey, reloadFn) {
  Events.On("mcp", async (evt) => {
    const mcp = evt?.data ?? {};
    const reply = (text) => {
      mcp.res = text;
      typeof McpFuncRes === "function" && McpFuncRes(mcp);
    };
    try {
      const page = String(mcp.page ?? "").toLowerCase();
      const tag = String(mcp.tag ?? "").toLowerCase();
      if (page !== "main" || tag !== "configreload") {
        return;
      }
      let body = {};
      try {
        body = JSON.parse(String(mcp.msg ?? "{}"));
      } catch (_) {
        return reply("invalid payload");
      }
      const rule = String(body.rule ?? "").toLowerCase();
      const mine = String(ruleKey ?? "").toLowerCase();
      if (rule !== "" && rule !== mine && rule !== "all") {
        return;
      }
      const action = String(body.action ?? "reload").toLowerCase();
      await reloadFn(action, body);
      return reply("ok");
    } catch (e) {
      try {
        reply("处理失败");
      } catch (_) {
      }
    }
  });
}
