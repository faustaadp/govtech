package handlers

import (
	"net/http"
	"govtech/database"
	"govtech/models"
	"io/ioutil"
	"encoding/json"
	"strings"
	"time"
)

func HandleProduct(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()

	id := strings.TrimPrefix(r.URL.Path, "/products/")

	switch r.Method {
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var updateData models.Product
		err = json.Unmarshal(body, &updateData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		insForm, err := db.Prepare("UPDATE Products SET sku=?, title=?, description=?, category=?, etalase=?, image=?, weight=?, average_rating=? WHERE id=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := insForm.Exec(updateData.Sku, updateData.Title, updateData.Description, updateData.Category, updateData.Etalase, updateData.Image, updateData.Weight, updateData.Price, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Product not found", http.StatusInternalServerError)
		}

	case http.MethodDelete:
		insForm, err := db.Prepare("DELETE FROM Products WHERE id=?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := insForm.Exec(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Product not found", http.StatusInternalServerError)
		}
		
	case http.MethodGet:
		result, err := db.Query("SELECT id, sku, title, description, category, etalase, image, weight, price, created_at, average_rating FROM Products WHERE id=?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		product := models.Product{}
		if result.Next() {
			var id int
			var sku, title, description, category, etalase, image, created_at_string string
			var weight, price, average_rating float32
			err = result.Scan(&id, &sku, &title, &description, &category, &etalase, &image, &weight, &price, &created_at_string, &average_rating)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			created_at, err := time.Parse("2006-01-02 15:04:05", created_at_string)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			product.Id = id
			product.Sku = sku
			product.Title = title
			product.Description = description
			product.Category = category
			product.Etalase = etalase
			product.Image = image
			product.Weight = weight
			product.Price = price
			product.CreatedAt = created_at
			product.AverageRating = average_rating
		} else {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		jsonRes, err := json.Marshal(product)
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