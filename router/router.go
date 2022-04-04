package router

import (
	"github.com/gin-gonic/gin"

	handlers "expense_tracker/handlers"
)

type Router struct {
	engine   *gin.Engine
	handlers *handlers.Handlers
}

func NewRouter(engine *gin.Engine, handlers *handlers.Handlers) Router {

	return Router{engine, handlers}
}

func (r *Router) SetupRouter() {

	r.engine.POST("/category", r.handlers.CreateCategory)
	r.engine.GET("/categories", r.handlers.GetCategoryList)
	r.engine.DELETE("/category/:id", r.handlers.DeleteCategory)
	r.engine.POST("/expense", r.handlers.CreateExpense)
	r.engine.DELETE("/expense/:id", r.handlers.DeleteExpense)
	r.engine.POST("/expenses/day", r.handlers.GetDayTotalExpenses)
	r.engine.POST("/expenses/month", r.handlers.GetMonthExpenses)
	r.engine.POST("/expenses/category", r.handlers.GetMonthExpensesByCategory)
}
