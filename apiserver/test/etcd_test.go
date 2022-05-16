package test

import (
	"github.com/stretchr/testify/assert"
	"minik8s/apiserver/src/etcd"
	"testing"
)

func TestEtcd(t *testing.T) {
	err := etcd.Put("22", "33")
	assert.Nil(t, err)
	err = etcd.Put("22", "1234")
	assert.Nil(t, err)
	var value string
	value, err = etcd.Get("22")
	assert.NotEqual(t, "22", value)
	value, err = etcd.Get("22")
	assert.NotEqual(t, "33", value)
	value, err = etcd.Get("22")
	assert.Equal(t, "1234", value)

	err = etcd.Delete("22")
	assert.Nil(t, err)
	value, err = etcd.Get("22")
	assert.Equal(t, "", value)
}
