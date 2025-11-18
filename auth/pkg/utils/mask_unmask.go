package utils

import (
	"errors"
	"fmt"
	"log"

	"github.com/sqids/sqids-go"
)

type MaskedId string

var (
	s *sqids.Sqids

	ErrInvalidMaskedID = errors.New("invalid masked ID")
)

func init() {
	var err error
	s, err = sqids.New(sqids.Options{
		MinLength: 8,
	})
	if err != nil {
		log.Fatalf("Failed to initialize sqids generator: %v", err)
	}
}

func Mask(id uint) MaskedId {
	if id == 0 {
		return ""
	}

	masked, _ := s.Encode([]uint64{uint64(id)})
	return MaskedId(masked)
}

func Unmask(sid MaskedId) uint {
	id, err := UnmaskWithError(sid)
	if err != nil {
		return 0
	}
	return id
}

func UnmaskWithError(sid MaskedId) (uint, error) {
	str := string(sid)
	if str == "" {
		return 0, ErrInvalidMaskedID
	}

	ids := s.Decode(str)
	if len(ids) == 0 {
		return 0, ErrInvalidMaskedID
	}

	return uint(ids[0]), nil
}

func (m MaskedId) String() string {
	return string(m)
}

func (m MaskedId) Unmask() uint {
	return Unmask(m)
}

func (m MaskedId) Valid() bool {
	_, err := UnmaskWithError(m)
	return err == nil
}

func GetMaskedId(s string) MaskedId {
	return MaskedId(s)
}

func (m MaskedId) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, m.String())), nil
}

func (m *MaskedId) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid JSON for MaskedId")
	}
	*m = MaskedId(data[1 : len(data)-1])
	return nil
}
