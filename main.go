package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path"

    "github.com/gin-gonic/gin"
    examples "kfunc/tools/examples"
    "kfunc/view"
)

//go:embed templates/*.html
var tmplFS embed.FS

func main() {
    // Register example tools from examples package
    examples.Register()

    r := gin.Default()

    // Provide renderer from main (embed templates here)
    view.RegisterRoutes(r, render)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    addr := ":" + port
    log.Printf("Listening on http://localhost%v\n", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal(err)
    }
}

func render(w http.ResponseWriter, name string, data any) {
    funcMap := template.FuncMap{
        "toJSON": func(v any) template.JS {
            b, err := json.Marshal(v)
            if err != nil {
                return template.JS("null")
            }
            return template.JS(b)
        },
    }
    t, err := template.New("layout").Funcs(funcMap).ParseFS(tmplFS, "templates/layout.html", path.Join("templates", name))
    if err != nil {
        http.Error(w, fmt.Sprintf("模板错误: %v", err), http.StatusInternalServerError)
        return
    }
    if err := t.ExecuteTemplate(w, "layout", data); err != nil {
        http.Error(w, fmt.Sprintf("渲染错误: %v", err), http.StatusInternalServerError)
        return
    }
}