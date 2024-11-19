package link

import (
	"arch/ikeppu/github.com/configs"
	"arch/ikeppu/github.com/pkg/event"
	"arch/ikeppu/github.com/pkg/middleware"
	"arch/ikeppu/github.com/pkg/req"
	"arch/ikeppu/github.com/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}
type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	linkHandler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}

	router.HandleFunc("POST /link", linkHandler.Create)
	router.Handle("PATCH /link/{id}", middleware.Auth(linkHandler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", linkHandler.Delete)
	router.HandleFunc("GET /{alias}", linkHandler.Get)

}

func (handler *LinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[LinkCreateRequest](&w, r)
	if err != nil {
		return
	}

	link := NewLink(body.Url)
	createdLink, err := handler.LinkRepository.Create(link)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.Json(w, createdLink, http.StatusCreated)
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(middleware.ContextEmailKey).(string)

		if ok {
		}

		body, err := req.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return
		}

		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{
				ID: uint(id),
			},
			Url:  body.Url,
			Hash: body.Hash,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		response.Json(w, link, 201)
	}
}

func (handler *LinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = handler.LinkRepository.GetById(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = handler.LinkRepository.Delete(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response.Json(w, nil, http.StatusOK)
}

func (handler *LinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("alias")

	link, err := handler.LinkRepository.GetByHash(hash)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	go handler.EventBus.Publish(event.Event{
		Type: event.EventLinkVisited,
		Data: link.ID,
	})

	http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
}
