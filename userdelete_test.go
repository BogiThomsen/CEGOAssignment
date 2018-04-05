
package main

import (
	"testing"	
)

//Testing User constructor
func TestNewUser(t *testing.T){
	id := "1"
	firstName := "firstName"
	lastName := "lastName"
	email := "email"

	user := NewUser(id, firstName, lastName, email)

	if user.ID != id {
		t.Errorf("Expected %s but got %s\n", id, user.ID)
	}
	if user.firstName != firstName {
		t.Errorf("Expected %s but got %s\n", firstName, user.firstName)
	}
	if user.lastName != lastName {
		t.Errorf("Expected %s but got %s\n", lastName, user.lastName)
	}
	if user.email != email {
		t.Errorf("Expected %s but got %s\n", email, user.email)
	}
}