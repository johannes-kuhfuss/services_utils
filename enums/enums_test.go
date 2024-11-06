package enums

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	EmptyEnum Enum
	TestEnum1 Enum
	TestEnum2 Enum
	eItemOne  EnumItem = EnumItem{
		Idx: 0,
		Val: "test",
	}
	eItemTwo EnumItem = EnumItem{
		Idx: 1,
		Val: "pinguin",
	}
)

func setup() func() {
	TestEnum1.Items = append(TestEnum1.Items, eItemOne)
	TestEnum1.Items = append(TestEnum1.Items, eItemTwo)
	return func() {
		TestEnum1.Items = nil
	}
}

func TestItemByValueNoItemReturnsNotFoundError(t *testing.T) {
	item, err := EmptyEnum.ItemByValue("test")

	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with value test found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func TestItemByValueTestItemReturnsTestItemValue(t *testing.T) {
	teardown := setup()
	defer teardown()

	item, err := TestEnum1.ItemByValue("test")

	assert.NotNil(t, item)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, item.Idx)
	assert.EqualValues(t, "test", item.Val)
}

func TestItemByIndexNoItemReturnsNotFoundError(t *testing.T) {
	item, err := EmptyEnum.ItemByIndex(1)

	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with index 1 found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func TestItemByIndexTestItemReturnsTestItemIndex(t *testing.T) {
	teardown := setup()
	defer teardown()

	item, err := TestEnum1.ItemByIndex(1)

	assert.NotNil(t, item)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, item.Idx)
	assert.EqualValues(t, "pinguin", item.Val)
}

func TestValueNoItemReturnsNotFoundError(t *testing.T) {
	val, err := EmptyEnum.AsValue(0)

	assert.Empty(t, val)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with index 0 found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func TestValueTestItemReturnsValue(t *testing.T) {
	teardown := setup()
	defer teardown()

	val, err := TestEnum1.AsValue(0)

	assert.NotNil(t, val)
	assert.Nil(t, err)
	assert.EqualValues(t, "test", val)
}

func TestIndexNoItemReturnsNotFoundError(t *testing.T) {
	val, err := EmptyEnum.AsIndex("test")

	assert.Empty(t, val)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with value test found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func TestIndexTestItemReturnsValue(t *testing.T) {
	teardown := setup()
	defer teardown()

	val, err := TestEnum1.AsIndex("pinguin")

	assert.NotNil(t, val)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, val)
}

func TestValuesNoItemsReturnsEmptyStringSlice(t *testing.T) {
	vals := EmptyEnum.Values()

	assert.IsType(t, []string{}, vals)
	assert.EqualValues(t, 0, len(vals))
}

func TestValuesTwoItemsReturnsStringSlice(t *testing.T) {
	teardown := setup()
	defer teardown()

	vals := TestEnum1.Values()

	assert.IsType(t, []string{}, vals)
	assert.EqualValues(t, 2, len(vals))
	assert.EqualValues(t, "test", vals[0])
	assert.EqualValues(t, "pinguin", vals[1])
}

func TestAsMapNoItemsReturnsEmptyMap(t *testing.T) {
	m := EmptyEnum.AsMap()

	assert.IsType(t, map[int32]string{}, m)
	assert.EqualValues(t, 0, len(m))
}

func TestAsMapTwoItemsReturnsMap(t *testing.T) {
	teardown := setup()
	defer teardown()

	m := TestEnum1.AsMap()

	assert.NotNil(t, m)
	assert.IsType(t, map[int32]string{}, m)
	assert.EqualValues(t, 2, len(m))
	assert.EqualValues(t, "test", m[0])
	assert.EqualValues(t, "pinguin", m[1])
}

func TestFromMapEmptyMapReturnsEmptyEnum(t *testing.T) {
	m := make(map[int32]string)

	TestEnum2.FromMap(m)

	assert.NotNil(t, TestEnum2)
	assert.EqualValues(t, 0, len(TestEnum2.Items))
}

func TestFromMapMapReturnsEnum(t *testing.T) {
	m := make(map[int32]string)
	m[3] = "extra"
	m[4] = "super"

	TestEnum2.FromMap(m)

	val3, _ := TestEnum2.AsValue(3)
	val4, _ := TestEnum2.AsValue(4)
	assert.NotNil(t, TestEnum2)
	assert.EqualValues(t, 2, len(TestEnum2.Items))
	assert.EqualValues(t, val3, m[3])
	assert.EqualValues(t, val4, m[4])
}
