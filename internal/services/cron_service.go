// services/cron_service.go

package services

import (
	"back/internal/moysklad"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
)

// CronService struct to hold necessary information for the cron job
type CronService struct {
	cronScheduler *cron.Cron
}

// NewCronService initializes and returns a new CronService
func NewCronService() *CronService {
	return &CronService{
		cronScheduler: cron.New(),
	}
}

// Start initializes the cron job and starts it
func (cs *CronService) Start(db *gorm.DB) {
	// Run the task immediately once at the start
	resultProducts := moysklad.GetProducts()
	err := moysklad.SaveProducts(resultProducts, db)
	if err != nil {
		log.Println("Error updating product:", err)
	} else {
		fmt.Println("Product update successful")
	}

	resultVariants := moysklad.GetModifications()
	err = moysklad.SaveModifications(resultVariants, db)
	if err != nil {
		log.Println("Error updating modifications:", err)
	} else {
		fmt.Println("Modifications update successful")
	}

	resultStock := moysklad.GetStock()
	err = moysklad.SaveStock(resultStock, db)
	if err != nil {
		log.Println("Error updating modifications:", err)
	} else {
		fmt.Println("Modifications update successful")
	}

	moysklad.GetSaveDownloadProductImages(db)
	moysklad.GetSaveDownloadModImages(db)

	// Schedule the update job to run every day at midnight (or any desired schedule)
	_, err = cs.cronScheduler.AddFunc("@every 24h", func() {
		fmt.Println("Starting daily update job")

		resultProducts := moysklad.GetProducts()
		err := moysklad.SaveProducts(resultProducts, db)
		if err != nil {
			log.Println("Error updating product:", err)
		} else {
			fmt.Println("Product update successful")
		}

		resultVariants := moysklad.GetModifications()
		err = moysklad.SaveModifications(resultVariants, db)
		if err != nil {
			log.Println("Error updating modifications:", err)
		} else {
			fmt.Println("Modifications update successful")
		}

		resultStock := moysklad.GetStock()
		err = moysklad.SaveStock(resultStock, db)
		if err != nil {
			log.Println("Error updating modifications:", err)
		} else {
			fmt.Println("Modifications update successful")
		}

		moysklad.GetSaveDownloadProductImages(db)
		moysklad.GetSaveDownloadModImages(db)
	})

	if err != nil {
		log.Fatal("Error scheduling cron job:", err)
	}

	// Start the cron job scheduler
	cs.cronScheduler.Start()

	fmt.Println("Cron job started.")
}
