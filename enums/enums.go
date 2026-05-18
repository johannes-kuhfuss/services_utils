package enums

import (
	"fmt"
	"sort"
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
	names := make([]string, 0, len(e.Items))
	for _, item := range e.Items {
		names = append(names, item.Val)
	}
	return names
}

func (e *Enum) AsMap() map[int32]string {
	m := make(map[int32]string, len(e.Items))
	for _, item := range e.Items {
		m[item.Idx] = item.Val
	}
	return m
}

func (e *Enum) FromMap(m map[int32]string) {
	e.Items = make([]EnumItem, 0, len(m))
	for index, item := range m {
		e.Items = append(e.Items, EnumItem{
			Idx: index,
			Val: item,
		})
	}
	sort.Slice(e.Items, func(i, j int) bool {
		return e.Items[i].Idx < e.Items[j].Idx
	})
}

func (e *Enum) ItemByValue(v string) (*EnumItem, api_error.ApiErr) {
	for i := range e.Items {
		if strings.EqualFold(v, e.Items[i].Val) {
			return &e.Items[i], nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("No item with value %v found", v))
}

func (e *Enum) ItemByIndex(i int32) (*EnumItem, api_error.ApiErr) {
	for iItem := range e.Items {
		if e.Items[iItem].Idx == i {
			return &e.Items[iItem], nil
		}
	}
	return nil, api_error.NewNotFoundError(fmt.Sprintf("No item with index %v found", i))
}
