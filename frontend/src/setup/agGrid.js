import {AllEnterpriseModule, IntegratedChartsModule, LicenseManager} from "ag-grid-enterprise";
import {AgChartsEnterpriseModule} from "ag-charts-enterprise";
import {
    AllCommunityModule,
    ModuleRegistry,
    themeAlpine,
    themeBalham,
    themeMaterial,
    themeQuartz,
} from "ag-grid-community";
import {LicenseKey} from "../AGLicenseKey";

let initialized = false;

export function setupAgGrid() {
    if (initialized) {
        return;
    }
    initialized = true;
    window.themes = [
        {id: "themeQuartz", theme: themeQuartz},
        {id: "themeBalham", theme: themeBalham},
        {id: "themeMaterial", theme: themeMaterial},
        {id: "themeAlpine", theme: themeAlpine},
    ];
    ModuleRegistry.registerModules([
        AllCommunityModule,
        IntegratedChartsModule.with(AgChartsEnterpriseModule),
        AllEnterpriseModule,
    ]);
    LicenseManager.setLicenseKey(LicenseKey);
}
