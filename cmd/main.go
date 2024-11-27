package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jake-young-dev/bed/cron"
	dotenv "github.com/joho/godotenv"
)

const (
	//rcon vars
	RCON_PASSWORD_ENV = "RCON_PASSWORD"
	MC_CONTAINER_ENV  = "RCON_MC_CONTAINER"

	//minio vars
	MINIO_URL_ENV        = "MINIO_URL"
	MINIO_BUCKET_ENV     = "MINIO_BUCKET"
	MINIO_ACCESS_KEY_ENV = "MINIO_KEY"
	MINIO_ACCESS_ID_ENV  = "MINIO_ID"
)

// entry point
func main() {
	dotenv.Load()
	err := validateEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting backup cron")
	backup := cron.NewCronHandler()
	backup.TakeBackup() //this func doesn't need to be exposed after testing
	// backup.Run()

	// //wait for interupt
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	// <-quit

	// log.Println("stopping backup cron")
	// backup.Stop()
}

// validates the existence of necessary environment variables required to run the cron
func validateEnvironment() error {
	var missingVars []string
	//validate rcon variables
	if _, exists := os.LookupEnv(RCON_PASSWORD_ENV); !exists {
		missingVars = append(missingVars, RCON_PASSWORD_ENV)
	}
	if _, exists := os.LookupEnv(MC_CONTAINER_ENV); !exists {
		missingVars = append(missingVars, MC_CONTAINER_ENV)
	}
	//validate minio variables
	if _, exists := os.LookupEnv(MINIO_BUCKET_ENV); !exists {
		missingVars = append(missingVars, MINIO_BUCKET_ENV)
	}
	if _, exists := os.LookupEnv(MINIO_ACCESS_KEY_ENV); !exists {
		missingVars = append(missingVars, MINIO_ACCESS_KEY_ENV)
	}
	if _, exists := os.LookupEnv(MINIO_URL_ENV); !exists {
		missingVars = append(missingVars, MINIO_URL_ENV)
	}
	if _, exists := os.LookupEnv(MINIO_ACCESS_ID_ENV); !exists {
		missingVars = append(missingVars, MINIO_ACCESS_ID_ENV)
	}

	//if any vars are missing list them here
	if len(missingVars) > 0 {
		missingValues := strings.Join(missingVars, ", ")
		return fmt.Errorf("missing environment variables: %s", missingValues)
	}

	return nil
}
