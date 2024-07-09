package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/lib/request"
	resp "time-tracker/internal/lib/response"
	"time-tracker/internal/models"
	taskService "time-tracker/internal/service/task"
	userService "time-tracker/internal/service/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	uuidlib "github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, passportSerie, passportNumber int) (*models.User, error)
	GetUsers(ctx context.Context, page int, filter string) ([]models.User, error)
	UpdateUserInfo(ctx context.Context, userInfo *models.User) (*models.User, error)
	RemoveUserByUUID(ctx context.Context, uuid string) error
}

type TaskService interface {
	GetTasksInRange(ctx context.Context, userUUID, startDate, endDate string) ([]models.Task, error)
}

type Handler struct {
	userService UserService
	taskService TaskService
	log         *slog.Logger
}

func New(userService UserService, taskService TaskService, log *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		taskService: taskService,
		log:         log,
	}
}

func (h *Handler) Register() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/", h.getUsers)
		r.Patch("/{uuid}", h.updateUser)
		r.Delete("/{uuid}", h.deleteUser)
		r.Get("/{user_id}/worklogs", h.getTasksInRange)
	}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.createUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	var credentials *request.CreateUser
	if err := render.DecodeJSON(r.Body, &credentials); err != nil {
		log.Error("failed to decode request body", sl.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Err("Internal Error"))
		return
	}

	if credentials.PassportNumber == "" {
		log.Debug("request body is empty")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err("Request body is empty"))
		return
	}

	log.Debug("creating new user", slog.String("passport_number", credentials.PassportNumber))

	passport := strings.Split(credentials.PassportNumber, " ")

	passportSerie, err := strconv.Atoi(passport[0])
	if err != nil {
		log.Error("failed to get passport serie", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err("Invalid passport serie"))
		return
	}

	passportNumber, err := strconv.Atoi(passport[1])
	if err != nil {
		log.Error("failed to get passport number", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err("Invalid passport number"))
		return
	}

	user, err := h.userService.CreateUser(r.Context(), passportSerie, passportNumber)
	if err != nil {
		if errors.Is(err, userService.ErrExists) {
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, resp.Err("User already exists"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("Internal error"))
			return
		}
	}

	log.Debug("user created successfully")

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.getUsers"

	log := h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	// Getting `page` param & validation
	var page int
	p := r.URL.Query().Get("page")
	if p == "" {
		page = 1
		log.Debug(`set "page" value to default`, slog.Int("page", page))
	} else {
		parsedPage, err := strconv.Atoi(p)
		if err != nil || parsedPage < 1 {
			log.Error(`error while parsing "page" param`, sl.Error(err))
			page = 1
			log.Debug(`set "page" value to default`, slog.Int("page", page))
		}
		page = parsedPage
		log.Debug(`validate "page" value`, slog.Int("page", page))
	}

	// Getting `filter` param
	filter := r.URL.Query().Get("filter")

	log.Debug("getting all users", slog.Int("page", page), slog.String("filter", filter))

	users, err := h.userService.GetUsers(r.Context(), page, filter)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Err("Internal error"))
		return
	}

	log.Debug("got all users successfully")

	render.JSON(w, r, users)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.updateUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	uuid := chi.URLParam(r, "uuid")

	_, err := uuidlib.Parse(uuid)
	if err != nil {
		log.Error("invalid userUUID", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err(`invalid user uuid format`))
	}

	log.Debug("patching user info", slog.String("user_uuid", uuid))

	var userInfo models.User
	if err := render.DecodeJSON(r.Body, &userInfo); err != nil {
		log.Error("failed to decode request body", sl.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Err("Internal Error"))
		return
	}

	// Set uuid from user
	userInfo.ID = uuid

	user, err := h.userService.UpdateUserInfo(r.Context(), &userInfo)
	if err != nil {
		if errors.Is(err, userService.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Err("User not found"))
			return
		} else if errors.Is(err, userService.ErrEmptyBody) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Err("Request body is empty"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("Internal error"))
			return
		}
	}

	log.Debug("user patched succesfully", slog.String("user_uuid", uuid))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.deleteUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	uuid := chi.URLParam(r, "uuid")

	_, err := uuidlib.Parse(uuid)
	if err != nil {
		log.Error("invalid userUUID", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err(`invalid user uuid format`))
	}

	log.Debug("removing user", slog.String("user_uuid", uuid))

	err = h.userService.RemoveUserByUUID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, userService.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, resp.Err("User not found"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Err("Internal error"))
			return
		}
	}

	log.Debug("user removed succesfully", slog.String("user_uuid", uuid))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp.Ok("User removed successfully"))
}

func (h *Handler) getTasksInRange(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.getTaskInRange"

	log := h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	userUUID := chi.URLParam(r, "user_id")
	if userUUID == "" {
		log.Error("missing user_id parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err(`'user_id' parameter is required`))
		return
	}

	startDate := r.URL.Query().Get("start_date")
	if startDate == "" {
		log.Error("missing start_date parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err(`'start_date' parameter is required`))
		return
	}

	endDate := r.URL.Query().Get("end_date")
	if endDate == "" {
		log.Error("missing end_date parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp.Err(`'end_date' parameter is required`))
		return
	}

	log.Debug("getting tasks in range", slog.String("user_id", userUUID), slog.String("start_date", startDate), slog.String("end_date", endDate))

	tasks, err := h.taskService.GetTasksInRange(r.Context(), userUUID, startDate, endDate)
	if err != nil {
		if errors.Is(err, taskService.ErrInvalidDateRange) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Err("Invalid date range"))
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, resp.Err("Internal error"))
		return
	}

	log.Debug("got tasks in range successfully")

	render.JSON(w, r, tasks)
}
