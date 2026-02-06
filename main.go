package main

import (
	"log"
	"net/http"
	"os"

	"cashier/database"
	"cashier/handlers"
	"cashier/repositories"
	"cashier/services"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	database.Connect()
	defer database.DB.Close()

	// Initialize layers for categories
	categoryRepo := repositories.NewCategoryRepository(database.DB)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Initialize layers for products
	productRepo := repositories.NewProductRepository(database.DB)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize layers for transactions
	transactionRepo := repositories.NewTransactionRepository(database.DB)
	transactionService := services.NewTransactionService(transactionRepo, productRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Initialize layers for reports
	reportService := services.NewReportService(transactionRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	r := mux.NewRouter()

	// Routes for categories
	r.HandleFunc("/categories", categoryHandler.GetCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	r.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Routes for products
	r.HandleFunc("/products", productHandler.GetProducts).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	r.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Routes for transactions
	r.HandleFunc("/checkout", transactionHandler.ProcessCheckout).Methods("POST")

	// Routes for reports
	r.HandleFunc("/report/hari-ini", reportHandler.GetTodaysSalesReport).Methods("GET")
	r.HandleFunc("/report", reportHandler.GetSalesReport).Methods("GET")

	log.Printf("Server is about to start on :%s", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Println("Server stopped")
}
