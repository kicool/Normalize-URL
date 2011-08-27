package normalize

import (
  "os"
)

func Normalize(rawurl string) (url string, err os.Error) {
  return "Not yet implemented.", *(NewNormalizeError("Not yet implmented"))
}

func NewNormalizeError(description string) *NormalizeError {
  err := new(NormalizeError)
  err.err = description
  return err
}

type NormalizeError struct {
  err string
}

func (err NormalizeError) String() string {
  return err.err
}
