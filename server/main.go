package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "example.com/go-grpc-product-management-system/proto"
)
//main function to start server
func main() {
	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterProductServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
//defining ports for server
var (
	port = flag.Int("port", 50051, "gRPC server port")
)
//implement productService
type server struct {
	pb.UnimplementedProductServiceServer
}

//init DB connection
func init() {
	DatabaseConnection()
}

var DB *gorm.DB
var err error
//product structure
type Product struct {
	ID          string `gorm:"primarykey"`
	Name        string
	Category    string
	SubCategory string
	Price       string
	CreatedAt   time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:true"`
	IsActive    string
	CreatedBy   string
	UpdatedBy   string
}

//database connection data
//improve with Secrets
func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "postgres"
	dbUser := "postgres"
	password := "Vasile94!"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(Product{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")
}

//Server function to CreateProduct
// TO DO
//transform CreatedAt and UpdatedAt from string to DATE.Time
func (*server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	fmt.Println("Create Product")
	product := req.GetProduct()
	product.Id = uuid.New().String()

	data := Product{
		ID:          product.GetId(),
		Name:        product.GetName(),
		Category:    product.GetCategory(),
		SubCategory: product.GetSubCategory(),
		Price:       product.GetCreatedDate(),
		//CreatedAt:   product.GetCreatedDate(),
		//UpdatedAt:   product.GetUpdatedDate(),
		IsActive:  product.GetIsActive(),
		CreatedBy: product.GetCreatedBy(),
		UpdatedBy: product.GetUpdatedBy(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("Product creation unsuccessful")
	}
	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          product.GetId(),
			Name:        product.GetName(),
			Category:    product.GetCategory(),
			SubCategory: product.GetSubCategory(),
			Price:       product.GetCreatedDate(),
			//CreatedAt:   product.GetCreatedDate(),
			//UpdatedAt:   product.GetUpdatedDate(),
			IsActive:  product.GetIsActive(),
			CreatedBy: product.GetCreatedBy(),
			UpdatedBy: product.GetUpdatedBy(),
		},
	}, nil
}
//Server function to GetProduct
// TO DO
//transform CreatedAt and UpdatedAt from string to DATE.Time
func (*server) GetProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.ReadProductResponse, error) {
	fmt.Println("Read Product", req.Product.GetId())
	var product Product
	res := DB.Find(&product, "id = ?", req.Product.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("Product not found")
	}
	return &pb.ReadProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Category:    product.Category,
			SubCategory: product.SubCategory,
			Price:       product.Price,
			IsActive:    product.IsActive,
			CreatedBy:   product.CreatedBy,
		},
	}, nil
}
//Server function to GetProducts list, returns all products in DB
// TO DO
func (*server) GetProducts(ctx context.Context, req *pb.ReadProductsRequest) (*pb.ReadProductsResponse, error) {
	fmt.Println("Read Products")
	products := []*pb.Product{}
	res := DB.Find(&products)
	if res.RowsAffected == 0 {
		return nil, errors.New("Product not found")
	}
	return &pb.ReadProductsResponse{
		Product: products,
	}, nil
}
//Server function to UpdateProduct
// TO DO
//transform CreatedAt and UpdatedAt from string to DATE.Time
func (*server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	fmt.Println("Update Product")
	var product Product
	reqProduct := req.GetProduct()

	res := DB.Model(&product).Where("id=?", reqProduct.Id).Updates(Product{Name: reqProduct.Name, Category: reqProduct.Category, SubCategory: reqProduct.SubCategory, Price: reqProduct.Price})

	if res.RowsAffected == 0 {
		return nil, errors.New("Products not found")
	}

	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Category:    product.Category,
			SubCategory: product.SubCategory,
			Price:       product.Price,
			IsActive:    product.IsActive,
			CreatedBy:   product.CreatedBy,
		},
	}, nil
}
//Server function to DeleteProduct from DB
func (*server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	fmt.Println("Delete Product")
	var Product Product
	res := DB.Where("id = ?", req.GetId()).Delete(&Product)
	if res.RowsAffected == 0 {
		return nil, errors.New("Product not found")
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

