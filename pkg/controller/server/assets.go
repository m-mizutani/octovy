package server

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/octovy/assets"
	"github.com/m-mizutani/octovy/pkg/utils"
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
			utils.Logger.With("path", fname).With("error", err).Debug("failed to open requested file")
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
	{
		ptn:   regexp.MustCompile("^vulnerability$"),
		fname: "vulnerability.html",
	},
	{
		ptn:   regexp.MustCompile("^vulnerability/[A-Za-z0-9-.]+$"),
		fname: "vulnerability/[id].html",
	},
	{
		ptn:   regexp.MustCompile("^login$"),
		fname: "login.html",
	},
	{
		ptn:   regexp.MustCompile("^config$"),
		fname: "config.html",
	},
}

func getStaticFile(c *gin.Context) {
	c.Next()
	logger := getLog(c)

	if c.Writer.Status() != http.StatusNotFound {
		return
	}

	fname := nextRoutes.Rewrite(strings.Trim(c.Request.URL.Path, "/"))
	logger.With("req", c.Request.URL).With("rewrite", fname).Debug("accessing static file")

	if fname == "" {
		fname = "index.html"
	}

	cache := assetCache.Read(fname)
	if cache == nil {
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("ETag", cache.eTag)

	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, cache.eTag) {
			c.Status(http.StatusNotModified)
			return
		}
	}

	c.Data(http.StatusOK, extMap.Find(fname), cache.data)
}
