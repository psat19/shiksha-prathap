package main

import (
	"net/url"
	"testing"

	models "github.com/psat/shiksha-prathap/pkg/models"
)

type fakeUserStore struct {
	insertFn func(email, password string) (int, error)
}

func (f fakeUserStore) Insert(email, password string) (int, error) {
	return f.insertFn(email, password)
}

func TestProcessSignup_InvalidMissingFields(t *testing.T) {
	post := url.Values{}
	post.Add("email", "test@example.com")
	// missing password

	fake := fakeUserStore{
		insertFn: func(email, password string) (int, error) {
			t.Fatal("Insert should not be called for invalid form")
			return 0, nil
		},
	}

	form, id, err := processSignup(post, fake)

	// basically insert method would not get called if the form is invalid
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 0 {
		t.Fatalf("expected id=0 for invalid form, got %d", id)
	}
	if form == nil {
		t.Fatalf("expected non-nil form")
	}
	if form.Valid() {
		t.Fatalf("expected form to be invalid (missing password)")
	}

	if form.Errors.Get("password") == "" {
		t.Fatalf("expected password validation error")
	}
}

func TestProcessSignup_InvalidEmailFormat(t *testing.T) {
	post := url.Values{}
	post.Add("email", "not-an-email")
	post.Add("password", "secret123")

	fake := fakeUserStore{
		insertFn: func(email, password string) (int, error) {
			t.Fatal("Insert should not be called for invalid email format")
			return 0, nil
		},
	}

	form, id, err := processSignup(post, fake)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 0 {
		t.Fatalf("expected id=0 for invalid form, got %d", id)
	}
	if form.Valid() {
		t.Fatalf("expected form to be invalid (bad email)")
	}
	if form.Errors.Get("email") == "" {
		t.Fatalf("expected email validation error")
	}
}

func TestProcessSignup_DuplicateEmail(t *testing.T) {
	post := url.Values{}
	post.Add("email", "dup@example.com")
	post.Add("password", "secret123")

	fake := fakeUserStore{
		insertFn: func(email, password string) (int, error) {
			return 0, models.ErrDuplicateEmail
		},
	}

	form, id, err := processSignup(post, fake)
	if err == nil {
		t.Fatalf("expected duplicate email error, got nil")
	}
	if err != models.ErrDuplicateEmail {
		t.Fatalf("expected ErrDuplicateEmail, got %v", err)
	}
	if id != 0 {
		t.Fatalf("expected id=0 on duplicate, got %d", id)
	}
	if form == nil || !form.Valid() {
		// For duplicate email the form was valid, insertion failed.
		t.Fatalf("expected valid form even when insert fails with duplicate")
	}
}

func TestProcessSignup_Success(t *testing.T) {
	post := url.Values{}
	post.Add("email", "ok@example.com")
	post.Add("password", "secret123")

	fake := fakeUserStore{
		insertFn: func(email, password string) (int, error) {
			return 42, nil
		},
	}

	form, id, err := processSignup(post, fake)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 42 {
		t.Fatalf("expected id=42, got %d", id)
	}
	if form == nil || !form.Valid() {
		t.Fatalf("expected valid form on success")
	}
	// optional: verify form contains the email
	if form.Get("email") != "ok@example.com" {
		t.Fatalf("expected form email preserved, got %s", form.Get("email"))
	}
}
