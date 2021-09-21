package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/assets"
)

type fileCache struct {
	data []byte
	eTag string
}
type cacheMap map[string]*fileCache

var assetCache = cacheMap{}

func (x cacheMap) Read(fname string) *fileCache {
	asset := assets.Assets()
	c, ok := x[fname]
	if !ok {
		data, err := asset.ReadFile(filepath.Join("out", fname))
		if err != nil {
			logger.Debug().Str("path", fname).Err(err).Msg("failed to open requested file")
			return nil
		}

		c = &fileCache{
			data: data,
			eTag: fmt.Sprintf("%x", sha256.Sum256(data)),
		}
	}

	return c
}

type extTypeMap map[string]string

var extMap = extTypeMap{
	".html": "text/html",
	".js":   "application/javascript",
}

func (x extTypeMap) Find(path string) string {
	for ext, contentType := range x {
		if strings.HasSuffix(path, ext) {
			return contentType
		}
	}

	return "text/html"
}

type rewriteRoute struct {
	ptn   *regexp.Regexp
	fname string
}

type rewriteRoutes []*rewriteRoute

func (x rewriteRoutes) Rewrite(path string) string {
	for _, route := range x {
		if route.ptn.MatchString(path) {
			return route.fname
		}
	}
	return path
}

var nextRoutes = rewriteRoutes{
	{
		ptn:   regexp.MustCompile("^scan/[a-z0-9-]+$"),
		fname: "scan/[id].html",
	},
}

func getStaticFile(ctx *gin.Context) {
	ctx.Next()

	if ctx.Writer.Status() != http.StatusNotFound {
		return
	}

	fname := nextRoutes.Rewrite(strings.Trim(ctx.Request.URL.Path, "/"))
	logger.Debug().Interface("req", ctx.Request.URL).Str("rewrite", fname).Msg("accessing static file")

	if fname == "" {
		fname = "index.html"
	}

	cache := assetCache.Read(fname)
	if cache == nil {
		return
	}

	ctx.Header("Cache-Control", "public, max-age=31536000")
	ctx.Header("ETag", cache.eTag)

	if match := ctx.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, cache.eTag) {
			ctx.Status(http.StatusNotModified)
			return
		}
	}

	ctx.Data(http.StatusOK, extMap.Find(fname), cache.data)
}
