# kfunc · Golang Tool Platform

A lightweight, tag-driven web platform to create and run business tools in Go. Define a struct with tags, auto-generate the web form, bind submitted data back to the struct, run your logic, and return results to the page. Supports organizing tools by business spaces.

Looking for Chinese docs? See: [中文说明 / README_zh.md](README_zh.md)

---

## Quick Start

- Requirements: Go 1.22+
- Run:

```bash
cd kfunc
go run .
# Open http://localhost:8080
```

## Features

- Tag-driven form generation from Go structs.
- Bind `url.Values` → struct with type conversion (string/int/float/bool/[]string).
- Tool interface to plug in custom logic; results rendered on a result page.
- Business spaces to group tools; see tools per space.
- Clean templates; minimal dependencies.

## Routes

- `/` — Home, list spaces
- `/spaces/{space}` — Tools in the space
- `/tools/{toolID}` — GET form, POST execute and show result

## Project Structure

```
.
├── internal/platform          # platform core
│   ├── form.go                # tags → fields, values → struct
│   └── registry.go            # tool registry, space index
├── templates                  # HTML templates
│   ├── layout.html            # base layout & styles
│   ├── index.html             # spaces list
│   ├── space.html             # tools list in a space
│   ├── form.html              # dynamic form
│   └── result.html            # result display
├── tools/examples             # example tools
│   ├── register.go
│   ├── text.go
│   └── calc.go
└── main.go                    # HTTP routes & boot
```

## Define a Tool (Interface)

A tool must implement:

```go
// internal/platform/registry.go
type FormTool interface {
    ID() string
    Name() string
    Description() string
    Space() string
    FormStruct() any
    Run(ctx context.Context, form any) (any, error)
}
```

## Struct Tag Spec (Form Tags)

Tags are key=value pairs separated by commas on exported fields. Supported keys:

- `type`: `text | number | textarea | select | multiselect | radio | checkbox`
- `label`: label text shown in the form
- `required`: `true | false`
- `options`: option list, separated by `|` (for `select`, `multiselect`, `radio`)
- `placeholder`: placeholder text
- `default`: default value
- `name`: field name in the form (override struct field name)

Example:

```go
// tools/examples/text.go
type TextForm struct {
    Content string `form:"type=textarea,label=Text,placeholder=Type here,required=true"`
    Action  string `form:"type=select,label=Action,options=Upper|Lower|Title,required=true"`
}
```

## Example Tool

```go
// tools/examples/text.go
type TextTool struct{}
func (t *TextTool) ID() string          { return "text_tool" }
func (t *TextTool) Name() string        { return "Text Processor" }
func (t *TextTool) Description() string { return "String case transforms" }
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
        return nil, fmt.Errorf("unknown action: %s", f.Action)
    }
}
```

Register tools (centralized):

```go
// tools/examples/register.go
func Register() {
    platform.RegisterTool(&TextTool{})
    platform.RegisterTool(&CalcTool{})
}

// main.go
examples.Register()
```

## Add Your Own Tool

1) Create a new package under `tools/yourpkg`, define form struct with `form` tags and a type implementing `FormTool`.
2) Add a `Register()` function in your package to call `platform.RegisterTool(...)` for each tool.
3) Import and call your `yourpkg.Register()` from `main.go`.
4) Visit `/spaces/{yourSpace}` to find your tool.

## Notes & Tips

- Only exported struct fields are read and bound.
- `multiselect` binds to `[]string`.
- `number` binds to `float64` or any numeric kinds; parsing errors return 400.
- Validation can be enhanced in `Run()` or by extending `form.go`.
 
---
 
## Framework & Run Notes
 
 - Routing is implemented with Gin (v1.9.1); HTML is rendered via `html/template` using embedded templates (`embed.FS`).
 - Port: default is `8080`. You can override via `PORT` env: `PORT=8081 go run .`.
 - Mainland China users: set Go module proxy if needed: `go env -w GOPROXY=https://goproxy.cn,direct`.
 - Trusted proxies: in production, configure explicitly to avoid trusting all proxies:
 
 ```go
 r := gin.Default()
 r.SetTrustedProxies(nil) // or set specific CIDR ranges
 ```