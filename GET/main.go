package main

import (
	"encoding/json"
	"log"

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
const dsn = "admin:Password123@tcp(mysql-rds-1.c5eooawm67do.us-east-1.rds.amazonaws.com:3306)/inventorysystem?charset=utf8mb4&parseTime=True&loc=Local"

// Initialize the database connection
func initializeDatabase() {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
}

// Lambda handler function
func GetProductsHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var products []Product

	// Fetch products from the database
	if err := DB.Find(&products).Error; err != nil {
		log.Printf("Error fetching products: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to fetch products"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Marshal products to JSON
	body, _ := json.Marshal(products)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	initializeDatabase()
	lambda.Start(GetProductsHandler)
}
