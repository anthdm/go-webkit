package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode"
)

var defaultMessages = map[string]string{
	"min": "must be at least %d characters",
	"max": "must not exceed %d characters",
}

type RuleFunc func() RuleSet

type RuleSet struct {
	Name  string
	Value any
}

type Fields map[string][]RuleSet

func Required() RuleSet {
	return RuleSet{
		Name: "required",
	}
}

func Email() RuleSet {
	return RuleSet{
		Name: "email",
	}
}

func Max(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:  "max",
			Value: n,
		}
	}
}

func Min(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:  "min",
			Value: n,
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

func Validate(v any, fields Fields) (map[string]string, bool) {
	errors := map[string]string{}
	ok := true
	for fieldName, ruleSets := range fields {
		// reflect panics on un-exported variables.
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldValue := getFieldValueByName(v, fieldName)
		for _, set := range ruleSets {
			if !validate(fieldValue, set) {
				msg := defaultMessages[set.Name]
				errors[fieldName] = fmt.Sprintf(msg, set.Value)
				ok = false
			}
		}
	}
	return errors, ok
}

func validate(value any, ruleSet RuleSet) bool {
	switch ruleSet.Name {
	case "required":
		str, ok := validateString(value)
		if !ok {
			return false
		}
		return len(str) > 0
	case "email":
		email, ok := validateString(value)
		if !ok {
			return false
		}
		return validateEmail(email)
	case "min":
		return validateMinMax(value, ruleSet.Value.(int), true)
	case "max":
		return validateMinMax(value, ruleSet.Value.(int), false)
	}
	return false
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

func validateMinMax(v any, n int, min bool) bool {
	switch t := v.(type) {
	case string:
		if min {
			return len(t) >= n
		}
		return len(t) <= n
	case int:
		if min {
			return t >= n
		}
		return t <= n
	default:
		return false
	}
}

func validateString(v any) (out string, ok bool) {
	out, ok = v.(string)
	return
}

func validateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
