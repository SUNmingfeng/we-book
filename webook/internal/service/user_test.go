package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("12345678#")
	encrypt, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	fmt.Printf("密码加密：%v\n", string(encrypt))
	err = bcrypt.CompareHashAndPassword(encrypt, []byte("123456789"))
	assert.NotNil(t, err)
}
