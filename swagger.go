package gin_enhance

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	//go:embed static
	front embed.FS
	s     service
)

type Config struct {
	RelativePath string
	DocJson      string
}

type service struct {
	Name           string `json:"name"`
	Url            string `json:"url"`
	SwaggerVersion string `json:"swaggerVersion"`
	Location       string `json:"location"`
}

func readFileFromEmbedFS(fs embed.FS, filename string) (string, error) {
	file, err := fs.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// RegisterSwaggerDoc registers swagger documentation
//
// app : the iris application
//
// jsonPath: the path to the swagger json file (e.g. ./docs/swagger.json)
//
// route: the path to register the swagger documentation (e.g. /doc)
//
// return: the path of the swagger documentation
func RegisterSwaggerDoc(app *gin.Engine, jsonPath string, route string) {
	route = strings.TrimPrefix(route, "/")
	route = strings.TrimSuffix(route, "/")
	route = strings.TrimSpace(route)
	if route == "" {
		log.Println("route is empty")
		return
	}
	route = "/" + route

	if jsonPath == "" {
		log.Println(jsonPath + "jsonPath is empty")
		return
	}
	docJson, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Println("no swagger.json found in " + jsonPath)
		return
	}
	strDocJson := string(docJson)
	app.GET(route+"/{any:path}", swagDocHandler(Config{RelativePath: route, DocJson: strDocJson}))
}

// swagDocHandler is a handler for swagger documentation
func swagDocHandler(config Config) gin.HandlerFunc {
	docJsonPath := config.RelativePath + "/docJson"
	servicesPath := config.RelativePath + "/static/service"
	docPath := config.RelativePath + "/index"
	appjsPath := config.RelativePath + "/static/webjars/js/app.42aa019b.js"

	s.Url = "/docJson"
	s.Location = "/docJson"
	s.Name = "API Documentation"
	s.SwaggerVersion = "2.0"

	appjs, err := readFileFromEmbedFS(front, "static/webjars/js/app.42aa019b.js")
	if err != nil {
		log.Println(err)
	}
	appjs = strings.ReplaceAll(appjs, "{[(.RelativePath)]}", config.RelativePath)

	doc, err := readFileFromEmbedFS(front, "static/doc.html")
	if err != nil {
		log.Println(err)
	}
	doc = strings.ReplaceAll(doc, "{[(.RelativePath)]}", config.RelativePath)

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}
		switch ctx.Request.RequestURI {
		case appjsPath:
			ctx.String(http.StatusOK, appjs)
		case servicesPath:
			ctx.JSON(http.StatusOK, []service{s})
		case docPath:
			ctx.HTML(http.StatusOK, "doc", doc)
		case docJsonPath:
			ctx.Header("Content-Type", "application/json")
			ctx.String(http.StatusOK, string(config.DocJson))
		default:
			filePath := strings.TrimPrefix(ctx.Request.RequestURI, config.RelativePath)
			filePath = strings.TrimPrefix(filePath, "/")
			file, err := front.Open(filePath)
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Error while opening file: %v", err)
				return
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				ctx.String(http.StatusInternalServerError, "Error while getting file info: %v", err)
				return
			}
			//ctx.ServeContent(reader, fileInfo.Name(), fileInfo.ModTime())
			ctx.DataFromReader(http.StatusOK, fileInfo.Size(), mime.TypeByExtension(filepath.Ext(fileInfo.Name())),
				file,
				map[string]string{
					"Content-Disposition": fmt.Sprintf("attachment; filename=%s", fileInfo.Name()),
					"Last-Modified":       fileInfo.ModTime().UTC().Format(http.TimeFormat),
				})
		}
		ctx.Next()
	}

}
