package handlers

import (
	"net/http"
	"govtech/database"
	"govtech/models"
	"encoding/json"
	"io/ioutil"
	"strings"
)

func HandleProductReviews(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()

    productId := strings.TrimPrefix(r.URL.Path, "/products/reviews/")
	var sum_rating, count_rating int

	result, err := db.Query("SELECT sum_rating, count_rating FROM Products WHERE id=?", productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if result.Next() {
		err = result.Scan(&sum_rating, &count_rating)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var reviewData models.Review
		err = json.Unmarshal(body, &reviewData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		db = database.DbConn()
		insForm, err := db.Prepare("INSERT INTO Reviews(product_id, rating, comment) VALUES(?,?,?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
		_, err = insForm.Exec(productId, reviewData.Rating, reviewData.Comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var newSumRating, newCountRating int
		var newAverageRating float32
		
		newSumRating = sum_rating + reviewData.Rating
		newCountRating = count_rating + 1
		newAverageRating = float32(newSumRating) / float32(newCountRating)

		db = database.DbConn()
		insForm, err = db.Prepare("UPDATE Products SET sum_rating=?, count_rating=?, average_rating=? WHERE id=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = insForm.Exec(newSumRating, newCountRating, newAverageRating, productId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
		http.Redirect(w, r, "/products/reviews/" + productId, http.StatusSeeOther)

	case http.MethodGet:
		db := database.DbConn()
		result, err := db.Query("SELECT id, rating, comment FROM Reviews WHERE product_id=?", productId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
		res := []models.Review{}
		for result.Next() {
			var id, rating int
			var comment string
			err = result.Scan(&id, &rating, &comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			review := models.Review{
				Id: id,
				Comment: comment,
				Rating: rating,
			}
			res = append(res, review)
		}

		jsonRes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonRes)
	default:
		http.Error(w, "Wrong method", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/products", http.StatusSeeOther)
}