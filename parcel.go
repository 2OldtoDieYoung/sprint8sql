package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	//добавляем строку
	res, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	// возвращаем идентификатор последней добавленной записи
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("посылка не найдена")
		}
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	// менять адрес можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("адрес может быть изменен только у посылки со статусом 'registered'")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {

	// удалять строку можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("посылку можно удалить только со статусом'registered'")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number",
		sql.Named("number", number))
	if err != nil {
		return err
	}

	return nil
}
