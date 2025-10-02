package examples

import (
    "context"
    "fmt"
)

// CalcTool 示例工具
type CalcForm struct {
    A      float64 `form:"type=number,label=A,required=true"`
    B      float64 `form:"type=number,label=B,required=true"`
    Op     string  `form:"type=radio,label=运算,options=加|减|乘|除,required=true"`
    Labels []string `form:"type=multiselect,label=标签,options=快速|准确|实验"`
}

type CalcTool struct{}

func (c *CalcTool) ID() string          { return "calc_tool" }
func (c *CalcTool) Name() string        { return "简易计算器" }
func (c *CalcTool) Description() string { return "执行基础四则运算并附加标签" }
func (c *CalcTool) Space() string       { return "math" }
func (c *CalcTool) FormStruct() any     { return &CalcForm{} }
func (c *CalcTool) Run(ctx context.Context, form any) (any, error) {
    f := form.(*CalcForm)
    var res float64
    switch f.Op {
    case "加":
        res = f.A + f.B
    case "减":
        res = f.A - f.B
    case "乘":
        res = f.A * f.B
    case "除":
        if f.B == 0 {
            return nil, fmt.Errorf("除数不能为0")
        }
        res = f.A / f.B
    default:
        return nil, fmt.Errorf("未知运算: %s", f.Op)
    }
    return map[string]any{
        "result":  res,
        "labels":  f.Labels,
        "explain": fmt.Sprintf("%v %s %v = %v", f.A, f.Op, f.B, res),
    }, nil
}