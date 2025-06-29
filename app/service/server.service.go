// services/server
package service

import (
	"log"
	"xi/app/lib"

	"github.com/gin-gonic/gin"
)

// InitServer start the web server
func InitServer(app *gin.Engine) error {
	appPort := lib.Env("APP_PORT", "5000")

	log.Printf("\a\033[1;94mServer started \033[0;34m'http://localhost:%s'%s\033[0m\n", appPort,
		func() string {
			if url := lib.Env("APP_URL"); url != "" {
				return ", 'http://" + url + "'"
			}
			return ""
		}())

	// Start Web-Server
	if err := app.Run(":" + appPort); err != nil {
		log.Fatalf("❌ Err starting server: %v", err)
		return err
	}

	return nil
}
