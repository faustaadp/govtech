package database

import (
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

func DbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbName := "govtech" 
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/")
    if err != nil {
        panic(err.Error())
    }

    _, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
    if err != nil {
        panic(err.Error())
    }

    _, err = db.Exec("USE " + dbName)
    if err != nil {
        panic(err.Error())
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS Products (
        id INT AUTO_INCREMENT PRIMARY KEY,
        sku VARCHAR(255) UNIQUE,
        title VARCHAR(255),
        description TEXT,
        category VARCHAR(255),
        etalase VARCHAR(255),
        image TEXT,
        weight FLOAT,
        price FLOAT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        sum_rating INT DEFAULT 0,
        count_rating INT DEFAULT 0,
        average_rating FLOAT DEFAULT 0.0
    )`)
    if err != nil {
        panic(err.Error())
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS Reviews (
        id INT AUTO_INCREMENT PRIMARY KEY,
        product_id INT NOT NULL,
        rating INT NOT NULL,
        comment TEXT,
        FOREIGN KEY (product_id) REFERENCES Products (id) ON DELETE CASCADE
    )`)
    if err != nil {
        panic(err.Error())
    }

    return db
}
