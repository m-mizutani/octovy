package api

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Assets
func getIndex(c *gin.Context) {
	config := getConfig(c)
	data, err := ioutil.ReadFile(filepath.Join(config.AssetDir, "index.html"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.Data(200, "text/html", data)
}

func getBundleJS(c *gin.Context) {
	config := getConfig(c)
	data, err := ioutil.ReadFile(filepath.Join(config.AssetDir, "bundle.js"))
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}
	c.Data(200, "application/javascript", data)
}
