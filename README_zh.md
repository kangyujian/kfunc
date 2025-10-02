# kfunc · Go 工具平台（中文）

这是一个基于结构体标签的轻量级 Web 工具平台。你只需要用 `tag` 描述表单字段，平台即可自动生成页面表单、读取用户输入并绑定到结构体，调用你的业务方法并把结果显示到页面，同时支持按“业务空间”归类工具。

如需英文版，请查看根目录的 [README.md](README.md)。

---

## 快速开始

- 环境：Go 1.22+
- 运行：

```bash
cd kfunc
go run .
# 打开 http://localhost:8080
```

## 功能点

- 基于 `struct tag` 自动生成表单（字符串、整数/浮点、单选、下拉、多选、复选）。
- 从 Web 表单读取数据并自动绑定到结构体，完成类型转换。
- 为结构体绑定处理方法（业务自行实现），执行后将结果返回页面。
- 平台支持定义不同的业务空间；空间内可查看已存在的工具。

## 路由说明

- `/`：首页，展示所有业务空间
- `/spaces/{space}`：空间内的工具列表
- `/tools/{toolID}`：GET 渲染表单，POST 执行并展示结果

## 目录结构

```
internal/platform  平台核心（字段解析、数据绑定、工具注册与空间索引）
templates         页面模板（布局、空间列表、表单、结果）
tools/examples    示例工具（文本处理、计算器）
main.go           路由与入口，集中注册示例工具
```

## 表单标签规范

- `type`：`text | number | textarea | select | multiselect | radio | checkbox`
- `label`：字段展示名称
- `required`：`true | false`
- `options`：选项列表（`select`、`multiselect`、`radio`），使用 `|` 分隔
- `placeholder`：占位提示
- `default`：默认值
- `name`：表单字段名（不填则使用结构体字段名）

示例：

```go
type CalcForm struct {
    A      float64 `form:"type=number,label=A,required=true"`
    B      float64 `form:"type=number,label=B,required=true"`
    Op     string  `form:"type=radio,label=运算,options=加|减|乘|除,required=true"`
    Labels []string `form:"type=multiselect,label=标签,options=快速|准确|实验"`
}
```

## 新增一个业务工具

1) 在 `tools/yourpkg` 下创建你的包，定义表单结构体与实现 `FormTool` 的类型。
2) 在你的包中提供 `Register()` 方法，调用 `platform.RegisterTool(&YourTool{})`。
3) 在 `main.go` 中 `import yourpkg` 并调用 `yourpkg.Register()`。
4) 访问 `/spaces/{Space}` 即可看到你的工具。

## 常见问题

- 选择类控件（`select`、`radio`、`multiselect`）在模板迭代中要注意上下文切换；已在模板中通过局部变量方式修正单选框 `name` 绑定问题。
- 默认值与校验可在 `Run()` 中自定义，也可扩展 `form.go` 实现更丰富的校验能力。

> 欢迎继续提出需求：支持嵌套结构体、分组字段、校验提示、文件上传等，都可以在现有架构上扩展。