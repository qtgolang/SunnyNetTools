import {ref} from "vue";
import {GetTour, GOOS} from "../../../bindings/changeme/Service/appmain";

export const IsOpenTour = ref(false);
export const TourList = ref([]);
const array = []

export function Tour_Add(target, index, title, description) {
    if (target instanceof HTMLElement) {
        if (!target.id) {
            target.id = 'tour-' + Math.random().toString(36).substr(2, 9)
        }
        target = `#${target.id}`
    }

    array.push({
        target,
        index,
        title,
        description
    })
}

export function Tour_Start() {
    GetTour(true).then(res => {
        if (res === false) {
            array.sort((a, b) => a.index - b.index)
            TourList.value = array
            IsOpenTour.value = true
        }
    })
}


let Setting = false
let Tools = false
let Jscript = false
let Dev = false
let isWindows = false
let xuHao = false

function check() {
    if (array.length >= 12) {
        if (Setting && Tools && Jscript) {
            if (isWindows && Dev) {
                Tour_Start()
                return
            } else if (!isWindows) {
                Tour_Start()
                return
            }
        }
    }
    requestAnimationFrame(() => {
        const arr = document.querySelectorAll('[data-ref="eToggleButton"]')
        for (let i = 0; i < arr.length; i++) {
            const text = arr[i].innerText;
            if (text.includes("设置") && !Setting) {
                Setting = true
                Tour_Add(arr[i], 4, "设置选项", "点击这里进行各种设置")
            }
            if (text.includes("常用工具") && !Tools) {
                Tools = true
                if (isWindows) {
                    Tour_Add(arr[i], 5, "常用工具", "这里有一些内置的工具\n\n你也可以添加一些工具到这里")
                } else {
                    Tour_Add(arr[i], 5, "常用工具", "这里有一些内置的工具")
                }
            }
            if (text.includes("脚本") && !Jscript) {
                Jscript = true
                Tour_Add(arr[i], 6, "脚本编辑", "在这里打开脚本编辑/脚本日志")
            }
            if (isWindows) {
                if (text.includes("驱动") && !Dev) {
                    Dev = true
                    Tour_Add(arr[i], 7, "驱动加载", "在这里加载驱动,加载后捕获进程数据")
                }
            }
        }
        if (!xuHao) {
            const arr2 = document.getElementsByClassName("ag-header ag-focus-managed ag-pivot-off ag-header-allow-overflow")
            for (let i = 0; i < arr2.length; i++) {
                const text = arr2[i].innerText;
                if (text.startsWith("序号\n") && arr2[i].clientHeight > 0) {
                    xuHao = true
                    Tour_Add(arr2[i], 11, "列编辑", "在这里 “右键” 点击\n可以隐藏某列，显示某列，调整顺序，重置列")
                }
            }
        }
        check()
    })
}

GOOS().then(_isWindows => {
    isWindows = _isWindows
    check()
})
