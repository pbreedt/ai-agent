package contacts

import (
	"testing"
)

var (
	isOpen     = false
	contactsDB *ContactsDB
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
	contactsDB = NewContactsDB(testDBPath)
	err := contactsDB.Open()
	if err != nil {
		isOpen = false
		t.Error(err)
	}
	isOpen = true
}

func TestSQLite(t *testing.T) {
	if !isOpen {
		contactsDB = NewContactsDB(testDBPath)
	}
	err := contactsDB.CreateContactTable()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	if !isOpen {
		contactsDB = NewContactsDB(testDBPath)
	}
	err := contactsDB.Insert(testPerson)
	if err != nil {
		t.Error(err)
	}
}

func TestGetByName(t *testing.T) {
	if !isOpen {
		contactsDB = NewContactsDB(testDBPath)
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
		contactsDB = NewContactsDB(testDBPath)
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
		contactsDB = NewContactsDB(testDBPath)
	}

	p, err := contactsDB.GetById(1)
	if err != nil {
		t.Error(err)
	}
	if p.Name != testPerson.Name || p.Surname != testPerson.Surname {
		t.Errorf("expected %s %s, got %s %s", testPerson.Name, testPerson.Surname, p.Name, p.Surname)
	}
}

func DeleteById(t *testing.T) {
	if !isOpen {
		contactsDB = NewContactsDB(testDBPath)
	}
	err := contactsDB.DeleteById(1)
	if err != nil {
		t.Error(err)
	}
}

func TestClose(t *testing.T) {
	if !isOpen {
		contactsDB = NewContactsDB(testDBPath)
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
