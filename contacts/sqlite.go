package contacts

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteContactsDB struct {
	path string
	db   *sql.DB
}

func NewSqliteContactsDB(path string) *SqliteContactsDB {
	return &SqliteContactsDB{path: path}
}

func (c *SqliteContactsDB) Open() error {
	db, err := sql.Open("sqlite3", c.path)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func (c *SqliteContactsDB) Close() error {
	return c.db.Close()
}

func (c *SqliteContactsDB) CreateContactTable() error {
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

func (c *SqliteContactsDB) GetAll() ([]Person, error) {
	var contacts []Person
	rows, err := c.db.Query("select * from contact")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var contact Person
		err := rows.Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (c *SqliteContactsDB) GetByPerson(p Person) (Person, error) {
	var found Person
	var err error
	if p.Id != nil {
		found, err = c.GetById(*p.Id)
		if found != (Person{}) || err != nil {
			return found, err
		}
	}
	if p.Name != "" && p.Surname != "" {
		return c.GetByNameSurname(p.Name, p.Surname)
	}
	return c.GetByName(p.Name)
}

func (c *SqliteContactsDB) GetById(id int) (Person, error) {
	var contact Person
	err := c.db.QueryRow("select * from contact where id = ?", id).Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
	return contact, err
}

func (c *SqliteContactsDB) GetByName(name string) (Person, error) {
	var contact Person
	err := c.db.QueryRow("select * from contact where name = ?", name).Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
	return contact, err
}

func (c *SqliteContactsDB) GetByFullname(fullname string) (Person, error) {
	var contact Person

	split := strings.Split(fullname, " ")
	if len(split) != 2 {
		return contact, fmt.Errorf("invalid fullname: %s", fullname)
	}
	name := split[0]
	surname := split[1]

	return c.GetByNameSurname(name, surname)
}

func (c *SqliteContactsDB) GetByNameSurname(name string, surname string) (Person, error) {
	var contact Person

	err := c.db.QueryRow("select * from contact where name = ? and surname = ?", name, surname).Scan(&contact.Id, &contact.Name, &contact.Surname, &contact.Nickname, &contact.Email, &contact.Mobile, &contact.TelegramID)
	return contact, err
}

func (s *SqliteContactsDB) Insert(person Person) error {
	_, err := s.db.Exec("insert into contact(name, surname, nickname, email, mobile, telegram_id) values(?, ?, ?, ?, ?, ?)", person.Name, person.Surname, person.Nickname, person.Email, person.Mobile, person.TelegramID)
	return err
}

// TODO: update is vulnerable to sql injection - need to use prepared statements or parameterized queries
func (s *SqliteContactsDB) Update(person Person) error {
	sql := "update contact"

	if person.Nickname != "" {
		sql += fmt.Sprintf(" set nickname = '%s'", person.Nickname)
	}
	if person.Email != "" {
		sql += fmt.Sprintf(" set email = '%s'", person.Email)
	}
	if person.Mobile != "" {
		sql += fmt.Sprintf(" set mobile = '%s'", person.Mobile)
	}
	if person.TelegramID != nil {
		sql += fmt.Sprintf(" set telegram_id = %d", *person.TelegramID)
	}

	if person.Id != nil {
		if person.Name != "" {
			sql += fmt.Sprintf(" set name = '%s'", person.Name)
		}
		if person.Surname != "" {
			sql += fmt.Sprintf(" set surname = '%s'", person.Surname)
		}
		sql += fmt.Sprintf(" where id = %d", *person.Id)
	} else if person.Name != "" && person.Surname != "" {
		sql += fmt.Sprintf(" where name = '%s' and surname = '%s'", person.Name, person.Surname)
	}

	log.Println("Executing SQL:", sql)
	_, err := s.db.Exec(sql)
	return err
}

func (s *SqliteContactsDB) DeleteById(id int) error {
	_, err := s.db.Exec("delete from contact where id = ?", id)
	return err
}

func (s *SqliteContactsDB) DeleteByNameSurname(name string, surname string) error {
	_, err := s.db.Exec("delete from contact where name = ? and surname = ?", name, surname)
	return err
}

func (s *SqliteContactsDB) Clear() error {
	_, err := s.db.Exec("delete from foo")
	return err
}

func (s *SqliteContactsDB) Count() (int, error) {
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
