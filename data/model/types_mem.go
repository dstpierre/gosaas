// +build mem

package model

type Connection = bool
type Key = int

func Open(options ...string) (bool, error) {
	return true, nil
}
