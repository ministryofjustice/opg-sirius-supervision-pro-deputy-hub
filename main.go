package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/logging"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/server"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-supervision-pro-deputy-hub/internal/util"
)

func main() {
	logger := logging.New(os.Stdout, "opg-sirius-pro-deputy-hub ")

	port := getEnv("PORT", "1234")
	webDir := getEnv("WEB_DIR", "web")
	siriusURL := getEnv("SIRIUS_URL", "http://localhost:8080")
	siriusPublicURL := getEnv("SIRIUS_PUBLIC_URL", "")
	firmHubURL := getEnv("FIRM_HUB_URL", "/supervision/deputies/firm")
	prefix := getEnv("PREFIX", "")

	layouts, _ := template.
		New("").
		Funcs(map[string]interface{}{
			"prefix": func(s string) string {
				return prefix + s
			},
			"sirius": func(s string) string {
				return siriusPublicURL + s
			},
			"firmhub": func(s string) string {
				return firmHubURL + s
			},
			"translate": util.Translate,
		}).
		ParseGlob(webDir + "/template/*/*.gotmpl")

	files, _ := filepath.Glob(webDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	client, err := sirius.NewClient(http.DefaultClient, siriusURL)
	if err != nil {
		logger.Fatal(err)
	}

	s := &http.Server{
		Addr:    ":" + port,
		Handler: server.New(logger, client, tmpls, prefix, siriusPublicURL, firmHubURL, webDir),
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Print("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Print("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(tc); err != nil {
		logger.Print(err)
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
