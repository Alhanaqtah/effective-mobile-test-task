package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"time-tracker/internal/lib/logger/sl"
	resp "time-tracker/internal/lib/response"
	service "time-tracker/internal/service/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Service interface {
	// GetUserWithPagination(offset int) (*models.User, error)
	RemoveUserByUUID(uuid string) error
}

type Handler struct {
	service Service
	log     *slog.Logger
}

func New(service Service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) Register() func(r chi.Router) {
	return func(r chi.Router) {
		// r.Post("/", h.createUser)
		// r.Get("/{page}", h.getUsers)
		r.Patch("/{uuid}", h.updateUser)
		r.Delete("/{uuid}", h.deleteUser)
	}
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.deleteUser"

	log := h.log.With("op", op)

	uuid := chi.URLParam(r, "uuid")

	log.Debug("removing user", slog.String("user_uuid", uuid))

	err := h.service.RemoveUserByUUID(uuid)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Error("failed to remove user", sl.Error(err))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Err("User not found"))
			return
		} else {
			log.Error("failed to remove user", sl.Error(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("Internal error"))
			return
		}
	}

	render.Status(r, http.StatusOK)
}
