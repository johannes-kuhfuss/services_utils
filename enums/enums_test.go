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
		idx: 0,
		val: "test",
	}
	eItemTwo EnumItem = EnumItem{
		idx: 1,
		val: "pinguin",
	}
)

func setup() func() {
	TestEnum1.items = append(TestEnum1.items, eItemOne)
	TestEnum1.items = append(TestEnum1.items, eItemTwo)
	return func() {
		TestEnum1.items = nil
	}
}

func Test_ItemByValue_NoItem_Returns_NotFoundError(t *testing.T) {
	item, err := EmptyEnum.ItemByValue("test")

	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with value test found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_ItemByValue_TestItem_Returns_TestItemValue(t *testing.T) {
	teardown := setup()
	defer teardown()

	item, err := TestEnum1.ItemByValue("test")

	assert.NotNil(t, item)
	assert.Nil(t, err)
	assert.EqualValues(t, 0, item.idx)
	assert.EqualValues(t, "test", item.val)
}

func TestItemByIndex_NoItem_Returns_NotFoundError(t *testing.T) {
	item, err := EmptyEnum.ItemByIndex(1)

	assert.Nil(t, item)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with index 1 found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func TestItemByIndex_TestItem_Returns_TestItemIndex(t *testing.T) {
	teardown := setup()
	defer teardown()

	item, err := TestEnum1.ItemByIndex(1)

	assert.NotNil(t, item)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, item.idx)
	assert.EqualValues(t, "pinguin", item.val)
}

func Test_Value_NoItem_Returns_NotFoundError(t *testing.T) {
	val, err := EmptyEnum.Value(0)

	assert.Empty(t, val)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with index 0 found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_Value_TestItem_Returns_Value(t *testing.T) {
	teardown := setup()
	defer teardown()

	val, err := TestEnum1.Value(0)

	assert.NotNil(t, val)
	assert.Nil(t, err)
	assert.EqualValues(t, "test", val)
}

func Test_Index_NoItem_Returns_NotFoundError(t *testing.T) {
	val, err := EmptyEnum.Index("test")

	assert.Empty(t, val)
	assert.NotNil(t, err)
	assert.EqualValues(t, "No item with value test found", err.Message())
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode())
}

func Test_Index_TestItem_Returns_Value(t *testing.T) {
	teardown := setup()
	defer teardown()

	val, err := TestEnum1.Index("pinguin")

	assert.NotNil(t, val)
	assert.Nil(t, err)
	assert.EqualValues(t, 1, val)
}

func Test_Values_NoItems_Returns_EmptyStringSlice(t *testing.T) {
	vals := EmptyEnum.Values()

	assert.IsType(t, []string{}, vals)
	assert.EqualValues(t, 0, len(vals))
}

func Test_Values_TwoItems_Returns_StringSlice(t *testing.T) {
	teardown := setup()
	defer teardown()

	vals := TestEnum1.Values()

	assert.IsType(t, []string{}, vals)
	assert.EqualValues(t, 2, len(vals))
	assert.EqualValues(t, "test", vals[0])
	assert.EqualValues(t, "pinguin", vals[1])
}

func Test_AsMap_NoItems_Returns_EmptyMap(t *testing.T) {
	m := EmptyEnum.AsMap()

	assert.IsType(t, map[int]string{}, m)
	assert.EqualValues(t, 0, len(m))
}

func Test_AsMap_TwoItems_Returns_Map(t *testing.T) {
	teardown := setup()
	defer teardown()

	m := TestEnum1.AsMap()

	assert.NotNil(t, m)
	assert.IsType(t, map[int]string{}, m)
	assert.EqualValues(t, 2, len(m))
	assert.EqualValues(t, "test", m[0])
	assert.EqualValues(t, "pinguin", m[1])
}

func Test_FromMap_EmptyMap_Returns_EmptyEnum(t *testing.T) {
	m := make(map[int]string)

	TestEnum2.FromMap(m)

	assert.NotNil(t, TestEnum2)
	assert.EqualValues(t, 0, len(TestEnum2.items))
}

func Test_FromMap_Map_Returns_Enum(t *testing.T) {
	m := make(map[int]string)
	m[3] = "extra"
	m[4] = "super"

	TestEnum2.FromMap(m)

	val3, _ := TestEnum2.Value(3)
	val4, _ := TestEnum2.Value(4)
	assert.NotNil(t, TestEnum2)
	assert.EqualValues(t, 2, len(TestEnum2.items))
	assert.EqualValues(t, val3, m[3])
	assert.EqualValues(t, val4, m[4])
}
