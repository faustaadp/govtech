package handlers

import (
	"net/http"
	"govtech/database"
	"govtech/models"
	"io/ioutil"
	"encoding/json"
	"time"
)

func HandleProducts(w http.ResponseWriter, r *http.Request) {
    db := database.DbConn()

	queryParams := r.URL.Query()
	sku := queryParams.Get("sku")
	title := queryParams.Get("title")
	category := queryParams.Get("category")
	etalase := queryParams.Get("etalase")
	sortBy := queryParams.Get("sort")

	switch r.Method {
	case http.MethodGet:
		var query, querySort string
		var args []interface{}

		switch sortBy {
		case "newest":
			querySort = " ORDER BY created_at DESC"
		case "oldest":
			querySort = " ORDER BY created_at ASC"
		case "highest-rated":
			querySort = " ORDER BY average_rating DESC"
		case "lowest-rated":
			querySort = " ORDER BY average_rating ASC"
		default:
			querySort = ""
		}

		if sku != "" {
			query += " AND sku = ?"
			args = append(args, sku)
		}
		if title != "" {
			query += " AND title = ?"
			args = append(args, title)
		}
		if category != "" {
			query += " AND category = ?"
			args = append(args, category)
		}
		if etalase != "" {
			query += " AND etalase = ?"
			args = append(args, etalase)
		}

		if query != "" {
			query = "WHERE" + query[4:]
		}

		result, err := db.Query("SELECT id, sku, title, description, category, etalase, image, weight, price, created_at, average_rating FROM Products " + query + querySort, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		
		res := []models.Product{}
		for result.Next() {
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

			product := models.Product{
				Id: id,
				Sku:   sku,
				Title: title,
				Description: description,
				Category: category,
				Etalase: etalase,
				Image: image,
				Weight: weight,
				Price: price,
				AverageRating: average_rating,
				CreatedAt: created_at,
			}
			res = append(res, product)
		}

		jsonRes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonRes)

	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var productData models.Product
		err = json.Unmarshal(body, &productData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		insForm, err := db.Prepare("INSERT INTO Products(sku, title, description, category, etalase, image, weight, price) VALUES(?,?,?,?,?,?,?,?)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		_, err = insForm.Exec(productData.Sku, productData.Title, productData.Description, productData.Category, productData.Etalase, productData.Image, productData.Weight, productData.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/products", http.StatusSeeOther)

	default:
		http.Error(w, "Wrong method", http.StatusInternalServerError)
	}
}