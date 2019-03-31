package cache

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type metricConfSuit struct {
	suite.Suite
}

func TestCacheSuit(t *testing.T) {
	suite.Run(t, new(metricConfSuit))
}

func (suite *metricConfSuit) TestCanSetValue() {
	// NOTE(denisacostaq@gmail.com): Giving
	mc := NewCache()
	key := "dfdhj&**"
	val := []byte("dfdfdfdfdfd4545")

	// NOTE(denisacostaq@gmail.com): When
	mc.Set(key, val)
	rVal, err := mc.Get(key)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(val, rVal)
}

func (suite *metricConfSuit) TestDoNotChangeInternalStateOnCopyByValue() {
	// NOTE(denisacostaq@gmail.com): Giving
	mc := NewCache()
	key := "dfdhj&**"
	val := []byte("dfdfdfdfdfd4545")

	// NOTE(denisacostaq@gmail.com): When
	mc.Set(key, val)
	var cache1, cache2 Cache
	cache1 = mc
	cache2 = cache1
	rVal, err := cache2.Get(key)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(val, rVal)
}

func (suite *metricConfSuit) TestCanResetCache() {
	// NOTE(denisacostaq@gmail.com): Giving
	mc := NewCache()
	key := "dfdhj&**"
	val := []byte("dfdfdfdfdfd4545")

	// NOTE(denisacostaq@gmail.com): When
	mc.Set(key, val)
	rVal1, err1 := mc.Get(key)
	mc.Reset()
	_, err2 := mc.Get(key)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err1)
	suite.Equal(val, rVal1)
	suite.NotNil(err2)
}
