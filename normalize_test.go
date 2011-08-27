package normalize

import (
  "testing"
  "strconv"
  "strings"
  "sort"
  "fmt"
)

var rawURLs = [...]string{
  "HtTp://spHela.com",
  "HTTps://www.EXAMPLE.COM/%2d%aD/MOO#smoo",
  "HTTps://www.EXAMPLE.COM/%2d%aD/?MO=O smoo",
  "HTTps://www.EXAMPLE.COM/%2d%aD/MOO ",
  "http://apphacker.com/moo/../doo/./baz/",
  "http://apphacker.com/moo//doo//baz/",
  "http://www.apphacker.com?",
  "http://www.apphacker.com?boo=fuzz",
  "http://apphacker.com/?",
  "http://apphacker.com/?%foo=bar",
}

var normalizedURLs = [...]string{
  "http://sphela.com/",
  "https://www.example.com/~%AD/MOO#smoo",
  "https://www.example.com/~%AD/?MO=O%20smoo",
  "https://www.example.com/~%AD/MOO",
  "http://apphacker.com/moo/doo/baz/",
  "http://apphacker.com/moo/doo/baz/",
  "http://www.apphacker.com/",
  "http://www.apphacker.com/?boo=fuzz",
  "http://apphacker.com/",
  "http://apphacker.com/?%25foo=bar",
}

func TestNormalize(t *testing.T) {
  for i, checkURL := range rawURLs {
    if receivedURL, err := Normalize(checkURL); err == nil {
      if receivedURL != normalizedURLs[i] {
        t.Error("Received URL not normalized", receivedURL, normalizedURLs[i])
      }
    } else {
      t.Error("Error while normalizing ", err)
    }
  }
}

func testChar(t *testing.T, val int) {
  var (
    char byte
    hex  string
  )
  testURL := "http://google.com?moo=doo"
  char = byte(val)
  hex = strings.ToUpper(strconv.Itob(val, 16))
  if len(hex) == 1 {
    hex = "0" + hex
  }
  checkURL := testURL + string(char)
  normalizedURL := testURL
  if char != ' ' {
    //Trailing whitespaces are removed.
    normalizedURL = normalizedURL + "%" + hex
  }
  if receivedURL, err := Normalize(checkURL); err == nil {
    if receivedURL != normalizedURL {
      t.Error("Character not escaped right.", checkURL,
        normalizedURL, receivedURL)
    }
  } else {
    t.Error("Error while normalizing ", err)
  }
}

func TestControlChars(t *testing.T) {
  for i := 0; i < controlCharEnd; i++ {
    testChar(t, i)
  }
}

func TestReservedChars(t *testing.T) {
  for _, val := range reservedChars {
    testChar(t, val)
  }
}

func TestSomeUnsafeChars(t *testing.T) {
  for _, val := range reservedChars {
    if val != 35 || val != 37 {
      testChar(t, val)
    }
  }
}

func TestNonASCIIChars(t *testing.T) {
  for i := nonASCIImin; i <= nonASCIImax; i++ {
    testChar(t, i)
  }
}

func TestUnescapeChars(t *testing.T) {
  searchSlice := make(sort.IntSlice, len(reservedChars)+len(unsafeChars))
  copy(searchSlice, reservedChars[:])
  copy(searchSlice, unsafeChars[:])
  searchSlice.Sort()
  for i := 0; i < 256; i++ {
    index := searchSlice.Search(i)
    fmt.Println("Index:", index, "searchSlice[index]:", searchSlice[index],
      "len(searchSlice):", len(searchSlice), "i:", i)
    switch {
    default:
      //Do something
      fmt.Println("Searching and testing against", i)
    case i <= controlCharEnd:
      fmt.Println("Less than controlCharEnd", i, controlCharEnd)
      continue
    case i >= nonASCIImin && i <= nonASCIImax:
      fmt.Println("In non-ASCII range", i, nonASCIImin, nonASCIImax)
      continue
    case index < len(searchSlice):
      fmt.Println("Found in searchSlice i:", i, "index:", index,
        "searchSlice[index]:", searchSlice[index])
      continue
    }
  }
}
