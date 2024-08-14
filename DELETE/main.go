package main

import (
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

// Handle DELETE requests to delete a product by ID
func deleteProductHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	// Delete the product from the database
	result := DB.Delete(&Product{}, id)
	if result.RowsAffected == 0 {
		log.Printf("Error deleting product with ID %d: %v", id, result.Error)
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "Product not found"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Return success response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message": "Product deleted successfully"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	initializeDatabase()
	lambda.Start(deleteProductHandler)
}
