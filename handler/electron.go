package handler

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

type closestElectronsResponse struct {
	Electrons []string `json:"electrons"`
}

func (h *Handler) handleGetClosestElectrons(message []byte, ws *ws.WSConnection) {
	var data shared.GetClosestElectronsRequest
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("failed to parse message", err)
		return
	}

	closestElectrons := h.Services.GetClosestElectrons(data.Origin, 2)

	stringifiedResponse, err := json.Marshal(closestElectronsResponse{
		Electrons: closestElectrons,
	})

	if err != nil {
		log.Println("failed to marshal message", err)
		return
	}

	ws.SendMessage([]byte(stringifiedResponse))

}
