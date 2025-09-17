package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQLドライバをインポート
)

// (Item構造体の定義はSQLite版と同じ)
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	// --- MySQL接続設定 ---
	// DSN (Data Source Name)
	// docker-compose.ymlで設定したユーザー、パスワード、DB名、サービス名(db)を指定
	dsn := "user:password@tcp(db:3306)/testdb?parseTime=true"

	var err error
	// MySQLコンテナが起動完了するまで待つため、接続をリトライする
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to MySQL!")
				break
			}
		}
		log.Printf("Could not connect to MySQL, retrying in 5 seconds... (%s)", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to MySQL after retries: %s", err)
	}
	defer db.Close()
	
	createTable() // テーブル作成と初期データ投入

	// (APIエンドポイントの設定はSQLite版と同じ)
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/items/", itemHandler)
	log.Println("API server with MySQL started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	createTableSQL := `CREATE TABLE IF NOT EXISTS items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	// (初期データ投入のロジックはSQLite版と同じ)
	var count int
	db.QueryRow("SELECT COUNT(*) FROM items").Scan(&count)
	if count == 0 {
		log.Println("Inserting initial data...")
		db.Exec(`INSERT INTO items (name) VALUES (?), (?), (?)`, "高性能マウス", "メカニカルキーボード", "4Kモニター")
	}
}

// (itemsHandler, itemHandler, 各API処理(getItems, getItemByID, createItem)のコードは
//  SQLite版と全く同じなので、ここでは省略します。そのままコピーしてください)
func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItems(w, r)
	case http.MethodPost:
		createItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func itemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/items/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	getItemByID(w, r, id)
}
func getItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
func getItemByID(w http.ResponseWriter, r *http.Request, id int) {
	var item Item
	err := db.QueryRow("SELECT id, name FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO items (name) VALUES (?)", item.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lastID, _ := result.LastInsertId()
	item.ID = int(lastID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}