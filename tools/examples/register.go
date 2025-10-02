package examples

import "kfunc/internal/platform"

// Register 注册所有示例工具
func Register() {
    platform.RegisterTool(&TextTool{})
    platform.RegisterTool(&CalcTool{})
}