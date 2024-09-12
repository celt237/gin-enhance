package tmp

var Gin = `
// Code generated by gin-enhance. DO NOT EDIT.
package {{.PackageName}}

import (
{{- range .Imports}}
	{{if .HasAlias}}{{.Name}} {{end -}}"{{.Path}}"
{{- end}}{{- range .ExtraImports}}
	_ "{{.}}"
{{- end}}
)

{{$svrType := .Name}}
{{$empty := "emptypb.Empty"}}
// declare handler interface
type {{$svrType}}HTTPHandler interface {
{{- range .Methods}}
    {{.Name}}(ctx context.Context{{- range .Params}}, {{.Name | firstToLower}} {{.DataType}}{{- end}}) ({{.Result.DataType}}, error)
{{- end}}
}

// Register{{$svrType}}HTTPHandler define http router handle by gin. 
// regist route handler
func Register{{$svrType}}HTTPHandler(group *gin.RouterGroup, api gin_enhance.ApiHandler, srv {{$svrType}}HTTPHandler) {
    {{- range .Methods}}
    group.{{.Method | toUpper}}("{{.Path}}", _{{$svrType}}_{{.Name}}_HTTP_Handler(api, srv))
    {{- end}}
}

// declare handler
// Traverse all previously parsed rpc method information

{{range $outerIndex, $outerElement := .Methods}}
// @Summary {{.Summary}}{{if ne .Description ""}}
// @Description {{.Description}}{{end}}{{if ne .Tags ""}}
// @Tags {{.Tags}}{{end}}
// @Accept {{.Accept}}
// @Produce {{.Produce}}{{range .Params}}
// @Param {{.Name}} {{.ParamType}} {{.ParamDataType}} {{.Required}} {{.Desc}}{{end}}
// @Success 200 {object} {{.ApiResultType}}[{{.ApiResultDataType}}]  "请求成功返回的结果"
// @Failure {{.ErrorCode}} {object} {{.ApiResultType}}[{{.ApiResultDataType}}] "请求失败返回的结果"
// @Router {{.Path}} [{{.Method}}]
func _{{$svrType}}_{{.Name}}_HTTP_Handler(api gin_enhance.ApiHandler, srv {{$svrType}}HTTPHandler) gin.HandlerFunc {
    return func(c *gin.Context) {
        wrapperCtx := api.WrapContext(c)
		var resp {{.Result.DataType}}
		{{range .Params}}
		{{.Name}}, err := gin_enhance.GetParamFromContext[{{.DataType}}](c, "{{.Name}}", "{{.DataType}}", "{{.ParamType}}", {{if .IsPtr}}true{{else}}false{{end -}}, {{if .Required}}true{{else}}false{{end -}})
		if err != nil {
			api.Error(c, "{{$outerElement.Produce}}", resp, err)
			return
		}{{end}}
		{{range $key, $value := .CustomAnnotations}}{{range $key2, $value2 :=  .Attributes}}err = api.HandleCustomerAnnotation(c, "{{$key}}"{{range $k, $v := $value2 }}, "{{$v}}"{{end}})
		if err != nil {
			api.Error(c, "{{$outerElement.Produce}}", resp, err)
			return
		}
		{{end}}{{end}}
        // 执行方法
        resp, err = srv.{{.Name}}(wrapperCtx{{- range $index, $value := .Params}}, {{.Name | firstToLower}}{{- end}})
        if err != nil {
            api.Error(c, "{{$outerElement.Produce}}", resp, err)
            return
        }
        api.Success(c, "{{$outerElement.Produce}}", resp)
	}
}
{{end}}

`
