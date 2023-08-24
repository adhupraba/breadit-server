package controllers

import (
	"net/http"

	"github.com/adhupraba/breadit-server/internal/types"
	"github.com/adhupraba/breadit-server/utils"
)

type HealthController struct{}

func (hc *HealthController) Heartbeat(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJson(w, http.StatusOK, types.Json{"message": "success"})
}
