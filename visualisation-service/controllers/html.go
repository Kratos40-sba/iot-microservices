package controllers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)
func Home(x *gin.Context)  {
	err := renderPage(x.Writer,"home.jet",nil)
	if err != nil {
		log.Println(err)
	}
}

func renderPage(w http.ResponseWriter, tmpl string , data jet.VarMap) error{
	v , err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	err = v.Execute(w,data,nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
