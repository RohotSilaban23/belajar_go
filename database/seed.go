package database

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(db *sql.DB) error {
	nama := "Administrator"
	email := "admin@belajar.com"
	password := "Admin123!"
	role := "admin"
	isActive := true

	// Periksa apakah akun admin sudah ada.
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM users WHERE email = ?",
		email,
	).Scan(&count)

	if err != nil {
		return err
	}

	// Jika sudah ada, jangan membuat akun lagi.
	if count > 0 {
		fmt.Println("Admin sudah tersedia")
		return nil
	}

	// Hash password awal.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	// Masukkan data admin.
	_, err = db.Exec(
		"INSERT INTO users (nama, email, password, role, isActive) VALUES (?, ?, ?, ?, ?)",
		nama,
		email,
		string(hashedPassword),
		role,
		isActive,
	)
	if err != nil {
		return err
	}

	fmt.Println("Admin awal berhasil dibuat")
	return nil
}
