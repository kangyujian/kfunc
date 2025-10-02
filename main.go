package main

import (
    "context"
    "embed"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "path"
    "strings"

    "kfunc/internal/platform"
    examples "kfunc/tools/examples"
)

//go:embed templates/*.html
var tmplFS embed.FS

func main() {
    // Register example tools from examples package
    examples.Register()

    mux := http.NewServeMux()

    // Routes
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        spaces := platform.ListSpaces()
        render(w, "index.html", map[string]any{
            "Spaces": spaces,
        })
    })

    mux.HandleFunc("/spaces/", func(w http.ResponseWriter, r *http.Request) {
        // /spaces/{space}
        p := strings.TrimPrefix(r.URL.Path, "/spaces/")
        if p == "" {
            http.NotFound(w, r)
            return
        }
        tools := platform.ListToolsBySpace(p)
        render(w, "space.html", map[string]any{
            "Space": p,
            "Tools": tools,
        })
    })

    mux.HandleFunc("/tools/", func(w http.ResponseWriter, r *http.Request) {
        // /tools/{id}
        id := strings.TrimPrefix(r.URL.Path, "/tools/")
        t := platform.GetTool(id)
        if t == nil {
            http.NotFound(w, r)
            return
        }
        switch r.Method {
        case http.MethodGet:
            form := t.FormStruct()
            fields := platform.ExtractFields(form)
            render(w, "form.html", map[string]any{
                "Tool":   t,
                "Fields": fields,
                "Action": path.Join("/tools", t.ID()),
            })
        case http.MethodPost:
            form := t.FormStruct()
            if err := r.ParseForm(); err != nil {
                http.Error(w, "无法解析表单", http.StatusBadRequest)
                return
            }
            if err := platform.BindFormValues(form, r.Form); err != nil {
                http.Error(w, fmt.Sprintf("绑定表单失败: %v", err), http.StatusBadRequest)
                return
            }
            res, err := t.Run(context.Background(), form)
            if err != nil {
                http.Error(w, fmt.Sprintf("处理失败: %v", err), http.StatusInternalServerError)
                return
            }
            render(w, "result.html", map[string]any{
                "Tool":   t,
                "Result": res,
            })
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    addr := ":8080"
    log.Printf("Listening on http://localhost%v\n", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}

func render(w http.ResponseWriter, name string, data any) {
    t, err := template.ParseFS(tmplFS, "templates/layout.html", path.Join("templates", name))
    if err != nil {
        http.Error(w, fmt.Sprintf("模板错误: %v", err), http.StatusInternalServerError)
        return
    }
    if err := t.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, fmt.Sprintf("渲染错误: %v", err), http.StatusInternalServerError)
        return
    }
}