package notify

import "github.com/gin-gonic/gin"

func GinRecovery() gin.HandlerFunc {
	return recoverError()
}

func recoverError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e, _ := err.(error)

				ErrorRequest(e, c.Request)

				panic(err)
			}
		}()

		c.Next()
	}
}
