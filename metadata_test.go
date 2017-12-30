package exiftool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var metaTest = &Metadata{raw: map[string]interface{}{
	"MIMEType": "test/data",
	"string":   "Hello",
	"float64":  float64(123.45),
	"int":      float64(10),
	"date":     "2010:03:04 01:02:03",
	"dateMS":   "2010:03:04 01:02:03.456",
}}

func TestMIMEType(t *testing.T) {
	assert.Equal(t, "test/data", metaTest.MIMEType())
}

func TestCreateDate(t *testing.T) {
	{ // test with no CreateDate key
		_, ok := metaTest.CreateDate()
		assert.False(t, ok)
	}

	{ // with with a CreateDate key
		tdata := &Metadata{raw: map[string]interface{}{
			"CreateDate": "2010:03:04 01:02:03",
		}}

		create, ok := tdata.CreateDate()
		assert.True(t, ok)
		assert.Equal(t, "2010:03:04 01:02:03", create.Format(TimeFormat))
	}

	{ // a "SubSecCreateDate" has precedence over "CreateDate"
		tdata := &Metadata{raw: map[string]interface{}{
			"CreateDate":       "2009:03:04 01:02:03",
			"SubSecCreateDate": "2010:03:04 01:02:03.456",
		}}

		create, ok := tdata.CreateDate()
		assert.True(t, ok)
		assert.Equal(t, "2010:03:04 01:02:03.456", create.Format(TimeFormatMS))
	}
}

func TestTypeConversion(t *testing.T) {
	assert := assert.New(t)

	var (
		tString  string
		tFloat64 float64
		tInt     int
		tTime    time.Time
		err      error
	)

	tString, err = metaTest.GetString("string")
	assert.Equal("Hello", tString)
	assert.NoError(err)

	tFloat64, err = metaTest.GetFloat64("float64")
	assert.Equal(float64(123.45), tFloat64)
	assert.NoError(err)

	tInt, err = metaTest.GetInt("int")
	assert.Equal(10, tInt)
	assert.NoError(err)

	tTime, err = metaTest.GetDate("date")
	assert.NoError(err)
	assert.Equal("2010:03:04 01:02:03.000", tTime.Format(TimeFormatMS))

	// parse with millisecond precision
	tTime, err = metaTest.GetDate("dateMS")
	assert.NoError(err)
	assert.Equal("2010:03:04 01:02:03.456", tTime.Format(TimeFormatMS))
}

func TestKeyHelpers(t *testing.T) {
	assert := assert.New(t)
	assert.False(metaTest.KeyExists("nope"))

	for _, key := range metaTest.Keys() {
		assert.True(metaTest.KeyExists(key))
	}
}
