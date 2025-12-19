package app

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		log.WithFields(log.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"remote":   r.RemoteAddr,
			"status":   rec.status,
			"size":     rec.size,
			"duration": time.Since(start).String(),
		}).Info("request completed")
	})
}

func newReverseProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.WithError(err).Fatalf("invalid upstream url: %s", target)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		// preserve original host header of upstream
		req.Host = u.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.WithError(err).WithField("upstream", target).Error("upstream error")
		w.WriteHeader(http.StatusBadGateway)
		io.WriteString(w, "bad gateway")
	}

	return proxy
}

func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// proxies
	orderProxy := newReverseProxy(cfg.Upstreams.Order)
	paymentProxy := newReverseProxy(cfg.Upstreams.Payment)
	analyticsProxy := newReverseProxy(cfg.Upstreams.Analytics)
	menuProxy := newReverseProxy(cfg.Upstreams.Menu)

	mux.Handle("/dishes/", menuProxy)

	mux.Handle("/orders", orderProxy)
	mux.Handle("/orders/", orderProxy)

	mux.Handle("/payments", paymentProxy)
	mux.Handle("/payments/", paymentProxy)

	mux.Handle("/analytics", analyticsProxy)
	mux.Handle("/analytics/", analyticsProxy)

	// catch-all: return 404 for unknown
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return loggingMiddleware(mux)
}
