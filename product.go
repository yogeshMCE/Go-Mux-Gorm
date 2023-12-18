package main

import (
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type Book struct {
	ID       int     `json:"id" `
	Name     string  `gorm:"" json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func GetAllBooks(db *gorm.DB) []Book {
	var books []Book
	db.Find(&books)

	return books
}
func GetBookById(Id int64, db *gorm.DB) *Book {
	var getbook Book
	db.Where("ID=?", Id).First(&getbook)
	return &getbook
}
func (b *Book) CreateBook(db *gorm.DB) (*Book, error) {

	book := db.Create(&b)
	log.Printf("Book created %v", book)
	if book.RowsAffected == 0 {
		return nil, book.Error
	}
	return b, nil
}
func (b *Book) UpdateBook(ID int64, db *gorm.DB) Book {
	var getbook Book
	db.Where("ID=?", ID).First(&getbook)
	getbook.Quantity = b.Quantity
	getbook.Name = b.Name
	getbook.Price = b.Price
	db.Save(&getbook)
	return getbook

}
func DeleteBook(ID int64, db *gorm.DB) Book {
	var book Book
	db.Where("ID=?", ID).First(&book)
	db.Delete(&Book{}, ID)
	return book
}
