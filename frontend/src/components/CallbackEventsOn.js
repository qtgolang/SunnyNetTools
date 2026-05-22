import {CallTools} from "../../bindings/changeme/Service/appmain.js";

export async function OpenTools(name, open,args) {
    if (open === undefined || open === null) {
        CallTools(name, true,args);
        return;
    }
    CallTools(name,  open,args);
}
