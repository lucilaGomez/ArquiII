package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB wrapper para base de datos MySQL
type DB struct {
	conn *sql.DB
}

// Connect establece conexión con MySQL
func Connect(uri string) (*DB, error) {
	// Conectar a MySQL
	conn, err := sql.Open("mysql", uri+"?parseTime=true&loc=Local")
	if err != nil {
		return nil, fmt.Errorf("error conectando a MySQL: %v", err)
	}

	// Configurar pool de conexiones
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Verificar conexión
	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("error haciendo ping a MySQL: %v", err)
	}

	log.Println("✅ Conectado a MySQL exitosamente")
	
	return &DB{conn: conn}, nil
}

// GetDB obtiene la conexión de base de datos
func (db *DB) GetDB() *sql.DB {
	return db.conn
}

// Close cierra la conexión
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Ping verifica que la conexión esté activa
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// Begin inicia una transacción
func (db *DB) Begin() (*sql.Tx, error) {
	return db.conn.Begin()
}

// Query ejecuta una consulta que devuelve filas
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

// QueryRow ejecuta una consulta que devuelve máximo una fila
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

// Exec ejecuta una consulta que no devuelve filas
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}