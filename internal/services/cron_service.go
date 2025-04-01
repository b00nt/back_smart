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
func (cs *CronService) Start(city string, db *gorm.DB) {
	// Run the task immediately once at the start
	resultProducts := moysklad.GetProducts(city)
	err := moysklad.SaveProducts(city, resultProducts, db)
	if err != nil {
		log.Println("Error updating product:", err)
	} else {
		fmt.Println("Product update successful")
	}

	resultVariants := moysklad.GetModifications(city, db)
	err = moysklad.SaveModifications(city, resultVariants, db)
	if err != nil {
		log.Println("Error updating modifications:", err)
	} else {
		fmt.Println("Modifications update successful")
	}

	resultStock := moysklad.GetStock(city)
	err = moysklad.SaveStock(city, resultStock, db)
	if err != nil {
		log.Println("Error updating stock:", err)
	} else {
		fmt.Println("Stock update successful")
	}

	// Download and save images for products and modifications
	moysklad.GetSaveDownloadProductImages(city, db)
	moysklad.GetSaveDownloadModImages(city, db)

	// Schedule the update job to run every day at midnight (or any desired schedule)
	_, err = cs.cronScheduler.AddFunc("@every 24h", func() {
		fmt.Println("Starting daily update job")

		resultProducts := moysklad.GetProducts(city)
		err := moysklad.SaveProducts(city, resultProducts, db)
		if err != nil {
			log.Println("Error updating product:", err)
		} else {
			fmt.Println("Product update successful")
		}

		resultVariants := moysklad.GetModifications(city, db)
		err = moysklad.SaveModifications(city, resultVariants, db)
		if err != nil {
			log.Println("Error updating modifications:", err)
		} else {
			fmt.Println("Modifications update successful")
		}

		resultStock := moysklad.GetStock(city)
		err = moysklad.SaveStock(city, resultStock, db)
		if err != nil {
			log.Println("Error updating modifications:", err)
		} else {
			fmt.Println("Modifications update successful")
		}

		// Download and save images for products and modifications
		moysklad.GetSaveDownloadProductImages(city, db)
		moysklad.GetSaveDownloadModImages(city, db)
	})

	if err != nil {
		log.Fatal("Error scheduling cron job:", err)
	}

	// Start the cron job scheduler
	cs.cronScheduler.Start()

	fmt.Println("Cron job started.")
}
