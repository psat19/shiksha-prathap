package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var PhoneRX = regexp.MustCompile(`^[0-9]{10}$`)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) IsAgeValid() int {
	value := f.Get("age")
	age, err := strconv.Atoi(value)

	if err != nil {
		f.Errors.Add("age", "Age should be a valid number")
	}

	if age <= 0 {
		f.Errors.Add("age", "Age should be a positive number")
	}
	return age
}

func (f *Form) Has(field string) bool {
	value := f.Get(field)
	return strings.TrimSpace(value) != ""
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) MustHaveLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) == d {
		f.Errors.Add(field, fmt.Sprintf("This field should be of %d characters)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) && field == "phone" {
		f.Errors.Add(field, "Phone field should have exactly 10 digits")
	}

	if !pattern.MatchString(value) && field == "email" {
		f.Errors.Add(field, "Invalid email address")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
