package handlers

import (
	"net/http"

	"github.com/geek-remote-id/tspay-subhub/utils"
)

// GenerateIncomingHandler godoc
// @Summary      Incoming Handler
// @Description  Handles incoming requests for Tspay Subhub
// @Tags         incoming
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.Response
// @Router       /incoming [get]
func GenerateIncomingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, utils.Response{
			Status:  "success",
			Message: "Incoming Handler",
		})
	}
}
