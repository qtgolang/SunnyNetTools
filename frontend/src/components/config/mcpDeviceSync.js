import {Events} from "@wailsio/runtime";
import {McpFuncRes} from "../../../bindings/changeme/Service/appmain.js";
import {ObjString} from "./Config.js";

function parseDeviceMcp(evt) {
  const mcp = evt?.data ?? {};
  const page = String(mcp.page ?? "").toLowerCase();
  const tag = String(mcp.tag ?? "").toLowerCase();
  if (page !== "main" || tag !== "device") {
    return null;
  }
  let body = {};
  try {
    body = JSON.parse(String(mcp.msg ?? "{}"));
  } catch (_) {
    body = {};
  }
  return {mcp, body};
}

/**
 * 进程名列表（Device/Name.vue）MCP 同步：增 / 删 / 清空 / 全量同步。
 */
export function attachMcpDeviceName(vm) {
  Events.On("mcp", async (evt) => {
    const parsed = parseDeviceMcp(evt);
    if (!parsed) {
      return;
    }
    const {mcp, body} = parsed;
    const reply = (text) => {
      mcp.res = text;
      typeof McpFuncRes === "function" && McpFuncRes(mcp);
    };
    const api = vm.agGridApi;
    if (!api) {
      return reply("grid not ready");
    }
    const action = String(body.action ?? "").toLowerCase();
    try {
      switch (action) {
        case "add_name": {
          const name = ObjString(body.name ?? "").trim();
          if (!name) {
            return reply("empty name");
          }
          let exists = false;
          api.forEachNode((node) => {
            if (ObjString(node.data["进程名称"]).trim().toLowerCase() === name.toLowerCase()) {
              exists = true;
            }
          });
          if (!exists) {
            vm.addLine(name);
          }
          return reply("ok");
        }
        case "del_name": {
          const name = ObjString(body.name ?? "").trim().toLowerCase();
          const remove = [];
          api.forEachNode((node) => {
            if (ObjString(node.data["进程名称"]).trim().toLowerCase() === name) {
              remove.push(node.data);
            }
          });
          if (remove.length > 0) {
            api.applyTransaction({remove});
          }
          return reply("ok");
        }
        case "clear_names":
          api.setGridOption("rowData", []);
          vm.RowNodes = [];
          return reply("ok");
        case "sync_names": {
          const names = Array.isArray(body.names) ? body.names : [];
          api.setGridOption("rowData", []);
          vm.RowNodes = [];
          for (const n of names) {
            const name = ObjString(n).trim();
            if (name) {
              vm.addLine(name);
            }
          }
          return reply("ok");
        }
        default:
          return;
      }
    } catch (e) {
      return reply("处理失败");
    }
  });
}

/**
 * 进程 PID 列表（Device/Pid.vue）MCP 同步选中状态。
 */
export function attachMcpDevicePid(vm) {
  Events.On("mcp", async (evt) => {
    const parsed = parseDeviceMcp(evt);
    if (!parsed) {
      return;
    }
    const {mcp, body} = parsed;
    const reply = (text) => {
      mcp.res = text;
      typeof McpFuncRes === "function" && McpFuncRes(mcp);
    };
    const api = vm.agGridApi;
    if (!api) {
      return reply("grid not ready");
    }
    const action = String(body.action ?? "").toLowerCase();
    try {
      switch (action) {
        case "add_pid": {
          const pid = parseInt(body.pid, 10);
          if (!Number.isFinite(pid)) {
            return reply("invalid pid");
          }
          const id = String(pid);
          api.forEachNode((node) => {
            if (String(node.data.id) === id || parseInt(node.data.PID, 10) === pid) {
              node.setSelected(true);
              vm.previousSelectedIds.add(id);
            }
          });
          return reply("ok");
        }
        case "del_pid": {
          const pid = parseInt(body.pid, 10);
          if (!Number.isFinite(pid)) {
            return reply("invalid pid");
          }
          const id = String(pid);
          api.forEachNode((node) => {
            if (String(node.data.id) === id || parseInt(node.data.PID, 10) === pid) {
              node.setSelected(false);
              vm.previousSelectedIds.delete(id);
            }
          });
          return reply("ok");
        }
        case "clear_pids":
          api.deselectAll();
          vm.previousSelectedIds = new Set();
          return reply("ok");
        default:
          return;
      }
    } catch (e) {
      return reply("处理失败");
    }
  });
}

/**
 * 设备页根组件：驱动加载、ProcessAny 等状态同步。
 */
export function attachMcpDeviceRoot(vm) {
  Events.On("mcp", async (evt) => {
    const parsed = parseDeviceMcp(evt);
    if (!parsed) {
      return;
    }
    const {mcp, body} = parsed;
    const reply = (text) => {
      mcp.res = text;
      typeof McpFuncRes === "function" && McpFuncRes(mcp);
    };
    const action = String(body.action ?? "").toLowerCase();
    try {
      switch (action) {
        case "device_loaded":
          if (body.loaded != null) {
            vm.LoadDrive = !!body.loaded;
            vm.DriveLoading = false;
          }
          if (body.mode != null) {
            const m = String(body.mode);
            if (m === "NFAPI" || m === "Proxifier" || m === "Tun") {
              vm.LoadMode = m;
            }
          }
          return reply("ok");
        case "process_any":
          if (body.open != null) {
            vm.ProcessesAny = !!body.open;
          }
          return reply("ok");
        default:
          return;
      }
    } catch (e) {
      return reply("处理失败");
    }
  });
}
