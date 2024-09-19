package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theweird-kid/message-queue/internals/queue"
	"github.com/theweird-kid/message-queue/utils"
)

func (h *Handler) GetTopics(w http.ResponseWriter, r *http.Request) {
	topics := h.e.GetTopics()
	utils.RespondWithJSON(w, 200, topics)
}

type CreateTopicRequest struct {
	Name       string `json:"name"`
	BufferSize int    `json:"buffer_size"`
}

func (h *Handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	// Extract Payload from request
	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	h.e.CreateTopic(req.Name, req.BufferSize)
	utils.RespondWithJSON(w, 200, fmt.Sprintf("New topic %s created", req.Name))
}

// Handler to publish message to topics
type PubMessage struct {
	Topic string `json:"topic"`
	Msg   string `json:"message"`
}

func (h *Handler) PublishMessage(w http.ResponseWriter, r *http.Request) {
	var req PubMessage
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	defer r.Body.Close()

	//log.Println(req)
	message := queue.Message{Content: req.Msg}
	//log.Println("topic:", req.Topic)

	err := h.e.Publish(req.Topic, message)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, 200, "Message Published Successfully")
}

// Handler to Get message from topic
