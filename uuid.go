package uuid

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/gofrs/uuid"
)

// swagger:strfmt uuid
type UUID [16]byte

var Zero = UUID{}

func (u UUID) IsZero() bool {
	return u == Zero
}

func (u UUID) String() string {
	buf := make([]byte, 32)
	hex.Encode(buf[:], u[:])
	return string(buf)
}

func (u UUID) Bytes() []byte {
	return u[:]
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (u UUID) MarshalBinary() (data []byte, err error) {
	data = u.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return error if the slice isn't 16 bytes long.
func (u *UUID) UnmarshalBinary(data []byte) (err error) {
	if len(data) != 16 {
		err = fmt.Errorf("uuid: UUID must be exactly 16 bytes long, got %d bytes", len(data))
		return
	}
	copy(u[:], data)
	return
}

func (u *UUID) UnmarshalText(text []byte) (err error) {
	if len(text) < 32 {
		err = fmt.Errorf("uuid: invalid UUID string: %s", text)
		return
	}
	if len(text) == 36 {
		text = bytes.Replace(text, []byte("-"), []byte(""), -1)
	}
	_, err = hex.Decode(u[:], text)
	if err != nil {
		return
	}
	return
}

func (u UUID) MarshalText() (data []byte, err error) {
	data = []byte(u.String())
	return
}

// Value implements the driver.Valuer interface.
func (u UUID) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return u.String(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice is handled by UnmarshalBinary, while
// a longer byte slice or a string is handled by UnmarshalText.
func (u *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		if len(src) == 16 {
			return u.UnmarshalBinary(src)
		} else {
			//return fmt.Errorf("uuid: it is not a length of 16, its length is %v and its value is %s", len(src), src)
		}
		return u.UnmarshalText(src)
	case string:
		return u.UnmarshalText([]byte(src))
	case nil:
		u = &UUID{}
		return nil
	}
	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// FromBytes returns UUID converted from raw byte slice input.
// It will return error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (u UUID, err error) {
	err = u.UnmarshalBinary(input)
	return
}

func FromString(input string) (u UUID, err error) {
	err = u.UnmarshalText([]byte(input))
	return
}

func NewV1() (UUID, error) {
	var u UUID
	n, err := uuid.NewV1()
	if err != nil {
		return u, err
	}
	u = UUID(n)
	return u, nil
}

func NewV7() (UUID, error) {
	var u UUID
	n, err := uuid.NewV7()
	if err != nil {
		return u, err
	}
	u = UUID(n)
	return u, nil
}

func NewV4() (UUID, error) {
	var u UUID
	n, err := uuid.NewV4()
	if err != nil {
		return u, err
	}
	u = UUID(n)
	return u, nil
}
