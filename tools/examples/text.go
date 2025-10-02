package examples

import (
    "context"
    "fmt"
    "strings"
)

// TextTool 示例工具
type TextForm struct {
    Content string `form:"type=textarea,label=文本内容,placeholder=在此输入文本,required=true"`
    Action  string `form:"type=select,label=操作,options=Upper|Lower|Title,required=true"`
}

type TextTool struct{}

func (t *TextTool) ID() string          { return "text_tool" }
func (t *TextTool) Name() string        { return "文本处理工具" }
func (t *TextTool) Description() string { return "对文本内容进行大小写转换等操作" }
func (t *TextTool) Space() string       { return "content" }
func (t *TextTool) FormStruct() any     { return &TextForm{} }
func (t *TextTool) Run(ctx context.Context, form any) (any, error) {
    f := form.(*TextForm)
    switch f.Action {
    case "Upper":
        return strings.ToUpper(f.Content), nil
    case "Lower":
        return strings.ToLower(f.Content), nil
    case "Title":
        return strings.Title(f.Content), nil
    default:
        return nil, fmt.Errorf("未知操作: %s", f.Action)
    }
}