package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	periodictaskdomain "example.com/taskservice/internal/domain/periodictask"
	periodictaskusecase "example.com/taskservice/internal/usecase/periodictask"
)

type PeriodicTaskHandler struct {
	usecase periodictaskusecase.Usecase
}

func NewPeriodicTaskHandler(usecase periodictaskusecase.Usecase) *PeriodicTaskHandler {
	return &PeriodicTaskHandler{usecase: usecase}
}

func (h *PeriodicTaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req periodicTaskCreateDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD"))
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD"))
			return
		}
		endDate = &parsed
	}

	created, err := h.usecase.Create(r.Context(), periodictaskusecase.CreateInput{
		Title:            req.Title,
		Description:      req.Description,
		RecurrenceType:   req.RecurrenceType,
		RecurrenceParams: req.RecurrenceParams,
		StartDate:        startDate,
		EndDate:          endDate,
	})
	if err != nil {
		writePeriodicTaskUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newPeriodicTaskResponseDTO(created))
}

func (h *PeriodicTaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	pt, err := h.usecase.GetByID(r.Context(), id)
	if err != nil {
		writePeriodicTaskUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newPeriodicTaskResponseDTO(pt))
}

func (h *PeriodicTaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req periodicTaskUpdateDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD"))
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD"))
			return
		}
		endDate = &parsed
	}

	updated, err := h.usecase.Update(r.Context(), id, periodictaskusecase.UpdateInput{
		Title:            req.Title,
		Description:      req.Description,
		RecurrenceType:   req.RecurrenceType,
		RecurrenceParams: req.RecurrenceParams,
		StartDate:        startDate,
		EndDate:          endDate,
		IsActive:         req.IsActive,
	})
	if err != nil {
		writePeriodicTaskUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newPeriodicTaskResponseDTO(updated))
}

func (h *PeriodicTaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.usecase.Delete(r.Context(), id); err != nil {
		writePeriodicTaskUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PeriodicTaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.usecase.List(r.Context())
	if err != nil {
		writePeriodicTaskUsecaseError(w, err)
		return
	}

	response := make([]periodicTaskResponseDTO, 0, len(tasks))
	for i := range tasks {
		response = append(response, newPeriodicTaskResponseDTO(&tasks[i]))
	}

	writeJSON(w, http.StatusOK, response)
}

func writePeriodicTaskUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, periodictaskdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, periodictaskusecase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		writeError(w, http.StatusInternalServerError, err)
	}
}
