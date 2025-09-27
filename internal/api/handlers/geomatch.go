package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gnostenoff/geomatch/internal/geomatch"
)

type GeoMatchHandler struct {
	GeoMatcher geomatch.GeoMatcher
}

type GeoMatchRequest struct {
	PointsOfInterest []geomatch.PointOfInterest `json:"points_of_interest"`
}

func (h GeoMatchHandler) Handle(c *gin.Context) {
	req, err := readBody(c)
	if err != nil || !isValid(req) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	matchResults, err := h.GeoMatcher.Match(req.PointsOfInterest)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, matchResults)
}

func readBody(c *gin.Context) (*GeoMatchRequest, error) {
	req := GeoMatchRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		slog.Warn("Unable to read Body", slog.String("msg_raw", err.Error()))
		return nil, err
	}
	return &req, nil
}

func isValid(req *GeoMatchRequest) bool {
	return req != nil && len(req.PointsOfInterest) > 0
}
