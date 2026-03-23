package params

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type QueryParams struct {
	PageNumber int
	PageSize   int
	Search     string
	OrderBy    string
	OrderDesc  bool
}

func NewQueryParams(c echo.Context) QueryParams {
	filters := make(map[string]string)
	filterKeys := []string{
		"user_id",
	}

	for _, key := range filterKeys {
		if value := c.QueryParam(key); value != "" {
			filters[key] = value
		}
	}

	return QueryParams{
		PageNumber: parseIntWithDefault(c.QueryParam("page"), 1),
		PageSize:   parseIntWithDefault(c.QueryParam("size"), 10),
		Search:     c.QueryParam("search"),
		OrderBy:    c.QueryParam("order_by"),
		OrderDesc:  parseBoolWithDefault(c.QueryParam("order_desc"), false),
	}

}

func parseIntWithDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	var value int
	_, err := fmt.Sscanf(s, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}

func parseBoolWithDefault(s string, defaultValue bool) bool {
	if s == "" {
		return defaultValue
	}
	var value bool
	_, err := fmt.Sscanf(s, "%t", &value)
	if err != nil {
		return defaultValue
	}
	return value
}
