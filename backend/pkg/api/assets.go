package api

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type cache struct {
	data []byte
	eTag string
}

type cacheMap map[string]*cache

func (x cacheMap) load(fpath string) (*cache, error) {
	c, ok := x[fpath]
	if !ok {
		data, err := ioutil.ReadFile(fpath)
		if err != nil {
			return nil, err
		}
		c = &cache{
			data: data,
			eTag: fmt.Sprintf("%x", md5.Sum(data)),
		}
		x[fpath] = c
	}
	return c, nil
}

var assetCache = cacheMap{}

func handleAsset(ctx *gin.Context, fname, contentType string) {
	config := getConfig(ctx)
	c, err := assetCache.load(filepath.Join(config.AssetDir, fname))
	if err != nil {
		_ = ctx.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	ctx.Header("Cache-Control", "public, max-age=31536000")
	ctx.Header("ETag", c.eTag)

	if match := ctx.GetHeader("If-None-Match"); match != "" {
		logger.With("match", match).Info("if-none-match")
		// TODO: Fix if-none-match parsing
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
