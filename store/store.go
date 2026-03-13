package store

import (
	"database/sql"
	"log"

	"destinyServer/config"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./destiny.db?_journal_mode=WAL")
	if err != nil {
		log.Fatal("open db:", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		open_id    TEXT PRIMARY KEY,
		free_count INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS orders (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		open_id      TEXT NOT NULL,
		out_trade_no TEXT UNIQUE NOT NULL,
		amount_fen   INTEGER NOT NULL,
		status       TEXT NOT NULL DEFAULT 'pending',
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS referrals (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		referrer_id TEXT NOT NULL,
		referee_id  TEXT NOT NULL,
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(referrer_id, referee_id)
	);
	`
	if _, err := DB.Exec(schema); err != nil {
		log.Fatal("init schema:", err)
	}
}

func GetOrCreateUser(openID string) (int, error) {
	var freeCount int
	err := DB.QueryRow("SELECT free_count FROM users WHERE open_id = ?", openID).Scan(&freeCount)
	if err == sql.ErrNoRows {
		_, err = DB.Exec("INSERT INTO users (open_id, free_count) VALUES (?, ?)", openID, config.Cfg.InitFreeUses)
		if err != nil {
			return 0, err
		}
		return config.Cfg.InitFreeUses, nil
	}
	return freeCount, err
}

func UseFreeCount(openID string) (int, error) {
	result, err := DB.Exec(
		"UPDATE users SET free_count = free_count - 1 WHERE open_id = ? AND free_count > 0",
		openID,
	)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return 0, nil
	}
	var remaining int
	DB.QueryRow("SELECT free_count FROM users WHERE open_id = ?", openID).Scan(&remaining)
	return remaining, nil
}

func GetFreeCount(openID string) int {
	var count int
	DB.QueryRow("SELECT free_count FROM users WHERE open_id = ?", openID).Scan(&count)
	return count
}

func AddReferralBonus(referrerID, refereeID string) error {
	if referrerID == "" || referrerID == refereeID {
		return nil
	}

	_, err := DB.Exec(
		"INSERT OR IGNORE INTO referrals (referrer_id, referee_id) VALUES (?, ?)",
		referrerID, refereeID,
	)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"UPDATE users SET free_count = free_count + 1 WHERE open_id = ?",
		referrerID,
	)
	return err
}

func CreateOrder(openID, outTradeNo string, amountFen int) error {
	_, err := DB.Exec(
		"INSERT INTO orders (open_id, out_trade_no, amount_fen, status) VALUES (?, ?, ?, 'pending')",
		openID, outTradeNo, amountFen,
	)
	return err
}

func CompleteOrder(outTradeNo string) error {
	_, err := DB.Exec(
		"UPDATE orders SET status = 'paid' WHERE out_trade_no = ? AND status = 'pending'",
		outTradeNo,
	)
	return err
}
