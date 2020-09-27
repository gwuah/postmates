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

func (h *Handler) handleElectronLocationUpdate(message []byte, ws *ws.WSConnection) {
	var data shared.UserLocationUpdate
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("Failed to parse message", err)
		return
	}

	electron, err := h.Services.IndexElectronLocation(data)
	if err != nil {
		log.Println("Failed to marshal message", err)
		return
	}

	_, err = json.Marshal(electron)

	if err != nil {
		log.Println("Failed to marshal message", err)
		return
	}

	// ws.SendMessage([]byte(stringifiedResponse))
}

func (h *Handler) handleGetClosestElectrons(message []byte, ws *ws.WSConnection) {
	var data shared.GetClosestElectronsRequest
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("Failed to parse message", err)
		return
	}

	closestElectrons := h.Services.GetClosestElectrons(data.Origin, 2)

	stringifiedResponse, err := json.Marshal(closestElectronsResponse{
		Electrons: closestElectrons,
	})

	if err != nil {
		log.Println("Failed to marshal message", err)
		return
	}

	ws.SendMessage([]byte(stringifiedResponse))

}
