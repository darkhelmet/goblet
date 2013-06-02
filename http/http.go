package http

import (
    "fmt"
    G "github.com/darkhelmet/goblet"
    "net/http"
)

type Handler struct {
    goblet *G.Goblet
    paths  map[string]*G.Asset
}

func NewHandler(prefix string, g *G.Goblet) http.Handler {
    paths := make(map[string]*G.Asset)
    for name, asset := range g.Files {
        paths[fmt.Sprintf("%s/%s", prefix, name)] = asset
    }
    return Handler{goblet: g, paths: paths}
}

func (h Handler) get(path string) *G.Asset {
    return h.paths[path]
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    asset := h.get(req.URL.Path)
    if asset == nil {
        http.NotFound(w, req)
        return
    }

    w.Header().Set("ETag", fmt.Sprintf(`"%s"`, asset.Sha1))

    http.ServeContent(w, req, req.URL.Path, asset.LastModified, asset.Reader())
}
