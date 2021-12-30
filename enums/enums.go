package enums

import (
	"fmt"
	"strings"

	"github.com/johannes-kuhfuss/services_utils/api_error"
)

// Usage: var AccountTypes = Enum{[]EnumItem{{0, "Basic"}, {1, "Advanced"}}}

type EnumItem struct {
	Idx int32
	Val string
}
type Enum struct {
	Items []EnumItem
}

func (e *Enum) AsValue(i int32) (string, api_error.ApiErr) {
	item, err := e.ItemByIndex(i)
	if err != nil {
		return "", err
	}
	return item.Val, nil
}

func (e *Enum) AsIndex(v string) (int32, api_error.ApiErr) {
	item, err := e.ItemByValue(v)
	if err != nil {
		return 0, err
	}
	return item.Idx, nil
}

func (e *Enum) Values() []string {
	var names []string
	for _, item := range e.Items {
		names = append(names, item.Val)
	}
	return names
}

func (e *Enum) AsMap() map[int32]string {
	m := make(map[int32]string)
	for _, item := range e.Items {
		m[item.Idx] = item.Val
	}
	return m
}

func (e *Enum) FromMap(m map[int32]string) {
	var eItem EnumItem
	for index, item := range m {
		eItem.Idx = index
		eItem.Val = item
		e.Items = append(e.Items, eItem)
	}
}

func (e *Enum) ItemByValue(v string) (*EnumItem, api_error.ApiErr) {
	for _, item := range e.Items {
		if strings.EqualFold(v, item.Val) {
			return &item, nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("No item with value %v found", v))
}

func (e *Enum) ItemByIndex(i int32) (*EnumItem, api_error.ApiErr) {
	for _, item := range e.Items {
		if item.Idx == i {
			return &item, nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("No item with index %v found", i))
}
