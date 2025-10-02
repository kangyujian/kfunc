package platform

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "net/url"
)

// FieldSpec 用于生成表单的字段描述
type FieldSpec struct {
    Name         string
    Label        string
    Type         string // text, number, textarea, select, multiselect, radio, checkbox
    Required     bool
    Options      []string
    Placeholder  string
    DefaultValue string
}

// ExtractFields 根据结构体的 tag 解析出 FieldSpec 列表
func ExtractFields(v any) []FieldSpec {
    rv := reflect.ValueOf(v)
    if rv.Kind() == reflect.Ptr {
        rv = rv.Elem()
    }
    rt := rv.Type()
    fields := make([]FieldSpec, 0, rt.NumField())
    for i := 0; i < rt.NumField(); i++ {
        sf := rt.Field(i)
        if sf.PkgPath != "" { // 非导出字段跳过
            continue
        }
        tag := sf.Tag.Get("form")
        fs := parseTag(sf.Name, tag)
        fields = append(fields, fs)
    }
    return fields
}

func parseTag(fieldName, tag string) FieldSpec {
    fs := FieldSpec{Name: fieldName, Type: inferTypeByKind(fieldName)}
    if tag == "" {
        return fs
    }
    parts := strings.Split(tag, ",")
    for _, p := range parts {
        kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
        if len(kv) != 2 {
            continue
        }
        k, v := kv[0], kv[1]
        switch k {
        case "name":
            fs.Name = v
        case "label":
            fs.Label = v
        case "type":
            fs.Type = v
        case "required":
            fs.Required = v == "true"
        case "options":
            if v != "" {
                fs.Options = strings.Split(v, "|")
            }
        case "placeholder":
            fs.Placeholder = v
        case "default":
            fs.DefaultValue = v
        }
    }
    return fs
}

func inferTypeByKind(fieldName string) string { return "text" }

// BindFormValues 将 url.Values 填充到结构体实例中
func BindFormValues(dst any, values url.Values) error {
    rv := reflect.ValueOf(dst)
    if rv.Kind() != reflect.Ptr {
        return fmt.Errorf("dst 必须是指针")
    }
    rv = rv.Elem()
    rt := rv.Type()
    for i := 0; i < rv.NumField(); i++ {
        sf := rt.Field(i)
        if sf.PkgPath != "" { // 非导出字段跳过
            continue
        }
        tag := sf.Tag.Get("form")
        spec := parseTag(sf.Name, tag)
        name := spec.Name
        vals, ok := values[name]
        if !ok || len(vals) == 0 {
            continue
        }
        fv := rv.Field(i)
        if !fv.CanSet() {
            continue
        }
        switch fv.Kind() {
        case reflect.String:
            fv.SetString(vals[0])
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            iv, err := strconv.ParseInt(vals[0], 10, 64)
            if err != nil { return err }
            fv.SetInt(iv)
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            uv, err := strconv.ParseUint(vals[0], 10, 64)
            if err != nil { return err }
            fv.SetUint(uv)
        case reflect.Float32, reflect.Float64:
            fv64, err := strconv.ParseFloat(vals[0], 64)
            if err != nil { return err }
            fv.SetFloat(fv64)
        case reflect.Slice:
            // 仅支持 []string
            if fv.Type().Elem().Kind() == reflect.String {
                fv.Set(reflect.ValueOf(vals))
            }
        case reflect.Bool:
            bv := vals[0] == "on" || vals[0] == "true"
            fv.SetBool(bv)
        default:
            // 暂不支持嵌套
        }
    }
    return nil
}