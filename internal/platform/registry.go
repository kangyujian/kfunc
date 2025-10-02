package platform

import (
    "context"
    "sync"
)

// FormTool 表示一个业务工具，提供元信息、表单结构体、以及处理方法
type FormTool interface {
    ID() string
    Name() string
    Description() string
    Space() string
    FormStruct() any
    Run(ctx context.Context, form any) (any, error)
}

var (
    mu       sync.RWMutex
    tools    = make(map[string]FormTool)
    bySpace  = make(map[string][]FormTool)
)

// RegisterTool 注册一个工具到平台
func RegisterTool(t FormTool) {
    mu.Lock()
    defer mu.Unlock()
    tools[t.ID()] = t
    s := t.Space()
    bySpace[s] = append(bySpace[s], t)
}

// GetTool 通过ID获取工具
func GetTool(id string) FormTool {
    mu.RLock()
    defer mu.RUnlock()
    return tools[id]
}

// ListSpaces 返回所有业务空间名称
func ListSpaces() []string {
    mu.RLock()
    defer mu.RUnlock()
    spaces := make([]string, 0, len(bySpace))
    for s := range bySpace {
        spaces = append(spaces, s)
    }
    return spaces
}

// ListToolsBySpace 列出指定空间下的工具
func ListToolsBySpace(space string) []FormTool {
    mu.RLock()
    defer mu.RUnlock()
    return bySpace[space]
}