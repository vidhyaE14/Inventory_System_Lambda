package main

import (
    "encoding/json"
    "log"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "strconv"
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

// Lambda handler function to fetch a product by ID
func GetProductByIDHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

    var product Product
    // Fetch the product by ID from the database
    if err := DB.First(&product, id).Error; err != nil {
        log.Printf("Error fetching product with ID %d: %v", id, err)
        return events.APIGatewayProxyResponse{
            StatusCode: 404,
            Body:       `{"error": "Product not found"}`,
            Headers:    map[string]string{"Content-Type": "application/json"},
        }, nil
    }

    // Marshal product to JSON
    body, _ := json.Marshal(product)
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       string(body),
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}

func main() {
    initializeDatabase()
    lambda.Start(GetProductByIDHandler)
}
