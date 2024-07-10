package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/lib/request"
	"time-tracker/internal/lib/response"
	"time-tracker/internal/models"
	service "time-tracker/internal/service/user"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	uuidlib "github.com/google/uuid"
)

type Service interface {
	CreateUser(ctx context.Context, passportSerie, passportNumber int) (*models.User, error)
	GetUsers(ctx context.Context, page int, filter string) ([]models.User, error)
	UpdateUserInfo(ctx context.Context, userInfo *models.User) (*models.User, error)
	RemoveUserByUUID(ctx context.Context, uuid string) error
}

type TaskService interface {
	GetTasksInRange(ctx context.Context, userUUID, startDate, endDate string) ([]models.Task, error)
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
		r.Post("/", h.createUser)
		r.Get("/", h.getUsers)
		r.Patch("/{uuid}", h.updateUser)
		r.Delete("/{uuid}", h.deleteUser)
	}
}

// @Summary Создание нового пользователя
// @Description Создает нового пользователя по паспортным данным
// @Tags users
// @Accept json
// @Produce json
// @Param CreateUser body request.CreateUser true "Данные для создания пользователя"
// @Success 201 {object} models.User "Пользователь создан успешно"
// @Failure 400 {object} response.Response "Некорректные данные запроса"
// @Failure 409 {object} response.Response "Пользователь уже существует"
// @Failure 500 {object} response.Response "Внутренняя ошибка сервера"
// @Router /users [post]
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
		render.JSON(w, r, response.Err("Internal Error"))
		return
	}

	if credentials.PassportNumber == "" {
		log.Debug("request body is empty")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err("Request body is empty"))
		return
	}

	log.Debug("creating new user", slog.String("passport_number", credentials.PassportNumber))

	passport := strings.Split(credentials.PassportNumber, " ")

	passportSerie, err := strconv.Atoi(passport[0])
	if err != nil {
		log.Error("failed to get passport serie", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err("Invalid passport serie"))
		return
	}

	passportNumber, err := strconv.Atoi(passport[1])
	if err != nil {
		log.Error("failed to get passport number", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err("Invalid passport number"))
		return
	}

	user, err := h.service.CreateUser(r.Context(), passportSerie, passportNumber)
	if err != nil {
		if errors.Is(err, service.ErrExists) {
			render.Status(r, http.StatusConflict)
			render.JSON(w, r, response.Err("User already exists"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Err("Internal error"))
			return
		}
	}

	log.Debug("user created successfully")

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

// @Summary Получить пользователей
// @Description Получить список пользователей с возможностью фильтрации и пагинации
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param filter query string false "Строка фильтра"
// @Success 200 {array} models.User
// @Failure 500 {object} response.Response "Внутренняя ошибка"
// @Router /users [get]
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

	users, err := h.service.GetUsers(r.Context(), page, filter)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Err("Internal error"))
		return
	}

	log.Debug("got all users successfully")

	render.JSON(w, r, users)
}

// @Summary Обновить пользователя
// @Description Обновить информацию о пользователе по UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "UUID пользователя"
// @Param user body models.User true "Информация о пользователе"
// @Success 200 {object} models.User
// @Failure 400 {object} response.Response "Неверный формат UUID или пустое тело запроса"
// @Failure 404 {object} response.Response "Пользователь не найден"
// @Failure 500 {object} response.Response "Внутренняя ошибка"
// @Router /users/{uuid} [put]
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
		render.JSON(w, r, response.Err(`invalid user uuid format`))
		return
	}

	log.Debug("patching user info", slog.String("user_uuid", uuid))

	var userInfo models.User
	if err := render.DecodeJSON(r.Body, &userInfo); err != nil {
		log.Error("failed to decode request body", sl.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Err("Internal Error"))
		return
	}

	// Set uuid from user
	userInfo.ID = uuid

	user, err := h.service.UpdateUserInfo(r.Context(), &userInfo)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Err("User not found"))
			return
		} else if errors.Is(err, service.ErrEmptyBody) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Err("Request body is empty"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Err("Internal error"))
			return
		}
	}

	log.Debug("user patched succesfully", slog.String("user_uuid", uuid))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

// @Summary Удалить пользователя
// @Description Удалить пользователя по UUID
// @Tags users
// @Accept json
// @Produce json
// @Param uuid path string true "UUID пользователя"
// @Success 200 {object} response.Response "Пользователь успешно удалён"
// @Failure 400 {object} response.Response "Неверный формат UUID"
// @Failure 404 {object} response.Response "Пользователь не найден"
// @Failure 500 {object} response.Response "Внутренняя ошибка"
// @Router /users/{uuid} [delete]
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
		render.JSON(w, r, response.Err(`invalid user uuid format`))
		return
	}

	log.Debug("removing user", slog.String("user_uuid", uuid))

	err = h.service.RemoveUserByUUID(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Err("User not found"))
			return
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Err("Internal error"))
			return
		}
	}

	log.Debug("user removed succesfully", slog.String("user_uuid", uuid))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.Ok("User removed successfully"))
}
