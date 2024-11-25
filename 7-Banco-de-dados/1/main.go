package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // O "_" no inicio do import é para não remover o import mesmo que não utilizado
	"github.com/google/uuid"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() { //usuario:senha@tcp(localhost:3303)/nome_database
	db, err := sql.Open("mysql", "rafael:vieira@tcp(localhost:3306)/goexpert")

	if err != nil {
		panic(err)
	}

	defer db.Close()
	product := NewProduct("Notebook 2", 1899.90)

	err = insertProduct(db, product)

	if err != nil {
		panic(err)
	}

	product.Price = 1920.50
	err = updateProduct(db, product)

	if err != nil {
		panic(err)
	}

	p, err := selectProduct(db, product.ID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Product: %v, possui o preço de %.2f\n", p.Name, p.Price)

	products, err := selectAllProducts(db)

	if err != nil {
		panic(err)
	}

	for _, p := range products {
		fmt.Printf("Produto: %v, possui preço de %.2f\n", p.Name, p.Price)
	}

	err = deleteProduct(db, product.ID)

	if err != nil {
		panic(err)
	}
}

func insertProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("INSERT INTO products(id, name, price) VALUES(?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(product.ID, product.Name, product.Price)

	if err != nil {
		return err
	}

	return nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(product.Name, product.Price, product.ID)

	if err != nil {
		return err
	}

	return nil
}

func selectProduct(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("SELECT id, name, price FROM products WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var p Product

	err = stmt.QueryRow(id).Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func selectAllProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price)

		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func deleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}
