package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID    int `gorm:"primaryKey"`
	Name  string
	Price float64
	gorm.Model
}

func main() {
	dsn := "rafael:vieira@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Product{})

	db.Create(&Product{
		Name:  "Monitor",
		Price: 325.99,
	})
	//-----------------
	products := []Product{
		{Name: "Celular", Price: 3500.90},
		{Name: "Mouse", Price: 25.00},
		{Name: "Keyboard", Price: 30.00},
	}
	db.Create(&products)
	//-----------------
	//var product Product
	//db.First(&product, 1)
	//db.First(&product, "name = ?", "Mouse")
	//fmt.Println(product)
	//-----------------
	/* 	var products []Product
	   	db.Find(&products)
	   	for _, product := range products {
	   		fmt.Println(product)
	   	} */
	/*var products []Product
	db.Limit(2).Offset(2).Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}*/

	/* var products []Product
	db.Where("price > ?", 28).Find(&products)
	for _, product := range products {
		fmt.Println(product)
	} */

	/* var products []Product
	db.Where("name like ?", "%o%").Find(&products)
	for _, product := range products {
		fmt.Println(product)
	} */
	/*
		var p Product
		db.First(&p, 1)
		p.Name = "New Monitor"
		db.Save(&p)*/
	/*
		var p2 Product
		db.First(&p2, 1)
		fmt.Println(p2)*/
	/*
		var p Product
		db.First(&p, 1)
		db.Delete(&p)*/
}
