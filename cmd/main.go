package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/jake-young-dev/bed/cron"
	dotenv "github.com/joho/godotenv"
)

const (
	//rcon variables to validate
	RCON_PASSWORD_ENV = "RCON_PASSWORD"
	MC_CONTAINER_ENV  = "RCON_MC_CONTAINER"
	RCON_PORT_ENV     = "RCON_PORT"

	//minio variables to validate
	MINIO_URL_ENV        = "MINIO_URL"
	MINIO_BUCKET_ENV     = "MINIO_BUCKET"
	MINIO_ACCESS_KEY_ENV = "MINIO_KEY"
	MINIO_ACCESS_ID_ENV  = "MINIO_ID"

	//server action variables to validate
	SERVER_RESTART_ENV = "SERVER_RESTART"
	SERVER_RESTART_YES = "yes"
	SERVER_RESTART_NO  = "no"
)

// simple array to allow for a much smaller validation function
var ENV_VARS_TO_VALIDATE = [8]string{
	RCON_PASSWORD_ENV,
	MC_CONTAINER_ENV,
	RCON_PORT_ENV,
	MINIO_URL_ENV,
	MINIO_BUCKET_ENV,
	MINIO_ACCESS_KEY_ENV,
	MINIO_ACCESS_ID_ENV,
	SERVER_RESTART_ENV,
}

// entry point
func main() {
	dotenv.Load()
	//ensure we have our environment variables
	err := validateEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting backup cron")
	backup := cron.NewCronHandler()
	backup.Run()

	//safe stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("stopping backup cron")
	backup.Stop()
}

// validates that we have the environment variables we need to work with rcon and minio
func validateEnvironment() error {
	var missingVars []string

	//loop through environment variables and validate value if needed
	for _, e := range ENV_VARS_TO_VALIDATE {
		if _, exists := os.LookupEnv(e); !exists {
			missingVars = append(missingVars, e)
		} else {
			if e == SERVER_RESTART_ENV {
				//validate acceptable values
				if os.Getenv(SERVER_RESTART_ENV) != SERVER_RESTART_YES && os.Getenv(SERVER_RESTART_ENV) != SERVER_RESTART_NO {
					missingVars = append(missingVars, fmt.Sprintf("%s is set to an invalid value, use '%s' or '%s'", SERVER_RESTART_ENV, SERVER_RESTART_YES, SERVER_RESTART_NO))
				}
			}
		}
	}

	if len(missingVars) > 0 {
		missingValues := strings.Join(missingVars, ", ")
		return fmt.Errorf("missing environment variables: %s", missingValues)
	}

	return nil
}
