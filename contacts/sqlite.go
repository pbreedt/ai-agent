package contacts

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type ContactsDB struct {
	path string
	db   *sql.DB
}

func NewContactsDB(path string) *ContactsDB {
	return &ContactsDB{path: path}
}

func (c *ContactsDB) Open() error {
	db, err := sql.Open("sqlite3", c.path)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *ContactsDB) Close() error {
	return c.db.Close()
}

func (c *ContactsDB) CreateContactTable() error {
	sqlStmt := `
	create table if not exists contact (
		id integer not null primary key, 
		name text,
		surname text,
		nickname text,
		email text,
		mobile text,
		telegram_id integer
	);
	`
	_, err := c.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

func DestroyDatabase(path string) error {
	return os.Remove(path)
}

func (c *ContactsDB) GetByPerson(p Person) (Person, error) {
	if p.Id != nil {
		return c.GetById(*p.Id)
	}
	if p.Name != "" && p.Surname != "" {
		return c.GetByNameSurname(p.Name, p.Surname)
	}
	return c.GetByName(p.Name)
}

func (c *ContactsDB) GetById(id int) (Person, error) {
	var contact Person
	err := c.db.QueryRow("select * from contact where id = ?", id).Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
	return contact, err
}

func (c *ContactsDB) GetByName(name string) (Person, error) {
	var contact Person
	err := c.db.QueryRow("select * from contact where name = ?", name).Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
	return contact, err
}

func (c *ContactsDB) GetByFullname(fullname string) (Person, error) {
	var contact Person

	split := strings.Split(fullname, " ")
	if len(split) != 2 {
		return contact, fmt.Errorf("invalid fullname: %s", fullname)
	}
	name := split[0]
	surname := split[1]

	return c.GetByNameSurname(name, surname)
}

func (c *ContactsDB) GetByNameSurname(name string, surname string) (Person, error) {
	var contact Person

	err := c.db.QueryRow("select * from contact where name = ? and surname = ?", name, surname).Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
	return contact, err
}

func (s *ContactsDB) Insert(person Person) error {
	_, err := s.db.Exec("insert into contact(name, surname, nickname, email, mobile, telegram_id) values(?, ?, ?, ?, ?, ?)", person.Name, person.Surname, person.Nickname, person.Email, person.Mobile, person.TelegramID)
	return err
}

func (s *ContactsDB) DeleteById(id int) error {
	_, err := s.db.Exec("delete from contact where id = ?", id)
	return err
}

func (s *ContactsDB) DeleteByNameSurname(name string, surname string) error {
	_, err := s.db.Exec("delete from contact where name = ? and surname = ?", name, surname)
	return err
}

func (s *ContactsDB) Clear() error {
	_, err := s.db.Exec("delete from foo")
	return err
}

func (s *ContactsDB) Count() (int, error) {
	var count int
	err := s.db.QueryRow("select count(*) from foo").Scan(&count)
	return count, err
}

func test() {
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("こんにちは世界%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
