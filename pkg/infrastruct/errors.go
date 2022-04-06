package infrastruct

import "github.com/gin-gonic/gin"

var (
	ErrorInternalServerError = gin.H{"error": "internal server error"}
	ErrorBadRequest          = gin.H{"error": "bad query input"}
	ErrorForbidden           = gin.H{"error": "forbidden"}
)
