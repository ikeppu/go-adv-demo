package stat

import (
	"arch/ikeppu/github.com/configs"
	"fmt"
	"net/http"
	"time"
)

const (
	FilterByDay   = "day"
	FilterByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}
type StatHandler struct {
	LinkRepository *StatRepository
	Config         *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	statHandler := &StatHandler{
		LinkRepository: deps.StatRepository,
		Config:         deps.Config,
	}

	router.HandleFunc("GET /stat", statHandler.GetStat())

}

func (h *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to param", http.StatusBadRequest)
			return
		}

		by := r.URL.Query().Get("by")

		if by != FilterByDay && by != FilterByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}

		fmt.Println(from, to, by)
	}
}
