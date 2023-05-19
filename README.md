## Govtech Task

This is a REST API application for managing products and reviews.

## Configuration

Before running the application, make sure to configure your MySQL database settings. Open the `database/database.go` file and update the MySQL configuration details.

## Build

To build the Go application, run the following command:

```bash
go build
```

## Run

After building the application, you can run it using the following command:

```bash
./govtech
```

## API Endpoints

- Get all products: GET /products
- Create a new product: POST /products
- Get a product by ID: GET /products/{id}
- Update a product by ID: PUT /products/{id}
- Delete a product by ID: DELETE /products/{id}
- Search products by SKU: GET /products?sku={sku}
- Search products by title: GET /products?title={title}
- Search products by category: GET /products?category={category}
- Search products by etalase: GET /products?etalase={etalase}
- Sort products by date (newest): GET /products?sort=newest
- Sort products by date (oldest): GET /products?sort=oldest
- Sort products by average higest rating: GET /products?sort=highest-rated
- Sort products by average lowest rating: GET /products?sort=lowest-rated
- Add a review for a product: POST /products/reviews/{id}
- Get reviews for a product: GET /products/reviews/{id}
