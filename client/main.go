package main

import (
	"flag"
	"log"
	"net/http"

	pb "example.com/go-grpc-product-management-system/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// define address for server connection
var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type Product struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	SubCategory string `json:"sub_category"`
	Price       string `json:"price"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	IsActive    string `json:"is_active"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
}

// main function that Start client Fizz
// Defining all routers
// This routers will be invoked by the real client, Postman , anything
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewProductServiceClient(conn)

	r := gin.Default()
	r.GET("/products", func(ctx *gin.Context) {
		res, err := client.GetProducts(ctx, &pb.ReadProductsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"product": res.Product,
		})
	})
	r.GET("/product/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		res, err := client.GetProduct(ctx, &pb.ReadProductRequest{Product: &pb.Product{Id: id}})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"product": res.Product,
		})
	})
	r.POST("/product", func(ctx *gin.Context) {
		var product Product

		err := ctx.ShouldBind(&product)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.Product{
			Name:        product.Name,
			Category:    product.Category,
			SubCategory: product.SubCategory,
			Price:       product.Price,
			IsActive:    product.IsActive,
			CreatedBy:   product.CreatedBy,
			UpdatedBy:   product.UpdatedBy,
		}
		res, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
			Product: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"product": res.Product,
		})
	})
	r.PUT("/product/:id", func(ctx *gin.Context) {
		var product Product
		err := ctx.ShouldBind(&product)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.UpdateProduct(ctx, &pb.UpdateProductRequest{
			Product: &pb.Product{
				Id:          product.Id,
				Name:        product.Name,
				Category:    product.Category,
				SubCategory: product.SubCategory,
				Price:       product.Price,
				IsActive:    product.IsActive,
				CreatedBy:   product.CreatedBy,
				UpdatedBy:   product.UpdatedBy,
			},
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"product": res.Product,
		})
		return

	})
	r.DELETE("/product/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if res.Success == true {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "product deleted successfully",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "error deleting product",
			})
			return
		}

	})

	r.Run(":5000")

}
