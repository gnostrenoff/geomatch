package api

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/gnostenoff/geomatch/configs"
	"github.com/gnostenoff/geomatch/internal/api/handlers"
	"github.com/gnostenoff/geomatch/internal/geomatch/haversine"
	"github.com/gnostenoff/geomatch/internal/pkg/datasource/csvloader"
)

func Run() {
	r := gin.New()

	// load config
	configs.Init(os.Getenv("ENV"))

	geoMatchHandler := initDependencies()

	r.POST("/events", geoMatchHandler.Handle)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	slog.Info("Starting server", slog.Int("port", 8080))
	_ = server.ListenAndServe()
}

func initDependencies() (geoMatchHandler handlers.GeoMatchHandler) {
	csvLoader := csvloader.NewEventCSVLoader(configs.Config.GetString("CSV_FILE_PATH"))
	geoMatcher := haversine.NewMatcher(csvLoader)
	geoMatchHandler = handlers.GeoMatchHandler{GeoMatcher: geoMatcher}

	// init datasource by loading events from CSV file
	err := csvLoader.Load()
	if err != nil {
		slog.Error("failed to load events from CSV file", slog.String("filePath", configs.Config.GetString("CSV_FILE_PATH")), slog.String("error", err.Error()))
	}

	return
}
