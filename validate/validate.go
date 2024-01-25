package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type RuleFunc func() RuleSet

type ValidateFunc func(RuleSet) bool

type RuleSet struct {
	Name         string
	RuleValue    any
	FieldValue   any
	FieldName    any
	MessageFunc  func(RuleSet) string
	ValidateFunc func(RuleSet) bool
}

type Fields map[string][]RuleSet

type Messages map[string]string

func Required() RuleSet {
	return RuleSet{
		Name: "required",
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("%s is a required field", set.FieldName)
		},
		ValidateFunc: func(rule RuleSet) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) > 0
		},
	}
}

func Message(msg string) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "message",
			RuleValue: msg,
		}
	}
}

func Email() RuleSet {
	return RuleSet{
		Name: "email",
		MessageFunc: func(set RuleSet) string {
			return "email address is invalid"
		},
		ValidateFunc: func(set RuleSet) bool {
			email, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return emailRegex.MatchString(email)
		},
	}
}

func Max(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "max",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) <= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be maximum %d characters long", set.FieldName, n)
			},
		}
	}
}

func Min(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "min",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) >= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be at least %d characters long", set.FieldName, n)
			},
		}
	}
}

func Rules(rules ...RuleFunc) []RuleSet {
	ruleSets := make([]RuleSet, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]()
	}
	return ruleSets
}

type Validator struct {
	data   any
	fields Fields
}

func New(data any, fields Fields) *Validator {
	return &Validator{
		fields: fields,
		data:   data,
	}
}

func (v *Validator) Validate(target any) bool {
	ok := true
	for fieldName, ruleSets := range v.fields {
		// reflect panics on un-exported variables.
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldValue := getFieldValueByName(v.data, fieldName)
		for _, set := range ruleSets {
			set.FieldValue = fieldValue
			set.FieldName = fieldName
			if set.Name == "message" {
				setErrorMessage(target, fieldName, set.RuleValue.(string))
				continue
			}
			if !set.ValidateFunc(set) {
				msg := set.MessageFunc(set)
				setErrorMessage(target, fieldName, msg)
				ok = false
			}
		}
	}
	return ok
}

func setErrorMessage(v any, fieldName string, msg string) {
	if v == nil {
		return
	}
	switch t := v.(type) {
	case map[string]string:
		t[fieldName] = msg
	default:
		structVal := reflect.ValueOf(v)
		structVal = structVal.Elem()
		field := structVal.FieldByName(fieldName)
		field.Set(reflect.ValueOf(msg))
	}
}

func getFieldValueByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil
	}
	return fieldVal.Interface()
}
