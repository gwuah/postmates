package handler

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

type closestCourierResponse struct {
	Couriers []string `json:"couriers"`
}

func (h *Handler) handleGetClosestCouriers(message []byte, ws *ws.WSConnection) {
	var data shared.GetClosestCouriersRequest
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("failed to parse message", err)
		return
	}

	closestCouriers := h.Services.GetClosestCouriers(data.Origin, 2)

	stringifiedResponse, err := json.Marshal(closestCourierResponse{
		Couriers: closestCouriers,
	})

	if err != nil {
		log.Println("failed to marshal message", err)
		return
	}

	ws.SendMessage([]byte(stringifiedResponse))

}
