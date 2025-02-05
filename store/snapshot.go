package store

import (
	"encoding/gob"
	"fmt"
	"github.com/nccapo/kv-store/internal/gen"
	"os"
)

func (kv *Store) saveSnapshot() error {
	s := gen.RandStringRunes(10)

	temp, err := os.CreateTemp("", fmt.Sprintf("sanpshot-%s.tmp", s))
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())

	encoder := gob.NewEncoder(temp)
	if err = encoder.Encode(kv.data); err != nil {
		return err
	}

	if err = os.Rename(temp.Name(), "snapshot.gob"); err != nil {
		return err
	}

	return nil
}
