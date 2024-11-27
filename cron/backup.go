package cron

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jake-young-dev/bed/file"
	"github.com/jake-young-dev/bed/minecraft"
	qz "github.com/jake-young-dev/quick-zip"
	cronjob "github.com/robfig/cron"
)

const (
	//used for easier log searching and to divide backup iterations in logs
	logDivider = "------------"

	//path to zipped world directory
	filePath = "/data/world-backup-%s.tar.gz"
	//world directory
	worldDir = "/data/world/"
	//backup name skeleton
	fileName = "world-backup-%s.tar.gz"

	//expire date range
	expirationRange = -1 //delete after upload (no extra zip files kept in minio, rolling backups)
)

type Cron struct {
	job *cronjob.Cron
}

type ICron interface {
	Run()
	Stop()
	takeBackup()
}

// creates a new cron handler and add backup function
func NewCronHandler() *Cron {
	c := cronjob.New()
	cr := &Cron{}
	c.AddFunc("@daily", cr.takeBackup)
	cr.job = c
	return cr
}

func (c *Cron) Run() {
	c.job.Start()
}

func (c *Cron) Stop() {
	c.job.Stop()
}

// the main backup function for the cronjob, this func is huge and needs to be split up
func (c *Cron) takeBackup() {
	//logging date for easy searching
	log.Printf("%s/%s/%s\n", logDivider, time.Now().Format("01-02-2006"), logDivider)

	//create rcon client
	port, err := strconv.Atoi(os.Getenv("RCON_PORT"))
	if err != nil {
		log.Println(err)
		return
	}
	rcon := minecraft.NewRconHandler(os.Getenv("RCON_MC_CONTAINER"), port)
	err = rcon.Connect(os.Getenv("RCON_PASSWORD"))
	if err != nil {
		log.Println(err)
		return
	}
	defer rcon.Close()

	log.Println("connected to server")

	//send alert to server
	err = rcon.AlertPlayers("backup's will be taken in 1 minute")
	if err != nil {
		log.Println(err)
		return
	}

	//waiting
	log.Println("waiting")
	time.Sleep(time.Minute * 1)

	err = rcon.AlertPlayers("backing up now, buckle up")
	if err != nil {
		log.Println(err)
		return
	}

	//disabling autosaves for a full world save
	log.Println("autosaves disabled")
	err = rcon.DisableAutosaves()
	if err != nil {
		log.Println(err)
		return
	}
	defer rcon.EnableAutosaves() //ensure to enable them after backup

	//full save
	log.Println("taking full save")
	err = rcon.WorldSave()
	if err != nil {
		log.Println(err)
		return
	}

	//time shenanigans
	curr := time.Now()
	currStr := curr.Format("02-01-2006")

	//expiration shenanigans
	expired := curr.AddDate(0, 0, expirationRange)
	expStr := expired.Format("02-01-2006")

	//filename shenanigans
	savePath := fmt.Sprintf(filePath, currStr)
	saveFile := fmt.Sprintf(fileName, currStr)
	expFile := fmt.Sprintf(fileName, expStr)

	//zip world file
	log.Println("zipping world folder")
	fly := qz.NewZipper(worldDir)
	_, err = fly.Zip(savePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer os.Remove(savePath)

	//minio upload
	mc, err := file.NewCloudHandler(os.Getenv("MINIO_URL"), os.Getenv("MINIO_ID"), os.Getenv("MINIO_KEY"), os.Getenv("MINIO_BUCKET"))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("uploading backup to storage")
	err = mc.Upload(saveFile)
	if err != nil {
		log.Println(err)
		return
	}

	//cleanup expired backups
	log.Println("cleanup")
	mc.Delete(expFile)

	log.Println("done")

	//do you want us to restart the server
	if os.Getenv("SERVER_RESTART") == "yes" {
		//restart server
		log.Println("restarting server")
		err = rcon.AlertPlayers("the server will be restarted in 5 seconds")
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(time.Second * 5)

		//trigger system restart to avoid lag
		err = rcon.RestartServer()
		if err != nil {
			log.Println(err)
			return
		}
	}
}
