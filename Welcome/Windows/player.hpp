#include <windows.h>
#include <stdbool.h> // 确保支持 C 中的 bool 类型

#pragma once

#ifdef __cplusplus
extern "C" {
#endif
void StopGIFWindow() ;
bool StartGIFWindow();
bool SetImgPath(const char* gifPath, const char* iconPath);
#ifdef __cplusplus
}
#endif
