package models

import (
	"database/sql"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH не установлен в переменных окружения")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) CreateTable() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_ids (
			chat_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS errors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			error_text TEXT,
			chat_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) SaveChatId(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	chatID := update.Message.Chat.ID
	username := update.Message.Chat.UserName
	_, err := d.db.Exec("INSERT INTO chat_ids (chat_id, username) VALUES (?, ?)", chatID, username)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(int64(1528078947), "@GetYoutubeImgBot Новый пользователь: "+update.Message.Contact.FirstName)
	bot.Send(msg)

	return nil
}

func (d *Database) SaveError(update tgbotapi.Update, error_text string) error {
	chatID := update.Message.Chat.ID
	_, err := d.db.Exec("INSERT INTO errors (chat_id, error_text) VALUES (?, ?)", chatID, error_text)
	if err != nil {
		log.Println("Error inserting data into database (errors):", err)
		return err
	}

	return nil
}

func (d *Database) GetUserByID(id int) (string, string, error) {
	var name, email string
	err := d.db.QueryRow("SELECT name, email FROM users WHERE id = ?", id).Scan(&name, &email)
	if err != nil {
		return "", "", err
	}

	return name, email, nil
}
