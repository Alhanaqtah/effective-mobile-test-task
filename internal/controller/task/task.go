package task

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/lib/response"
	"time-tracker/internal/models"
	service "time-tracker/internal/service/task"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	uuidlib "github.com/google/uuid"
)

type Service interface {
	GetTasksInRange(ctx context.Context, userUUID, startDate, endDate string) ([]models.Task, error)
	StartTask(ctx context.Context, uuid string) (*models.Task, error)
	FinishTask(ctx context.Context, uuid string) (*models.Task, error)
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
		r.Get("/{user_id}/worklogs", h.getTasksInRange)
		r.Post("/{task_id}/start", h.startTask)
		r.Post("/{task_id}/finish", h.finishTask)
	}
}

// @Summary Получить задачи в диапазоне дат
// @Description Возвращает задачи пользователя в заданном диапазоне дат
// @Tags tasks
// @Accept json
// @Produce json
// @Param user_uuid path string true "UUID пользователя"
// @Param start_date query string true "Дата начала в формате YYYY-MM-DD"
// @Param end_date query string true "Дата окончания в формате YYYY-MM-DD"
// @Success 200 {array} models.Task "Список задач"
// @Failure 400 {object} response.Response "Некорректный запрос"
// @Failure 500 {object} response.Response "Внутренняя ошибка"
// @Router /users/{user_uuid}/tasks [get]
func (h *Handler) getTasksInRange(w http.ResponseWriter, r *http.Request) {
	const op = "controller.user.getTaskInRange"

	log := h.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	userUUID := chi.URLParam(r, "user_uuid")
	if userUUID == "" {
		log.Error("missing user_id parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err(`'user_id' parameter is required`))
		return
	}

	startDate := r.URL.Query().Get("start_date")
	if startDate == "" {
		log.Error("missing start_date parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err(`'start_date' parameter is required`))
		return
	}

	endDate := r.URL.Query().Get("end_date")
	if endDate == "" {
		log.Error("missing end_date parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err(`'end_date' parameter is required`))
		return
	}

	log.Debug("getting tasks in range", slog.String("user_id", userUUID), slog.String("start_date", startDate), slog.String("end_date", endDate))

	tasks, err := h.service.GetTasksInRange(r.Context(), userUUID, startDate, endDate)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateRange) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Err("Invalid date range"))
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Err("Internal error"))
		return
	}

	log.Debug("got tasks in range successfully")

	render.JSON(w, r, tasks)
}

// @Summary Запуск задачи
// @Description Запускает задачу по ее UUID
// @Tags tasks
// @Accept json
// @Produce json
// @Param task_id path string true "UUID задачи"
// @Success 200 {object} models.Task "Запущенная задача"
// @Failure 400 {object} response.Response "Неверный формат UUID или пустое тело запроса"
// @Failure 404 {object} response.Response "Задача не найдена"
// @Failure 500 {object} response.Response "Внутренняя ошибка сервера"
// @Router /tasks/{task_uuid} [post]
func (h *Handler) startTask(w http.ResponseWriter, r *http.Request) {
	const op = "controller.task.startTask"

	log := h.log.With(slog.String("op", op))

	uuid := chi.URLParam(r, "task_id")

	_, err := uuidlib.Parse(uuid)
	if err != nil {
		log.Error("invalid userUUID", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err(`invalid user uuid format`))
		return
	}

	log.Debug("starting task", slog.String("uuid", uuid))

	task, err := h.service.StartTask(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Err("Task not found"))
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Err("Internal error"))
		return
	}

	log.Debug("task started successfully")
	render.JSON(w, r, task)
}

// @Summary Завершение задачи
// @Description Отметить задачу как завершенную
// @Tags tasks
// @Accept json
// @Produce json
// @Param task_id path string true "ID задачи"
// @Success 200 {object} models.Task
// @Failure 400 {object} response.Response "Неверный формат UUID"
// @Failure 404 {object} response.Response "Задача не найдена"
// @Failure 500 {object} response.Response "Внутренняя ошибка"
// @Router /tasks/{task_id}/finish [post]
func (h *Handler) finishTask(w http.ResponseWriter, r *http.Request) {
	const op = "controller.task.finishTask"

	log := h.log.With(slog.String("op", op))

	uuid := chi.URLParam(r, "task_id")

	_, err := uuidlib.Parse(uuid)
	if err != nil {
		log.Error("invalid userUUID", sl.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err(`invalid user uuid format`))
		return
	}

	log.Debug("finishing task", slog.String("uuid", uuid))

	task, err := h.service.FinishTask(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Err("Task not found"))
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Err("Internal error"))
		return
	}

	log.Debug("task finished successfully")
	render.JSON(w, r, task)
}
