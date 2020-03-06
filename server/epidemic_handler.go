package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NeoHuang/bit-hedge/api"
)

const (
	InternalErrorResponse = `{"error":"internal error"}`
)

type EpidemicHandler struct {
	api api.Api
}

func NewEpidemicHandler(rkiApi api.Api) *EpidemicHandler {
	return &EpidemicHandler{
		api: rkiApi,
	}
}

func (handler *EpidemicHandler) ServeHTTP(writer http.ResponseWriter, httpRequest *http.Request) {
	epidemicMap, err := handler.api.GetCurrent()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(InternalErrorResponse))
		log.Printf("Failed to get current data. err:%s", err)
		return
	}

	bytes, err := json.Marshal(epidemicMap)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(InternalErrorResponse))
		log.Printf("Failed to marshal current data:%v err:%s", epidemicMap, err)
		return
	}

	writer.Write(bytes)
}
