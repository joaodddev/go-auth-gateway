package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/joaodddev/go-auth-gateway/pkg/logger"
	"go.uber.org/zap"
)

type ReverseProxy struct {
	targets map[string]*httputil.ReverseProxy
}

func NewReverseProxy(serviceURLs map[string]string) (*ReverseProxy, error) {
	proxies := make(map[string]*httputil.ReverseProxy)

	for name, serviceURL := range serviceURLs {
		url, err := url.Parse(serviceURL)
		if err != nil {
			return nil, err
		}

		proxy := httputil.NewSingleHostReverseProxy(url)

		// Customize director to preserve original headers
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Header.Set("X-Proxy", "Go-Auth-Gateway")
			req.Host = url.Host
		}

		proxies[name] = proxy
	}

	return &ReverseProxy{
		targets: proxies,
	}, nil
}

func (rp *ReverseProxy) Route(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy, exists := rp.targets[serviceName]
		if !exists {
			logger.Log.Warn("Service not found", zap.String("service", serviceName))
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		// Remove service prefix from path
		path := strings.TrimPrefix(r.URL.Path, "/"+serviceName)
		r.URL.Path = path

		logger.Log.Debug("Proxying request",
			zap.String("service", serviceName),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
		)

		proxy.ServeHTTP(w, r)
	}
}
