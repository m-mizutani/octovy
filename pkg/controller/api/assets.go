package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/assets"
	"github.com/m-mizutani/octovy/pkg/domain/model"
)

type cache struct {
	data []byte
	eTag string
}
type cacheMap map[string]*cache

var assetCache = cacheMap{}

func initAsset() {
	assets := assets.Assets()

	indexHTML, err := assets.ReadFile("dist/index.html")
	if err != nil {
		panic("Open dist/index.html: " + err.Error())
	}
	bundleJS, err := assets.ReadFile("dist/bundle.js")
	if err != nil {
		panic("Open dist/bundle.js: " + err.Error())
	}

	assetCache["index.html"] = &cache{
		data: indexHTML,
		eTag: fmt.Sprintf("%x", sha256.Sum256(indexHTML)),
	}

	assetCache["bundle.js"] = &cache{
		data: bundleJS,
		eTag: fmt.Sprintf("%x", sha256.Sum256(bundleJS)),
	}
}

func handleAsset(ctx *gin.Context, fname, contentType string) {
	c, ok := assetCache[fname]
	if !ok {
		ctx.Error(goerr.Wrap(model.ErrItemNotFound))
		return
	}

	ctx.Header("Cache-Control", "public, max-age=31536000")
	ctx.Header("ETag", c.eTag)

	if match := ctx.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, c.eTag) {
			ctx.Status(http.StatusNotModified)
			return
		}
	}

	ctx.Data(http.StatusOK, contentType, c.data)
}

// Assets
func getIndex(c *gin.Context) {
	handleAsset(c, "index.html", "text/html")
}

func getBundleJS(c *gin.Context) {
	handleAsset(c, "bundle.js", "application/javascript")
}
