package start

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status": "OK"}`))
}

func info(env *envContext) http.HandlerFunc {
	items := map[string]string{
		"name":      env.AppName,
		"version":   env.AppVersion,
		"commit":    env.AppCommit,
		"builtAt":   env.AppBuildAt,
		"startedAt": env.AppStartedAt,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		e := json.NewEncoder(w)
		_ = e.Encode(items)
	}
}

func Run(dirToServe string, addr string, tlsAddr string, certFile string, keyFile string, overridden []string) error {
	env := NewEnvContext(overridden)

	mux := chi.NewRouter()
	mux.Use(
		middleware.Recoverer,       // Recover errors.
		middleware.DefaultCompress, // GZIP
		middleware.RealIP,
		middleware.Logger,
	)

	configHandler, err := configJs(env)
	if err != nil {
		return err
	}

	mux.Get("/config.js", configHandler.ServeHTTP)
	// Admin

	mux.Get("/@/health", health)
	mux.Get("/@/info", info(env))

	if err := buildFolderListener(mux, env, dirToServe); err != nil {
		return err
	}

	errs := make(chan error)

	go func() {
		log.Println(fmt.Sprintf("Start to serve %s, binded to %s", dirToServe, addr))
		if err := http.ListenAndServe(addr, mux); err != nil {
			errs <- err
		}
	}()

	_, keyErr := os.Stat(keyFile)
	_, certErr := os.Stat(certFile)
	if keyErr == nil && certErr == nil {
		go func() {
			log.Println(fmt.Sprintf("Start to serve %s, binded to %s", dirToServe, tlsAddr))
			if err := http.ListenAndServeTLS(tlsAddr, certFile, keyFile, mux); err != nil {
				errs <- err
			}
		}()
	}

	return <-errs
}

func buildFolderListener(mux *chi.Mux, env *envContext, folder string) error {
	dir := http.Dir(folder)
	base := filepath.Clean(folder)

	var rootFile *fileHandler = nil
	// Explore the folder to create specific route for each files under the tree.
	res := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		httpPath := strings.Replace(strings.TrimPrefix(path, base), string(os.PathSeparator), "/", -1)
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			file, err := dir.Open(httpPath)
			s1 := sha1.New()
			_, err = s1.Write(b)
			if err != nil {
				return err
			}

			noCache := false
			for _, item := range env.NoCacheFiles {
				if item == httpPath {
					noCache = true
					break
				}
			}
			h := &fileHandler{
				Name:    info.Name(),
				modtime: info.ModTime(),
				file:    file,
				path:    path,
				noCache: noCache,
				hash:    fmt.Sprintf("%x", s1.Sum(nil)),
			}
			if httpPath == "/index.html" {
				rootFile = h
			} else if strings.HasSuffix(httpPath, "index.html") {
				mux.Handle(strings.TrimSuffix(httpPath, "index.html"), h)
				mux.Handle(strings.TrimSuffix(httpPath, "/index.html"), h)
			}
			mux.Handle(httpPath, h)
		}
		return nil
	})
	if rootFile != nil {
		mux.Handle("/*", rootFile)
	}
	return res
}

type fileHandler struct {
	Name    string
	path    string
	modtime time.Time
	file    http.File
	noCache bool
	hash    string
}

func (t *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO Inspect index.html to send all resources via http2
	/*
	push, isPusher := w.(http.Pusher)
	if isPusher {
		// Do Something
	}
	 */
	if t.noCache {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Expires", "0")
		w.Header().Add("Pragma", "no-cache")
	} else {
		w.Header().Add("ETag", t.hash)
	}
	http.ServeContent(w, r, t.Name, t.modtime, t.file)
}
