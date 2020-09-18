package handler

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

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

	stringifiedResponse, err := json.Marshal(electron)

	if err != nil {
		log.Println("Failed to marshal message", err)
		return
	}

	ws.SendMessage([]byte(stringifiedResponse))
}
