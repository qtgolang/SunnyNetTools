#include <stdio.h>
#include <windows.h>
#include <gdiplus.h>
#include <pthread.h>
#include <time.h>

#define global_var static

// 全局变量
global_var UINT currentFrameIndex = 0;  // 当前 GIF 动画帧索引
global_var Gdiplus::GdiplusStartupInput gdiPlusStartupInput;  // GDI+ 启动输入参数
global_var ULONG_PTR gdiPlusToken;  // GDI+ 启动令牌
global_var HWND mainWindowHandle;  // 主窗口句柄
global_var HDC mainDeviceContext;  // 主设备上下文句柄
global_var bool isAppRunning = false;  // 程序运行标志
global_var struct timespec sleepRemaining, sleepRequest;  // 时间控制结构体
global_var bool hasCenteredWindow = false;  // 是否已居中窗口
global_var pthread_t drawThreadId;  // 绘制线程 ID
wchar_t gifImagePath[MAX_PATH];  // GIF 文件路径（宽字符）
wchar_t iconWidePath[MAX_PATH];  // 宽字符图标路径
// 绘制图像函数
long DrawGIFImage(HWND hwnd, HDC hdc, const WCHAR* gifPath) {
    Gdiplus::Image* gifImage = new Gdiplus::Image(gifPath);  // 加载 GIF 图像
    int topOffset = 50;  // 图片偏移 因为我提供的gif 底部有logo 所以高度减少50 用于隐藏logo
    int imageHeight = gifImage->GetHeight() - topOffset;  // 获取图像高度
    int imageWidth = gifImage->GetWidth();  // 获取图像宽度

    // 首帧绘制时居中窗口
    if (currentFrameIndex == 0 && !hasCenteredWindow) {
        hasCenteredWindow = true;  // 设置已居中标志
        int screenWidth = GetSystemMetrics(SM_CXSCREEN);  // 获取屏幕宽度
        int screenHeight = GetSystemMetrics(SM_CYSCREEN);  // 获取屏幕高度

        int windowWidth = imageWidth;  // 窗口宽度等于图片宽度
        int windowHeight = imageHeight - topOffset;  // 窗口高度为图片高度减去偏移

        int posX = (screenWidth - windowWidth) / 2;  // 计算居中位置 X
        int posY = (screenHeight - windowHeight) / 2;  // 计算居中位置 Y

        // 设置窗口位置和大小
        SetWindowPos(hwnd, HWND_TOPMOST, posX, posY, windowWidth, windowHeight, SWP_SHOWWINDOW);
    }

    // 创建临时 Bitmap 和绘图对象
    Gdiplus::Bitmap* frameBitmap = new Gdiplus::Bitmap(imageWidth, imageHeight, PixelFormat32bppARGB);
    Gdiplus::Graphics* frameGraphics = new Gdiplus::Graphics(frameBitmap);

    UINT totalFrames = gifImage->GetFrameCount(&Gdiplus::FrameDimensionTime);  // 获取总帧数
    gifImage->SelectActiveFrame(&Gdiplus::FrameDimensionTime, currentFrameIndex);  // 选择当前帧

    int propertySize = gifImage->GetPropertyItemSize(PropertyTagFrameDelay);  // 获取帧延迟属性大小
    Gdiplus::PropertyItem* frameDelayItem = (Gdiplus::PropertyItem*)malloc(propertySize);  // 分配属性内存
    if (!frameDelayItem) {
        printf("Error: Memory allocation failed for frameDelayItem!\n");  // 内存分配失败提示
        isAppRunning = false;  // 停止运行
        Gdiplus::GdiplusShutdown(gdiPlusToken);  // 关闭 GDI+
        exit(0);  // 退出程序
    }

    gifImage->GetPropertyItem(PropertyTagFrameDelay, propertySize, frameDelayItem);  // 获取帧延迟属性
    long frameDelayMs = ((long*)frameDelayItem->value)[currentFrameIndex] * 10;  // 计算帧延迟时间（毫秒）

    currentFrameIndex = (currentFrameIndex + 1) % totalFrames;  // 更新帧索引（循环播放）

    frameGraphics->Clear(Gdiplus::Color(0, 0, 0, 0));  // 清除背景
    frameGraphics->DrawImage(gifImage, 0, 0);  // 绘制当前帧图像

    HBITMAP hFrameBitmap;  // 位图句柄
    frameBitmap->GetHBITMAP(Gdiplus::Color(0, 0, 0, 0), &hFrameBitmap);  // 获取 HBITMAP

    HDC memDC = CreateCompatibleDC(hdc);  // 创建兼容 DC
    HGDIOBJ originalObj = SelectObject(memDC, hFrameBitmap);  // 选择位图到 DC 中

    BLENDFUNCTION blendFunc = { 0 };  // 混合函数结构体
    blendFunc.BlendOp = AC_SRC_OVER;
    blendFunc.SourceConstantAlpha = 255;  // 不透明
    blendFunc.AlphaFormat = AC_SRC_ALPHA;  // 使用 alpha 通道

    POINT srcPoint = { 0, 0 };  // 图像源位置
    SIZE bmpSize = { imageWidth, imageHeight };  // 图像大小
    UpdateLayeredWindow(hwnd, hdc, 0, &bmpSize, memDC, &srcPoint, 0, &blendFunc, ULW_ALPHA);  // 更新分层窗口

    // 清理资源
    SelectObject(hdc, originalObj);  // 恢复原对象
    DeleteObject(hFrameBitmap);  // 删除位图对象
    DeleteObject(memDC);  // 删除内存 DC

    delete gifImage;  // 释放图像对象
    delete frameBitmap;  // 释放 Bitmap
    delete frameGraphics;  // 释放 Graphics
    free(frameDelayItem);  // 释放属性内存

    return frameDelayMs;  // 返回帧延迟时间
}

// 绘图线程函数
void* DrawGIFThread(void* arg) {
    long delayMs, elapsedMs;
    clock_t startTime, endTime, timeDiff;

    while (isAppRunning) {
        startTime = clock();  // 获取开始时间
        delayMs = DrawGIFImage(mainWindowHandle, mainDeviceContext, gifImagePath);  // 绘制一帧
        timeDiff = clock() - startTime;  // 计算耗时
        elapsedMs = timeDiff * 1000 / CLOCKS_PER_SEC;  // 转换为毫秒
        delayMs = (delayMs - elapsedMs >= 0) ? delayMs - elapsedMs : 0;  // 计算剩余延迟
        endTime = clock() + (delayMs * CLOCKS_PER_SEC) / 1000;  // 计算目标时间
        while (clock() < endTime);  // 等待剩余时间
    }

    pthread_exit(NULL);  // 退出线程
    return NULL;  // 返回空指针
}

// 窗口过程函数
LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
    LRESULT result = 0;
    switch (uMsg) {
        case WM_DESTROY:
            PostQuitMessage(0);  // 发出退出消息
            break;
        case WM_CLOSE:
            DestroyWindow(hwnd);  // 销毁窗口
            break;
        case WM_NCHITTEST:
            return HTCAPTION;  // 支持窗口拖动
        default:
            result = DefWindowProc(hwnd, uMsg, wParam, lParam);  // 默认消息处理
            break;
    }
    return result;  // 返回结果
}

// 初始化 GIF 播放窗口
int InitGIFWindow(HINSTANCE hInstance ) {

    WNDCLASS wc = {};  // 注册窗口类
    wc.lpfnWndProc = WindowProc;
    wc.hInstance = hInstance;
    wc.lpszClassName = "SunnyNet";

    Gdiplus::GdiplusStartup(&gdiPlusToken, &gdiPlusStartupInput, NULL);  // 启动 GDI+

    if (RegisterClass(&wc)) {  // 注册成功
        int screenWidth = GetSystemMetrics(SM_CXSCREEN);  // 获取屏幕宽度
        int screenHeight = GetSystemMetrics(SM_CYSCREEN);  // 获取屏幕高度

        int windowWidth = 800;
        int windowHeight = 600;

        int posX = (screenWidth - windowWidth) / 2;  // 窗口 X 位置
        int posY = (screenHeight - windowHeight) / 2;  // 窗口 Y 位置

        mainWindowHandle = CreateWindowEx(
            WS_EX_LAYERED | WS_EX_TOPMOST,
            wc.lpszClassName,
            "SunnyNet Welcome",
            WS_POPUP,
            posX, posY,
            windowWidth, windowHeight,
            0, 0, hInstance, 0
        );  // 创建分层窗口

        mainDeviceContext = GetDC(mainWindowHandle);  // 获取窗口 DC
        HICON windowIcon = (HICON)LoadImageW(NULL, iconWidePath, IMAGE_ICON, 0, 0, LR_LOADFROMFILE | LR_DEFAULTSIZE);  // 加载图标
        SendMessage(mainWindowHandle, WM_SETICON, ICON_BIG, (LPARAM)windowIcon);  // 设置大图标
        SendMessage(mainWindowHandle, WM_SETICON, ICON_SMALL, (LPARAM)windowIcon);  // 设置小图标

        pthread_create(&drawThreadId, NULL, DrawGIFThread, NULL);  // 启动绘图线程

        ShowWindow(mainWindowHandle, SW_SHOW);  // 显示窗口

        MSG msg = {};
        while (GetMessage(&msg, NULL, 0, 0)) {  // 消息循环
            TranslateMessage(&msg);
            DispatchMessage(&msg);
        }
    } else {
        return -1;  // 注册失败
    }

    return 0;  // 成功
}

// 停止 GIF 播放窗口
extern "C" void StopGIFWindow() {
    if (isAppRunning == false) {
        return;  // 已停止则不执行
    }

    isAppRunning = false;  // 设置停止标志
    pthread_join(drawThreadId, NULL);  // 等待绘图线程结束

    if (mainWindowHandle != NULL) {
        SendMessage(mainWindowHandle, WM_CLOSE, 0, 0);  // 关闭窗口
    }

    Gdiplus::GdiplusShutdown(gdiPlusToken);  // 关闭 GDI+

    if (mainDeviceContext != NULL) {
        ReleaseDC(mainWindowHandle, mainDeviceContext);  // 释放 DC
        mainDeviceContext = NULL;
    }
}

// 启动 GIF 播放窗口（导出接口）
extern "C" bool StartGIFWindow() {
    isAppRunning = true;  // 设置运行标志
    return InitGIFWindow(GetModuleHandle(0)) == 0;  // 初始化窗口
}
// 启动 GIF 播放窗口（导出接口）
extern "C" bool SetImgPath(const char* gifPath, const char* iconPath) {
    int len = MultiByteToWideChar(CP_UTF8, 0, gifPath, -1, gifImagePath, sizeof(gifImagePath) / sizeof(gifImagePath[0]));  // 转换 GIF 路径为宽字符
    if (len == 0) return false;
    len = MultiByteToWideChar(CP_UTF8, 0, iconPath, -1, iconWidePath, sizeof(iconWidePath) / sizeof(iconWidePath[0]));  // 转换图标路径
    if (len == 0) return false;
    return true;
}
