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
	r.engine.POST("/login", r.handlers.Login)

	r.engine.POST("/category", r.handlers.AuthGuard(), r.handlers.CreateCategory)
	r.engine.GET("/categories", r.handlers.AuthGuard(), r.handlers.GetCategoryList)
	r.engine.DELETE("/category/:id", r.handlers.AuthGuard(), r.handlers.DeleteCategory)
	r.engine.POST("/expense", r.handlers.AuthGuard(), r.handlers.CreateExpense)
	r.engine.DELETE("/expense/:id", r.handlers.AuthGuard(), r.handlers.DeleteExpense)
	r.engine.POST("/expenses/day", r.handlers.AuthGuard(), r.handlers.GetDayTotalExpenses)
	r.engine.POST("/expenses/month", r.handlers.AuthGuard(), r.handlers.GetMonthExpenses)
	r.engine.POST("/expenses/category", r.handlers.AuthGuard(), r.handlers.GetMonthExpensesByCategory)
	r.engine.GET("/rates", r.handlers.AuthGuard(), r.handlers.GetCurrencyRate)

	r.engine.POST("/sync", r.handlers.AuthGuard(), r.handlers.SyncData)
}
