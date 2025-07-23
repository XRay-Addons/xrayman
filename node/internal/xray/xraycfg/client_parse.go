package xraycfg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func extractNameField(cfg string) (string, error) {
	users := getUsers(gjson.Get(cfg, `outbounds`))
	userIDs, err := extractFields(users, "email")
	if err != nil {
		return "", fmt.Errorf("extract Name field: %w", err)
	}
	userID, err := getSingleValue(userIDs)
	if err != nil {
		return "", fmt.Errorf("extract Name field: %w", err)
	}
	return userID, nil
}

func extractVlessUUIDField(cfg string) (string, error) {
	users := getUsers(gjson.Get(cfg, `outbounds.#(protocol=="vless")#`))
	userIDs, err := extractFields(users, "id")
	if err != nil {
		return "", fmt.Errorf("extract VlessUUID field: %w", err)
	}
	userID, err := getSingleValue(userIDs)
	if err != nil {
		return "", fmt.Errorf("extract VlessUUID field: %w", err)
	}
	return userID, nil
}

// get outbounds users
func getUsers(outs gjson.Result) []gjson.Result {
	users := make([]gjson.Result, 0)
	outs.ForEach(func(_, o gjson.Result) bool {
		o.Get("settings.vnext").ForEach(func(_, v gjson.Result) bool {
			v.Get("users").ForEach(func(_, u gjson.Result) bool {
				users = append(users, u)
				return true
			})
			return true
		})
		return true
	})
	return users
}

// extract fields by names
func extractFields(items []gjson.Result, name string) ([]string, error) {
	uniqueFields := make(map[string]struct{}, 0)
	var errs []error
	for _, item := range items {
		val := item.Get(name)
		if !val.Exists() {
			continue
		}
		field, err := extractTemplateVar(val.String())
		if err != nil {
			errs = append(errs, err)
			continue
		}
		uniqueFields[field] = struct{}{}
	}

	// return all errors on error
	if len(errs) > 0 {
		return nil, fmt.Errorf("extract vless uuid field: %w", errors.Join(errs...))
	}

	// get fields list
	fields := make([]string, 0, len(uniqueFields))
	for f, _ := range uniqueFields {
		fields = append(fields, f)
	}

	return fields, nil
}

func extractTemplateVar(s string) (string, error) {
	// trim spaces
	templateVar := strings.TrimSpace(s)
	// trim "{{", "}}"
	if !strings.HasPrefix(templateVar, "{{") || !strings.HasSuffix(templateVar, "}}") {
		return "", fmt.Errorf("invalid template format: %s", s)
	}
	templateVar = templateVar[2 : len(templateVar)-2]
	// trim spaces again
	templateVar = strings.TrimSpace(templateVar)
	// trim "."
	if !strings.HasPrefix(templateVar, ".") {
		return "", fmt.Errorf("template variable should start with '.'")
	}
	templateVar = templateVar[1:]
	if templateVar == "" {
		return "", fmt.Errorf("empty variable name")
	}
	return templateVar, nil
}

func getSingleValue(values []string) (string, error) {
	if len(values) > 1 {
		return "", fmt.Errorf("multiple name field templates found: %v", values)
	}
	for _, value := range values {
		return value, nil
	}
	return "", nil
}
