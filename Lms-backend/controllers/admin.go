package controllers

import (
	"errors"
	"fmt"
	"lms/backend/initializers"
	"lms/backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddBooks struct {
	ISBN      string
	Title     string
	Author    string
	Publisher string
	Version   int
}

// AddBook function - fixed version with better error handling
func AddBook(c *gin.Context) {
	// Getting admin information
	email, _ := c.Get("email")
	adminID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verifying admin access
	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Parsing and validating book input
	var bookInput models.BookInventory
	if err := c.ShouldBindJSON(&bookInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input validation - important to check before database operations
	if bookInput.ISBN == "" || bookInput.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ISBN and Title are required fields"})
		return
	}

	// Check if book with this ISBN exists in any library
	var existingBook models.BookInventory
	result := initializers.DB.Where("isbn = ?", bookInput.ISBN).First(&existingBook)

	if result.Error == nil {
		// Book exists somewhere
		if existingBook.LibID == adminUser.LibID {
			// Book exists in this library, update copy count
			existingBook.TotalCopies++
			existingBook.AvailableCopies++
			if err := initializers.DB.Save(&existingBook).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book copies"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Book copies updated successfully", "book": existingBook, "email": email})
			return
		} else {
			// Book exists in a different library
			c.JSON(http.StatusBadRequest, gin.H{"error": "Same ISBN already exists in another library, cannot add book"})
			return
		}
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Database error that's not just "not found"
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + result.Error.Error()})
		return
	}

	// Book doesn't exist, create new one
	newBook := models.BookInventory{
		ISBN:            bookInput.ISBN,
		Title:           bookInput.Title,
		Authors:         bookInput.Authors,
		Publisher:       bookInput.Publisher,
		Version:         bookInput.Version,
		TotalCopies:     1,
		AvailableCopies: 1,
		LibID:           adminUser.LibID,
	}

	// Create the book in DB
	if err := initializers.DB.Create(&newBook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book: " + err.Error()})
		return
	}

	// Successfully added the book
	c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully", "book": newBook, "email": email})
}
func RemoveBook(c *gin.Context) {

	// get the email of the admin logged in and check the book which needs to updated is from the same library
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user models.User

	if err := initializers.DB.Where("email=?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": err.Error(),
		})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	isbn := c.Param("id")
	if isbn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ISBN"})
		return
	}

	var book models.BookInventory
	if err := initializers.DB.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found."})
		return
	}
	//check admin and bookid
	if user.LibID != book.LibID {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Libraries are different",
		})
		return
	}

	if book.AvailableCopies <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No available copies to remove."})
		return
	}
	// if book.TotalCopies <= 1 {
	// 	initializers.DB.Delete(&book)
	// 	c.JSON(http.StatusOK, gin.H{"message": "Book removed completely."})
	// 	return
	// }

	book.TotalCopies -= 1
	book.AvailableCopies -= 1
	if err := initializers.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book copies"})
		return
	}
	c.JSON(http.StatusOK,
		gin.H{"book": book, "Message": "Book deleted Successfully"})
}

type UpdatorBook struct {
	ISBN      string
	Title     string
	Authors   string
	Publisher string
	Version   int
}

func UpdateBook(c *gin.Context) {

	// get the email of the admin logged in and check the book which needs to updated is from the same library
	email, _ := c.Get("email")
	var user models.User

	if err := initializers.DB.Where("email=?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error": err.Error(),
		})
		return
	}
	userLibID := user.LibID

	isbn := c.Param("id")
	// fmt.Println(isbn)
	var upBook UpdatorBook

	if err := c.ShouldBindJSON(&upBook); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var book models.BookInventory
	if err := initializers.DB.Where("isbn = ?", isbn).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found."})
		return
	}

	//check admin and bookid
	if userLibID != book.LibID {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Libraries are different",
		})
		return
	}

	book.Title = upBook.Title
	book.Authors = upBook.Authors
	book.Publisher = upBook.Publisher
	book.Version = upBook.Version
	initializers.DB.Save(&book)
	c.JSON(http.StatusOK, book)
}

func ListRequests(c *gin.Context) {
	// Getting the admin ID from the context
	adminID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch the admin details
	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// var onerequest models.RequestEvent
	var requests []models.RequestEvent

	if err := initializers.DB.
		Joins("JOIN book_inventories ON request_events.book_id = book_inventories.isbn").
		Where("book_inventories.lib_id = ? AND request_events.request_type = ?", adminUser.LibID, "Requested").
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Could not find the requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": requests,
	})
}

type UpdateRequest struct {
	RequestType  string     `json:"request_type"`
	ApprovalDate *time.Time `json:"approval_date,omitempty"`
	ApproverID   *uint      `json:"approver_id,omitempty"`
}

func ApproveRequest(c *gin.Context) {
	adminID, _ := c.Get("id")
	//	fmt.Println(adminID)

	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	//extract the request id and check its availablity or not
	requestId := c.Param("id")
	//fmt.Println(requestId)
	var request models.RequestEvent

	if err := initializers.DB.Where("req_id=?", requestId).Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnot find the request with this id",
		})
		return
	}
	//time to extract bookId from request id
	bookId := request.BookID
	//fmt.Println(bookId)

	var bookexists models.RequestEvent
	if err := initializers.DB.Where("book_id=? AND req_id=?", bookId, requestId).First(&bookexists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnt find the book id with this isbn",
		})
		return
	}
	var bookCopies models.BookInventory

	if err := initializers.DB.Where("isbn=?", bookexists.BookID).Find(&bookCopies).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Coudlnt find the book",
		})
		return
	}

	fmt.Println("Book copies are", bookCopies.AvailableCopies)

	if bookCopies.AvailableCopies < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "No copies available",
		})
		return
	}

	updateCopies := bookCopies.AvailableCopies - 1
	if err := initializers.DB.Model(&bookCopies).Where("isbn=?", bookId).Update("available_copies", updateCopies).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnot update copies",
		})
		return
	}

	// to do jo request id aai hai 10 usko update kr de requested to issued

	var handlereq UpdateRequest

	if err := c.ShouldBindJSON(&handlereq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	//time to update the bookrequest
	bookexists.RequestType = "Issued"
	now := time.Now()
	bookexists.ApprovalDate = &now
	bookexists.ApproverID = &adminUser.ID

	if err := initializers.DB.Save(&bookexists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnt update book in request events",
		})
		return
	}
	// c.JSON(http.StatusOK, bookexists)
	c.JSON(http.StatusAccepted, gin.H{
		"message":      "updation successfully done",
		"updated book": bookexists,
	})

	//now setup the issue registry accordingly
	issueReg := models.IssueRegistry{
		ISBN:               bookId,
		ReaderID:           request.ReaderID,
		IssueApproverID:    adminUser.ID,
		IssueStatus:        "Issued",
		ExpectedReturnDate: time.Now().AddDate(0, 0, 14), // 2 weeks later

	}
	initializers.DB.Create(&issueReg)

	c.JSON(http.StatusCreated, gin.H{
		"Message":   "Issue registry created successfully",
		"Issue reg": issueReg,
	})

}
func RejectRequest(c *gin.Context) {
	adminID, _ := c.Get("id")
	//	fmt.Println(adminID)

	var adminUser models.User
	if err := initializers.DB.Where("id = ? AND role = ?", adminID, "admin").First(&adminUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	//extract the request id and check its availablity or not
	requestId := c.Param("id")
	// fmt.Println("Req")
	// fmt.Println(requestId)
	var request models.RequestEvent

	if err := initializers.DB.Where("req_id=?", requestId).Find(&request).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnot find the request with this id",
		})
		return
	}
	//time to extract bookId from request id
	bookId := request.BookID
	// fmt.Println("book")
	// fmt.Println(bookId)

	var bookexists models.RequestEvent
	if err := initializers.DB.Where("book_id=? AND req_id=?", bookId, requestId).Find(&bookexists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Couldnt find the book id with this isbn",
		})
		return
	}

	// readerID := request.ReaderID
	// fmt.Println(readerID)

	var handlereq UpdateRequest

	if err := c.ShouldBindJSON(&handlereq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	//time to update the bookrequest
	fmt.Println(handlereq.RequestType)
	bookexists.RequestType = "Rejected"
	now := time.Now()
	bookexists.ApprovalDate = &now
	bookexists.ApproverID = &adminUser.ID

	initializers.DB.Save(&bookexists)
	// c.JSON(http.StatusOK, bookexists)
	c.JSON(http.StatusAccepted, gin.H{
		"message":      "updation successfully done",
		"updated book": bookexists,
	})
}

func GetAllBooks(c *gin.Context) {
	// var book models.BookInventory
	adminID, _ := c.Get("id")
	fmt.Println(adminID)

	var user models.User

	if err := initializers.DB.Where("ID=?", adminID).Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"Message": "User not found",
		})
		return
	}
	// userLibId := user.LibID

	var getBooks []models.BookInventory

	// if err := initializers.DB.Where("lib_id=? AND available_copies > ?", user.LibID, 0).Find(&getBooks).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"Error":   err.Error(),
	// 		"Message": "Unable to fetch all books",
	// 	})
	// 	return
	// }
	if err := initializers.DB.Where("lib_id=?", user.LibID).Find(&getBooks).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Error":   err.Error(),
			"Message": "Unable to fetch all books",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Books": getBooks,
	})
}
func IssueInfo(c *gin.Context) {
	readerID := c.Param("id")  //reader ke libID
	email, _ := c.Get("email") //admin ki libId

	//check all where admin Lib Id will be equal to reader LibId
	var adminDetails models.User

	if err := initializers.DB.Where("email=?", email).First(&adminDetails).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"Error": err.Error(),
		})
		return
	}

	var readerDetails models.User
	if err := initializers.DB.Where("id=?", readerID).First(&readerDetails).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"Error": err.Error(),
		})
		return
	}

	adminLibId := adminDetails.LibID
	readerLibId := readerDetails.LibID

	if adminLibId != readerLibId {
		c.JSON(http.StatusNotFound, gin.H{
			"Message": "Id not same so no registry to show for the reader by the admin",
		})
		return
	}

	var info []models.IssueRegistry

	if err := initializers.DB.Where("reader_id=?", readerID).Find(&info).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{
			"Message": "Coulnt find any issue registry",
			"Error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"info": info,
	})

}
