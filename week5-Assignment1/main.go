package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Product struct {
    ID       		string  	`json:"id"`
    Name     		string  	`json:"name"`
	Type 			string 		`json:"type"`
	Manufacturer 	string 		`json:"manufacturer"`
    Price     		float64 	`json:"price"`
    Store      		int 		`json:"store"`
	Sold			int 		`json:"Sold"`
	
}

var stocks = []Product{
    {ID: "1", Name: "Pepsi", Type: "Drinks", Manufacturer: "PepsiCo.", Price: 24.0, Store: 100, Sold: 60},
    {ID: "2", Name: "Coca-Cola", Type: "Drinks", Manufacturer: "The Coca-Cola Company", Price: 20.0, Store: 200, Sold: 90},
	{ID: "3", Name: "Puriku", Type: "Drinks", Manufacturer: "TCP.", Price: 12.0, Store: 200, Sold: 150},
	{ID: "4", Name: "Lay's", Type: "Snacks", Manufacturer: "Frito-Lay", Price: 20.0, Store: 500, Sold: 290},
	{ID: "5", Name: "Tasto", Type: "Snacks", Manufacturer: "BJC.", Price: 20.0, Store: 200, Sold: 90},
}

func getProduct_Name(c *gin.Context){
	NameQuery := c.Query("Name")

	if NameQuery != "" {
		filter := []Product{}
		for _, Product := range stocks{
			if fmt.Sprint(Product.Name) == NameQuery{
				filter = append(filter, Product)
			}
		}
		c.JSON(http.StatusOK, filter)
		return 
	}
	c.JSON(http.StatusOK, stocks)
}

func getProduct_Type(c *gin.Context){
		TypeQuery := c.Query("Type")

	if TypeQuery != "" {
		filter := []Product{}
		for _, Product := range stocks{
			if fmt.Sprint(Product.Type) == TypeQuery{
				filter = append(filter, Product)
			}
		}
		c.JSON(http.StatusOK, filter)
		return 
	}
	c.JSON(http.StatusOK, stocks)

}

func getProduct_Sold(c *gin.Context){
		SoldQuery := c.Query("Sold")

	if SoldQuery != "" {
		filter := []Product{}
		for _, Product := range stocks{
			if fmt.Sprint(Product.Sold) == SoldQuery{
				filter = append(filter, Product)
			}
		}
		c.JSON(http.StatusOK, filter)
		return 
	}
	c.JSON(http.StatusOK, stocks)
}

func main (){
	r := gin.Default()

	r.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{"message" : "healthy"})
	})

	api := r.Group("/api/v1/stocks")
	{
		api.GET("/name", getProduct_Name)
		api.GET("/type", getProduct_Type)
		api.GET("/sold", getProduct_Sold)
	}

	r.Run(":8080")
}