package models_test

import (
	"testing"

	"time"

	"github.com/juang77/GoMSArch/profile-service/app/models"
)

func getSimpleUser() models.UserCreate {
	user := models.UserCreate{}
	user.ID = 1
	user.Username = "NICOLAS GOMEZ PARRA"
	user.Email = "nicogopa2018@hotmail.com"
	user.Password = "12345678"
	user.Mobile = "0313166937870"
	return user
}

func TestGetUsername(t *testing.T) {
	user := models.UserCreate{}
	user.Username = "user"

	expected := "user"
	actual := user.GetUsername()

	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestPrint(t *testing.T) {
	now := time.Now()

	user := getSimpleUser()
	user.CreatedAt = now

	expected := "user (1) - " + now.String()
	actual := user.Print()

	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestValidate(t *testing.T) {
	user := getSimpleUser()

	err := user.Validate()
	if err != nil {
		t.Fatalf("Expected too pass but got %s", err.Error())
	}

	err = user.ValidatePassword()
	if err != nil {
		t.Fatalf("Expected too pass but got %s", err.Error())
	}
}

func TestTooShortUsername(t *testing.T) {
	user := models.UserCreate{}
	user.Username = ""
	user.Password = "pass"
	err := user.Validate()

	if err == nil {
		t.Fatalf("Expected too receive an error.")
	} else {
		if err.Error() != models.ErrUsernameTooShort.Error() {
			t.Fatalf("Expected too receive a ErrUsernameTooShort error but got %v.", err.Error())
		}
	}
}

func TestTooShortEMail(t *testing.T) {
	user := models.UserCreate{}
	user.Username = "user"
	user.Email = ""
	err := user.Validate()

	if err == nil {
		t.Fatalf("Expected too receive an error.")
	} else {
		if err.Error() != models.ErrEmailTooShort.Error() {
			t.Fatalf("Expected too receive a ErrUsernameTooShort error but got %v.", err.Error())
		}
	}
}

func TestTooShortPassword(t *testing.T) {
	user := models.UserCreate{}
	user.Username = "user"
	user.Password = ""
	err := user.ValidatePassword()

	if err == nil {
		t.Fatalf("Expected too receive an error.")
	} else {
		if err.Error() != models.ErrPasswordTooShort.Error() {
			t.Fatalf("Expected too receive a ErrUsernameTooShort error but got %v.", err.Error())
		}
	}
}
