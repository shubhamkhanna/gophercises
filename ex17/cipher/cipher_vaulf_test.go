package cipher

import (
	"crypto/cipher"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

var home, _ = homedir.Dir()

type tempStruct struct {
	err   error
	block cipher.Block
}

func (m *tempStruct) dummyNewCipherBlock(key string) (cipher.Block, error) {
	return nil, m.err
}

func TestSet(t *testing.T) {

	t.Run("it set the keys to voult", func(t *testing.T) {
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		v.Set("fb_api_key", "itssecret")
		assert.Equal(t, v.keyValues["fb_api_key"], "itssecret")
	})

	t.Run("it returns an error", func(t *testing.T) {
		v := Vault{}
		err := v.Set("fb_api_key", "itssecret")
		assert.NotEqual(t, err, nil)
	})

	t.Run("it returns an error while cipher process", func(t *testing.T) {
		m := &tempStruct{err: errors.New("failed")}
		mocknewCipherBlock = m.dummyNewCipherBlock
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		err := v.Set("fb_api_key", "itssecret")
		assert.NotEqual(t, err, nil)
	})
	defer func() {
		mocknewCipherBlock = newCipherBlock // set back original func at end of test
	}()
}

func TestGet(t *testing.T) {

	t.Run("it get the key values from voult", func(t *testing.T) {
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		value, _ := v.Get("fb_api_key")
		assert.Equal(t, value, "itssecret")
	})

	t.Run("it returns an error", func(t *testing.T) {
		v := Vault{}
		_, err := v.Get("fb_api_key")
		assert.NotEqual(t, err, nil)
	})

	t.Run("it returns an error while cipher process", func(t *testing.T) {
		m := &tempStruct{err: errors.New("failed")}
		mocknewCipherBlock = m.dummyNewCipherBlock
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		_, err := v.Get("fb_api_key")
		assert.NotEqual(t, err, nil)
	})
	defer func() {
		mocknewCipherBlock = newCipherBlock // set back original func at end of test
	}()
}

func TestRemove(t *testing.T) {

	t.Run("it remove the key values from voult", func(t *testing.T) {
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		v.Set("remove_key", "remove_me")
		v.Remove("remove_key")
		assert.Equal(t, v.keyValues["remove_key"], "")
	})

	t.Run("it returns an error", func(t *testing.T) {
		v := Vault{}
		v.Set("remove_key", "remove_me")
		err := v.Remove("remove_key")
		assert.NotEqual(t, err, nil)
	})

	t.Run("it returns an error while cipher process", func(t *testing.T) {
		m := &tempStruct{err: errors.New("failed")}
		mocknewCipherBlock = m.dummyNewCipherBlock
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		err := v.Remove("remove_key")
		assert.NotEqual(t, err, nil)
	})
	defer func() {
		mocknewCipherBlock = newCipherBlock // set back original func at end of test
	}()
}

func TestFile(t *testing.T) {
	path := filepath.Join(home, ".secrets")
	v := File("encoding_key", path)
	assert.Equal(t, v.encodingKey, "encoding_key")
	assert.Equal(t, v.filepath, path)
}

func TestEncryptWriter(t *testing.T) {
	t.Run("it returns stream writter", func(t *testing.T) {
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		f, _ := os.Open(v.filepath)
		w, _ := EncryptWriter("encyp_writer_key", f)
		assert.Equal(t, reflect.TypeOf(w).String(), "*cipher.StreamWriter")
	})

	t.Run("it returns an error", func(t *testing.T) {
		m := &tempStruct{err: errors.New("failed")}
		mocknewCipherBlock = m.dummyNewCipherBlock
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		f, _ := os.Open(v.filepath)
		_, err := EncryptWriter("", f)
		assert.NotEqual(t, err, nil)

		defer func() {
			mocknewCipherBlock = newCipherBlock // set back original func at end of test
		}()
	})
}

func TestDecryptReader(t *testing.T) {
	t.Run("it returns stream reader", func(t *testing.T) {
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		f, _ := os.Open(v.filepath)
		r, _ := DecryptReader("encyp_writer_key", f)
		assert.Equal(t, reflect.TypeOf(r).String(), "*cipher.StreamReader")
	})

	t.Run("it returns an error", func(t *testing.T) {
		m := &tempStruct{err: errors.New("failed")}
		mocknewCipherBlock = m.dummyNewCipherBlock
		v := Vault{filepath: filepath.Join(home, ".secrets")}
		f, _ := os.Open(v.filepath)
		_, err := DecryptReader("", f)
		assert.NotEqual(t, err, nil)

		defer func() {
			mocknewCipherBlock = newCipherBlock // set back original func at end of test
		}()
	})
}
