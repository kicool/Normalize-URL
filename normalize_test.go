package normalize

import (
  "testing"
  "strconv"
  "strings"
  "url"
)

func TestNormalize(t *testing.T) {
  rawURLs := [...]string{
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
  normalizedURLs := [...]string{
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
  for i, checkURL := range rawURLs {
    if URL, err := url.Parse(checkURL); err == nil {
      Normalize(URL)
      receivedURL := URL.String()
      if receivedURL != normalizedURLs[i] {
        t.Error("Received URL not normalized", receivedURL, normalizedURLs[i])
      }
    } else {
      t.Error("Error while parsing ", err)
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
  if URL, err := url.Parse(checkURL); err == nil {
    Normalize(URL)
    receivedURL := URL.String()
    if receivedURL != normalizedURL {
      t.Error("Character not escaped right.", checkURL,
        normalizedURL, receivedURL)
    }
  } else {
    t.Error("Error while parsing ", err)
  }
}

func TestControlChars(t *testing.T) {
  for i := 0; i < controlCharEnd; i++ {
    testChar(t, i)
  }
}

func TestReservedChars(t *testing.T) {
  for val, _ := range reservedChars {
    testChar(t, val)
  }
}

func TestSomeUnsafeChars(t *testing.T) {
  for val, _ := range unsafeChars {
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
  for i := 0; i < 256; i++ {
    _, reserved := reservedChars[i]
    _, unsafe := unsafeChars[i]
    switch {
    default:
      //Test char to make sure it's not escaped after normalized.
      t.Log("Searching and testing against", i)
    case i <= controlCharEnd:
      t.Log("Less than controlCharEnd", i, controlCharEnd)
      continue
    case i >= nonASCIImin && i <= nonASCIImax:
      t.Log("In non-ASCII range", i, nonASCIImin, nonASCIImax)
      continue
    case reserved:
      t.Log("In reservedChars", i, reservedChars[i])
      continue
    case unsafe:
      t.Log("In unsafeChars", i, unsafeChars[i])
      continue
    }
  }
}
