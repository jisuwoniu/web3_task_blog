package handler

import (
	"html/template"
)

var htmlTemplates = make(map[string]*template.Template)

// InitTemplates 初始化HTML模板
func InitTemplates() {
	// 这里将在实际应用中加载模板文件
	// 现在我们只是初始化map
	//htmlTemplates["index.html"] = template.Must(template.ParseFiles("web/static/index.html"))
	//htmlTemplates["hello.html"] = template.Must(template.ParseFiles("web/static/hello.html"))
	//htmlTemplates["login.html"] = template.Must(template.ParseFiles("web/static/login.html"))
	//htmlTemplates["register.html"] = template.Must(template.ParseFiles("web/static/register.html"))
	//htmlTemplates["post.html"] = template.Must(template.ParseFiles("web/static/post.html"))
}