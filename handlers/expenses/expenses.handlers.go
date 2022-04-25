package expenses_handlers

import (
	"expense_tracker/models"
	"expense_tracker/pkg/infrastruct"
	services "expense_tracker/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ExpensesHandlers struct {
	s services.ExpensesService
}

func NewExpensesHandlers(s services.ExpensesService) *ExpensesHandlers {
	return &ExpensesHandlers{
		s: s,
	}
}

type CreateCategoryDTO struct {
	Title string `json:"title" binding:"required"`
}

func (h *ExpensesHandlers) CreateCategory(ctx *gin.Context) {

	userID := ctx.GetString("user")

	var dto CreateCategoryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	category, err := h.s.CreateCategory(dto.Title, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (h *ExpensesHandlers) DeleteCategory(ctx *gin.Context) {

	userID := ctx.GetString("user")
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	deleted, err := h.s.DeleteCategory(id, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}

func (h *ExpensesHandlers) GetCategoryList(ctx *gin.Context) {
	userID := ctx.GetString("user")
	query := ctx.Query("query")

	list, err := h.s.GetCategoryList(query, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *ExpensesHandlers) DeleteExpense(ctx *gin.Context) {
	userID := ctx.GetString("user")
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	deleted, err := h.s.DeleteExpense(id, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, deleted)
}

func (h *ExpensesHandlers) CreateExpense(ctx *gin.Context) {
	userID := ctx.GetString("user")
	var dto models.ExpenseDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	if dto.Amount == 0 {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	if !(dto.Currency == models.USD || dto.Currency == models.EUR || dto.Currency == models.RUB || dto.Currency == models.TBH) {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	category, err := h.s.CreateExpense(dto, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, category)
}

type DayTotalDTO struct {
	Date     time.Time       `json:"date"`
	Currency models.Currency `json:"currency"`
}

func (h *ExpensesHandlers) GetDayTotalExpenses(ctx *gin.Context) {
	userID := ctx.GetString("user")
	var dto DayTotalDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	if !(dto.Currency == models.USD || dto.Currency == models.EUR || dto.Currency == models.RUB || dto.Currency == models.TBH) {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	result, err := h.s.GetDayTotalExpenses(dto.Date, dto.Currency, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, *result)
}

func (h *ExpensesHandlers) GetMonthExpenses(ctx *gin.Context) {
	userID := ctx.GetString("user")
	var dto models.MonthDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	if !(dto.Currency == models.USD || dto.Currency == models.EUR || dto.Currency == models.RUB || dto.Currency == models.TBH) {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	result, err := h.s.GetMonthExpenses(&dto, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, *result)
}

func (h *ExpensesHandlers) GetMonthExpensesByCategory(ctx *gin.Context) {
	userID := ctx.GetString("user")
	var dto DayTotalDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	if !(dto.Currency == models.USD || dto.Currency == models.EUR || dto.Currency == models.RUB || dto.Currency == models.TBH) {
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	result, err := h.s.GetMonthExpensesByCategory(dto.Date, dto.Currency, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, *result)
}

func (h *ExpensesHandlers) GetCurrencyRate(ctx *gin.Context) {

	result, err := h.s.GetCurrencyRate()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, *result)
}

func (h *ExpensesHandlers) SyncData(ctx *gin.Context) {

	userID := ctx.GetString("user")

	var dto models.SyncDTO

	if err := ctx.ShouldBindJSON(&dto); err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorBadRequest)
		return
	}

	result, err := h.s.SyncDatabase(dto, userID)

	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, infrastruct.ErrorInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, result)

}
