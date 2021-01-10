package plg

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gwuah/postmates/database/models"
	"github.com/gwuah/postmates/shared"
	"gorm.io/gorm"
)

func S(db *gorm.DB) {

	var data shared.GetDeliveryCostRequest
	err := json.Unmarshal([]byte(
		`{
			"origin": {
				"latitude": 5.677474538991623,
				"longitude": -0.24460022375167725
			},
			"destination": {
				"latitude": 5.6796946725653745,
				"longitude": -0.2447180449962616
			}
		}`,
	), &data)

	if err != nil {
		log.Fatal("error", err)
	}

	for i := 0; i < 10; i++ {
		cId := uint(1)
		delivery := models.Delivery{
			OriginLatitude:       data.Origin.Latitude,
			OriginLongitude:      data.Origin.Longitude,
			DestinationLatitude:  data.Destination.Latitude,
			DestinationLongitude: data.Destination.Longitude,
			Notes:                "Hello",
			CustomerID:           1,
			State:                models.Pending,
			CourierID:            &cId,
			CustomerRating:       1,
			ProductID:            1,
		}

		if err := db.Create(&delivery).Error; err != nil {
			log.Fatal("error", err)
		}

	}

	fmt.Println("seed complete")
}

func C(db *gorm.DB) {

	var data shared.GetDeliveryCostRequest
	err := json.Unmarshal([]byte(
		`{
			"origin": {
				"latitude": 5.677474538991623,
				"longitude": -0.24460022375167725
			},
			"destination": {
				"latitude": 5.6796946725653745,
				"longitude": -0.2447180449962616
			}
		}`,
	), &data)

	if err != nil {
		log.Fatal("error", err)
	}

	for i := 0; i < 10; i++ {
		cId := uint(1)
		delivery := models.Delivery{
			OriginLatitude:       data.Origin.Latitude,
			OriginLongitude:      data.Origin.Longitude,
			DestinationLatitude:  data.Destination.Latitude,
			DestinationLongitude: data.Destination.Longitude,
			Notes:                "Hello",
			CustomerID:           1,
			State:                models.Pending,
			CourierID:            &cId,
			CourierRating:        1,
			ProductID:            1,
		}

		if err := db.Create(&delivery).Error; err != nil {
			log.Fatal("error", err)
		}

	}

	fmt.Println("seed complete")
}
