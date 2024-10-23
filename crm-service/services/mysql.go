package services

import (
	"context"
	"crm-service/config"
	"crm-service/internals/dto"
	"crm-service/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var ctx = context.Background()

func AddCustomer(customer *models.Customer, config *config.Config) error {
	db := config.DB
	if err := db.Create(customer).Error; err != nil {
		return err
	}
	redis := config.Redis
	return CacheCustomer(customer, redis)
}

func GetCustomers(db *gorm.DB) ([]models.Customer, error) {
	var customers []models.Customer
	if err := db.Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

func FetchById(id int, db *gorm.DB) (models.Customer, error) {
	var customer models.Customer
	if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
		return customer, err
	}
	return customer, nil
}

func UpdateCustomer(id int, post dto.CustomerRequest, config *config.Config) (models.Customer, error) {
	db := config.DB
	customer := models.Customer{
		ID:        id,
		FirstName: post.FirstName,
		LastName:  post.LastName,
		Company:   post.Company,
		City:      post.City,
		Postal:    post.Postal,
		Phone:     post.Phone,
		Email:     post.Email,
		Address:   post.Address,
		Web:       post.Web,
		County:    post.County,
		UpdatedAt: time.Now(),
	}
	if err := db.Model(&models.Customer{}).Where("id = ?", id).Updates(&customer).Error; err != nil {
		return models.Customer{}, err
	}
	redis := config.Redis
	return customer, CacheCustomer(&customer, redis)
}

func DeleteCustomer(id int, config *config.Config) error {
	db := config.DB
	if err := db.Delete(&models.Customer{}, id).Error; err != nil {
		return err
	}

	redisKey := fmt.Sprintf("customer:%d", id)
	redis := config.Redis
	if err := redis.Del(context.Background(), redisKey).Err(); err != nil {
		return err
	}

	return nil
}

func GetAllCustomersFromCache(db *gorm.DB, redisClient *redis.Client) ([]models.Customer, error) {
	var customers []models.Customer
	var customerIDs []int

	// Fetch all customer IDs from the database
	if err := db.Model(&models.Customer{}).Pluck("id", &customerIDs).Error; err != nil {
		return nil, err
	}

	// Build Redis keys for each customer ID
	customerKeys := make([]string, len(customerIDs))
	for i, id := range customerIDs {
		customerKeys[i] = "customer:" + strconv.Itoa(id)
	}

	// Use MGET to retrieve all customers from Redis in a single call
	customerJSONs, err := redisClient.MGet(ctx, customerKeys...).Result()
	if err != nil {
		return nil, err
	}
	// Track IDs of customers that are not found in the cache
	var missingCustomerIDs []int
	for i, customerJSON := range customerJSONs {
		if customerJSON == nil {
			// If the result is nil, it means the customer was not found in the cache
			missingCustomerIDs = append(missingCustomerIDs, customerIDs[i])
		} else {
			// Deserialize the JSON into a Customer object
			var customer models.Customer
			if err := json.Unmarshal([]byte(customerJSON.(string)), &customer); err == nil {
				customers = append(customers, customer)
			}
		}
	}

	// Define a chunk size to avoid too many placeholders
	chunkSize := 1000

	// Fetch missing customers in chunks
	for i := 0; i < len(missingCustomerIDs); i += chunkSize {
		end := i + chunkSize
		if end > len(missingCustomerIDs) {
			end = len(missingCustomerIDs)
		}
		chunk := missingCustomerIDs[i:end]

		var missingCustomers []models.Customer
		if err := db.Where("id IN ?", chunk).Find(&missingCustomers).Error; err != nil {
			return nil, err
		}

		// Add the missing customers to the result list
		customers = append(customers, missingCustomers...)

		// Cache the missing customers back into Redis in a batch
		pipe := redisClient.Pipeline()
		for _, customer := range missingCustomers {
			customerJSON, _ := json.Marshal(customer)
			customerKey := "customer:" + strconv.Itoa(int(customer.ID))
			pipe.Set(ctx, customerKey, customerJSON, 5*time.Minute)
		}
		_, _ = pipe.Exec(ctx)
	}

	return customers, nil
}
