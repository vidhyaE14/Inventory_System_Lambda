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
	// Auto-migrate schema
	if err := DB.AutoMigrate(&Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// Handle POST requests to add a new product
func postProductHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Printf("Received request body: %s", request.Body)

    var newProduct Product
    if err := json.Unmarshal([]byte(request.Body), &newProduct); err != nil {
        log.Printf("Error unmarshaling request body: %v", err)
        return events.APIGatewayProxyResponse{
            StatusCode: 400,
            Body:       `{"error": "Invalid request body"}`,
            Headers:    map[string]string{"Content-Type": "application/json"},
        }, nil
    }

    log.Printf("Parsed product: %+v", newProduct)

    if err := DB.Create(&newProduct).Error; err != nil {
        log.Printf("Error creating product: %v", err)
        return events.APIGatewayProxyResponse{
            StatusCode: 500,
            Body:       `{"error": "Failed to create product"}`,
            Headers:    map[string]string{"Content-Type": "application/json"},
        }, nil
    }

    body, _ := json.Marshal(newProduct)
    return events.APIGatewayProxyResponse{
        StatusCode: 201,
        Body:       string(body),
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}

func main() {
	initializeDatabase()
	lambda.Start(postProductHandler)
}
