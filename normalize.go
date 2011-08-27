package normalize

import (
  "os"
)

func Normalize(rawurl string) (url string, err os.Error) {
  url = rawurl
  return url, nil
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

//Character values 0-31 need to be escaped in query strings:
var controlCharEnd int = 31
var reservedChars = [...]int{
  36,  //$
  38,  //&
  43,  //+
  44,  //,
  47,  ///
  58,  //:
  59,  //;
  61,  //=
  63,  //?
  64,  //@
  127, //.
}
var unsafeChars = [...]int{
  32,  //space
  34,  //"
  35,  //#
  37,  //%
  60,  //<
  62,  //
  91,  //[
  92,  //\
  93,  //]
  94,  //^
  96,  //`
  123, //{
  124, //|
  125, //}
  126, //~
}
//Character values 128-255 need to be escaped.
var nonASCIImin int = 128
var nonASCIImax int = 255
