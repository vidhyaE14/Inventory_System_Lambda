package main

import (
    "encoding/json"
    "log"
    "os"
	"fmt"
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

// Retrieve the database connection string from environment variables
func getDSN() string {
    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    database := os.Getenv("DB_NAME")

    if username == "" || password == "" || host == "" || port == "" || database == "" {
        log.Fatalf("One or more environment variables are not set")
    }

    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        username, password, host, port, database)
}

// Initialize the database connection
func initializeDatabase() {
    var err error
    dsn := getDSN() // Use environment variables to get DSN
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
