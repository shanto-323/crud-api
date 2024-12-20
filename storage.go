package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	CreateAccount(*Member) error
	UpdateAccount(*Member, int) (*Member, error)
	GetAllData() ([]*Member, error)
	GetAccountById(int) (*Member, error)
	DeleteAccount(int) error
}

type DB struct {
	db *sql.DB
}

func CreateStorage() (*DB, error) {
	dsn := "user1:12345678@tcp(localhost:3306)/mydb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (p *DB) init() error {
	return p.createAccountTable()
}

func (p *DB) createAccountTable() error {
	quary := `  
    CREATE TABLE IF NOT EXISTS membersAccount (
      id SERIAL PRIMARY KEY,
      name VARCHAR(30),
      active BOOLEAN,
      subscription_type VARCHAR(20),
      join_date VARCHAR(50) 
    )
  `

	_, err := p.db.Exec(quary)
	return err
}

func (db *DB) CreateAccount(member *Member) error {
	query := `
      INSERT INTO members (name, active, subscription_type, join_date) 
      VALUES (?, ?, ?, ?) 
  `
	_, err := db.db.Exec(query, member.Name, member.Active, member.Subscription_type, member.Join_date)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetAllData() ([]*Member, error) {
	query := `
        SELECT * FROM members
  `
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*Member{}
	for rows.Next() {
		account, err := getAccount(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, account)
	}
	return members, nil
}

func (db *DB) GetAccountById(id int) (*Member, error) {
	query := `
        SELECT * FROM members
        WHERE id = ?
  `
	row := db.db.QueryRow(query, id)

	member := &Member{}
	err := row.Scan(&member.Id, &member.Name, &member.Active, &member.Subscription_type, &member.Join_date)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (db *DB) UpdateAccount(member *Member, id int) (*Member, error) {
	query := `
      UPDATE members SET
      name = ?, active = ?, subscription_type = ?, join_date = ?
      WHERE id = ?
  `
	_, err := db.db.Exec(query, member.Name, member.Active, member.Subscription_type, member.Join_date, id)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (db *DB) DeleteAccount(id int) error {
	query := `
      DELETE FROM members
      WHERE id = ?
  `
	_, err := db.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func getAccount(row *sql.Rows) (*Member, error) {
	member := &Member{}
	err := row.Scan(
		&member.Id,
		&member.Name,
		&member.Active,
		&member.Subscription_type,
		&member.Join_date)

	if err != nil {
		return nil, err
	}
	return member, nil
}
