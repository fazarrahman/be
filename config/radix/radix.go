package radix

import (
	"log"

	"github.com/mediocregopher/radix/v3"
)

func New() (*radix.Pool, error) {
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return pool, nil
}
