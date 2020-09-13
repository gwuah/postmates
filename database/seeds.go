package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gwuah/api/database/models"
	"github.com/kylelemons/go-gypsy/yaml"
	"gorm.io/gorm"
)

type SeedFn func(db *gorm.DB, path string)

func RunSeeds(db *gorm.DB, seeds []SeedFn) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for _, seed := range seeds {
		seed(db, path)
	}
}

func SeedProducts(DB *gorm.DB, path string) {
	config, err := yaml.ReadFile(path + "/database/products.yml")
	if err != nil {
		panic(err)
	}

	productList, ok := config.Root.(yaml.List)
	if !ok {
		panic("failed to parse product.yml")
	}

	for i := 0; i < productList.Len(); i++ {
		productName := strings.ToLower(fmt.Sprintf("%s", productList.Item(i)))
		var product models.Product

		if err := DB.Where("name = ?", productName).First(&product).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&models.Product{Name: productName})
			} else {
				log.Printf("Product [ %s ] lookup failed", productName)
				log.Println(err)
			}
		}
	}
}

func SeedElectrons(DB *gorm.DB, path string) {
	config, err := yaml.ReadFile(path + "/database/electrons.yml")
	if err != nil {
		panic(err)
	}

	electronList, ok := config.Root.(yaml.List)
	if !ok {
		panic("failed to parse product.yml")
	}

	for i := 0; i < electronList.Len(); i++ {
		name := strings.ToLower(fmt.Sprintf("%s", electronList.Item(i)))
		firstName := strings.Split(name, " ")[0]
		lastName := strings.Split(name, " ")[1]

		var electron models.Electron

		if err := DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&electron).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&models.Electron{FirstName: firstName, LastName: lastName})
			} else {
				log.Printf("Electron [ %s ] lookup failed", name)
				log.Println(err)
			}
		}
	}
}

func SeedCustomers(DB *gorm.DB, path string) {
	config, err := yaml.ReadFile(path + "/database/customers.yml")
	if err != nil {
		panic(err)
	}

	customerList, ok := config.Root.(yaml.List)
	if !ok {
		panic("failed to parse product.yml")
	}

	for i := 0; i < customerList.Len(); i++ {
		name := strings.ToLower(fmt.Sprintf("%s", customerList.Item(i)))
		firstName := strings.Split(name, " ")[0]
		lastName := strings.Split(name, " ")[1]

		var customer models.Customer

		if err := DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&customer).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				DB.Create(&models.Customer{FirstName: firstName, LastName: lastName})
			} else {
				log.Printf("Customer [ %s ] lookup failed", name)
				log.Println(err)
			}
		}
	}
}
