package view

import (
    "context"
    "fmt"
    "net/http"
    "path"

    "github.com/gin-gonic/gin"
    "kfunc/internal/platform"
)

// Renderer is a function that renders the given template name with data
type Renderer func(http.ResponseWriter, string, any)

// RegisterRoutes sets up Gin handlers for views
func RegisterRoutes(r *gin.Engine, render Renderer) {
    r.GET("/", func(c *gin.Context) {
        spaces := platform.ListSpaces()
        render(c.Writer, "index.html", map[string]any{
            "Spaces": spaces,
        })
    })

    r.GET("/spaces/:space", func(c *gin.Context) {
        p := c.Param("space")
        if p == "" {
            c.Status(http.StatusNotFound)
            return
        }
        tools := platform.ListToolsBySpace(p)
        render(c.Writer, "space.html", map[string]any{
            "Space": p,
            "Tools": tools,
        })
    })

    r.GET("/tools/:id", func(c *gin.Context) {
        id := c.Param("id")
        t := platform.GetTool(id)
        if t == nil {
            c.Status(http.StatusNotFound)
            return
        }
        form := t.FormStruct()
        fields := platform.ExtractFields(form)
        render(c.Writer, "form.html", map[string]any{
            "Tool":   t,
            "Fields": fields,
            "Action": path.Join("/tools", t.ID()),
        })
    })

    r.POST("/tools/:id", func(c *gin.Context) {
        id := c.Param("id")
        t := platform.GetTool(id)
        if t == nil {
            c.Status(http.StatusNotFound)
            return
        }
        form := t.FormStruct()
        if err := c.Request.ParseForm(); err != nil {
            http.Error(c.Writer, "无法解析表单", http.StatusBadRequest)
            return
        }
        if err := platform.BindFormValues(form, c.Request.Form); err != nil {
            http.Error(c.Writer, fmt.Sprintf("绑定表单失败: %v", err), http.StatusBadRequest)
            return
        }
        res, err := t.Run(context.Background(), form)
        if err != nil {
            http.Error(c.Writer, fmt.Sprintf("处理失败: %v", err), http.StatusInternalServerError)
            return
        }
        render(c.Writer, "result.html", map[string]any{
            "Tool":   t,
            "Result": res,
        })
    })
}