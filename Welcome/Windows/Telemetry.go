package Windows

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

// 确保注册表路径存在，如果不存在则创建该路径
func ensureRegistryPath(root registry.Key, path string) error {
	// 尝试打开注册表键
	k, openedExisting, err := registry.CreateKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("Error creating or opening registry key: %v", err)
	}
	defer k.Close()

	// 如果创建的是新的键，输出一些信息
	if !openedExisting {
		fmt.Printf("Created registry path: %s\n", path)
	}

	return nil
}

// 设置注册表值
func setRegistryValue(root registry.Key, path string, key string, value uint32) {
	// 确保路径存在
	err := ensureRegistryPath(root, path)
	if err != nil {
		fmt.Printf("Error ensuring registry path: %v\n", err)
		return
	}

	// 打开注册表键，设置值
	k, err := registry.OpenKey(root, path, registry.SET_VALUE)
	if err != nil {
		fmt.Printf("Error opening registry key: %v\n", err)
		return
	}
	defer k.Close()

	// 设置 DWORD 类型的值
	err = k.SetDWordValue(key, value)
	if err != nil {
		fmt.Printf("Error setting registry value: %v\n", err)
	}
}

func DisableEdgeTelemetry() {
	// 禁用个性化报告
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Edge`, "PersonalizationReportingEnabled", uint32(0))
	setRegistryValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Edge`, "PersonalizationReportingEnabled", uint32(0))
	setRegistryValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Edge`, "UserFeedbackAllowed", uint32(0))
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Edge`, "UserFeedbackAllowed", uint32(0))

	// 禁用度量报告
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Edge`, "MetricsReportingEnabled", uint32(0))
	setRegistryValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Edge`, "MetricsReportingEnabled", uint32(0))

	// 禁用扩展书籍遥测
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\MicrosoftEdge\BooksLibrary`, "EnableExtendedBooksTelemetry", uint32(0))
	setRegistryValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\MicrosoftEdge\BooksLibrary`, "EnableExtendedBooksTelemetry", uint32(0))

	// 禁用智能屏幕功能
	setRegistryValue(registry.CURRENT_USER, `Software\Microsoft\Edge`, "SmartScreenEnabled", uint32(0))
	setRegistryValue(registry.CURRENT_USER, `Software\Microsoft\Edge`, "SmartScreenPuaEnabled", uint32(0))

	// 禁用扩展清单 v2 可用性
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Edge`, "ExtensionManifestV2Availability", uint32(2))

	// 禁用第三方 serp 遥测
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Edge`, "Edge3PSerpTelemetryEnabled", uint32(0))

	// 禁用 Spotlight 推荐
	setRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Policies\Microsoft\Edge`, "SpotlightExperiencesAndRecommendationsEnabled", uint32(0))
	setRegistryValue(registry.CURRENT_USER, `SOFTWARE\Policies\Microsoft\Edge`, "SpotlightExperiencesAndRecommendationsEnabled", uint32(0))

	fmt.Println("Edge telemetry has been disabled.")
}
