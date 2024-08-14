package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ProductName     string  `json:"prodname"`
	ProductCategory string  `json:"prodcategory"`
	ProductPrice    float64 `json:"prodprice"`
	ProductStock    int     `json:"prodstock"`
}

var DB *gorm.DB

// Database connection string
const dsn = "admin:Vidhya_14@tcp(inventorysystem.cbmag2acul28.us-east-1.rds.amazonaws.com:3306)/inventorysystem?charset=utf8mb4&parseTime=True&loc=Local"

// Initialize the database connection
func initializeDatabase() {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	// Auto-migrate schema
	if err := DB.AutoMigrate(&Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// Handle PUT requests to update a product by ID
func updateProductHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid product ID"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var updatedProduct Product
	if err := json.Unmarshal([]byte(request.Body), &updatedProduct); err != nil {
		log.Printf("Error unmarshaling request body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `{"error": "Invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Fetch the existing product from the database
	var existingProduct Product
	if err := DB.First(&existingProduct, id).Error; err != nil {
		log.Printf("Error fetching product with ID %d: %v", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "Product not found"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Update the existing product with the new data
	existingProduct.ProductName = updatedProduct.ProductName
	existingProduct.ProductCategory = updatedProduct.ProductCategory
	existingProduct.ProductPrice = updatedProduct.ProductPrice
	existingProduct.ProductStock = updatedProduct.ProductStock

	if err := DB.Save(&existingProduct).Error; err != nil {
		log.Printf("Error updating product with ID %d: %v", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to update product"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	body, _ := json.Marshal(existingProduct)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	initializeDatabase()
	lambda.Start(updateProductHandler)
}
