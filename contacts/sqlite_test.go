package contacts

import (
	"testing"
)

var (
	isOpen     = false
	contactsDB *SqliteContactsDB
	testPerson = Person{
		Id:         IntPointer(1),
		Name:       "name",
		Surname:    "surname",
		Nickname:   "nick",
		Email:      "email@domain.com",
		Mobile:     "123123123",
		TelegramID: Int64Pointer(123123123),
	}
)

const (
	testDBPath = "./test_db.db"
)

func TestOpen(t *testing.T) {
	contactsDB = NewSqliteContactsDB(testDBPath)
	err := contactsDB.Open()
	if err != nil {
		isOpen = false
		t.Error(err)
	}
	isOpen = true
}

func TestSQLite(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}
	err := contactsDB.CreateContactTable()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}
	err := contactsDB.Insert(testPerson)
	if err != nil {
		t.Error(err)
	}
}

func TestGetByName(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}

	p, err := contactsDB.GetByName("name")
	if err != nil {
		t.Error(err)
	}
	if p.Name != testPerson.Name {
		t.Errorf("expected %s, got %s", testPerson.Name, p.Name)
	}
}

func TestGetByNameSurname(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}

	p, err := contactsDB.GetByNameSurname("name", "surname")
	if err != nil {
		t.Error(err)
	}
	if p.Name != testPerson.Name || p.Surname != testPerson.Surname {
		t.Errorf("expected %s %s, got %s %s", testPerson.Name, testPerson.Surname, p.Name, p.Surname)
	}
}

func TestGetById(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}

	p, err := contactsDB.GetById(1)
	if err != nil {
		t.Error(err)
	}
	if p.Name != testPerson.Name || p.Surname != testPerson.Surname {
		t.Errorf("expected %s %s, got %s %s", testPerson.Name, testPerson.Surname, p.Name, p.Surname)
	}
}

func TestGetByPerson(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}

	// Name and Surname
	person := Person{
		Name:    "name",
		Surname: "surname",
	}

	p, err := contactsDB.GetByPerson(person)
	if err != nil {
		t.Error(err)
	}
	if p.Name != person.Name || p.Surname != person.Surname {
		t.Errorf("expected %s %s, got %s %s", testPerson.Name, testPerson.Surname, p.Name, p.Surname)
	}
	t.Logf("Person: %v", p)

	// ID
	person = Person{
		Id: IntPointer(1),
	}

	p, err = contactsDB.GetByPerson(person)
	if err != nil {
		t.Error(err)
	}
	if *(p.Id) != *(person.Id) {
		t.Errorf("expected %d, got %d", *(testPerson.Id), *(p.Id))
	}
	t.Logf("Person: %v", p)
}

func DeleteById(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}
	err := contactsDB.DeleteById(1)
	if err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	if !isOpen {
		contactsDB = NewSqliteContactsDB(testDBPath)
	}
	err := contactsDB.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestDestroy(t *testing.T) {
	err := DestroyDatabase(testDBPath)
	if err != nil {
		t.Error(err)
	}
}
