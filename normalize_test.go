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
		"HtTp://spHela.com:80/foo?baz=moo",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:7334:80/path/tostuff",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:80/path/tostuff",
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
		"http://sphela.com/foo?baz=moo",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:7334/path/tostuff",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:80/path/tostuff",
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
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			Normalize(URL)
			receivedURL := URL.String()
			if receivedURL != normalizedURLs[i] {
				t.Error("Received URL not normalized", receivedURL,
					normalizedURLs[i])
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
	if URL, err := url.ParseWithReference(checkURL); err == nil {
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

func TestNormalizeDomain(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?foo=bar",
		"https://gooooogle.com/search/",
		"http://gogl.net/",
	}
	formatted := [...]string{
		"http://www.google.com/path/tostuff/?foo=bar",
		"https://www.google.com/search/",
		"http://www.google.com/",
	}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			NormalizeDomain(URL, "www.google.com")
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("NormalizeDomain failed", checkURL, receivedURL,
					formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestNormalizeQuery(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?foo=bar&delete=this",
		"https://gooooogle.com/search/?nothing&wtf=this&fuzz=bar",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/?foo=bar",
		"https://gooooogle.com/search/?fuzz=bar",
		"http://gogl.net/?fuzz=baz&foo=bar",
	}
	params := [...]string{"fuzz", "foo"}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			NormalizeQuery(URL, params[:])
			formattedURL, _ := url.ParseWithReference(formatted[i])
			NormalizeQueryVariableOrder(URL)
			NormalizeQueryVariableOrder(formattedURL)
			receivedURL := URL.String()
			expected := formattedURL.String()
			if receivedURL != expected {
				t.Error("NormalizeQuery failed", checkURL, receivedURL,
					expected)
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestNormalizeQueryVariableOrder(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"https://gooooogle.com/search/?nothing&flow=this&car=bar",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/?atari=this&zar=bar",
		"https://gooooogle.com/search/?car=bar&flow=this&nothing",
		"http://gogl.net/?foo=bar&fuzz=baz&snow=cold",
	}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			NormalizeQueryVariableOrder(URL)
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("NormalizeQueryVariableOrder failed", checkURL,
					receivedURL, formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestNormalizeScheme(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"https://gooooogle.com/search/?nothing&flow=this&car=bar",
		"file://gogl.net/?fuzz=baz&snow=cold&foo=bar",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"http://gooooogle.com/search/?nothing&flow=this&car=bar",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
	}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			NormalizeScheme(URL, "http")
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("NormalizeScheme failed", checkURL,
					receivedURL, formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestNormalizeWWWShow(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:7334/path/tostuff",
		"https://www.gooooogle.com/search/?nothing&flow=this&car=bar",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
		"http://gogl.net:8080/?fuzz=baz&snow=cold&foo=bar",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:7334/path/tostuff",
		"https://www.gooooogle.com/search/?nothing&flow=this&car=bar",
		"http://www.gogl.net/?fuzz=baz&snow=cold&foo=bar",
		"http://www.gogl.net:8080/?fuzz=baz&snow=cold&foo=bar",
	}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			NormalizeWWW(URL, true)
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("NormalizeWWW show failed", checkURL,
					receivedURL, formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestNormalizeWWWHide(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:7334/path/tostuff",
		"https://www.gooooogle.com/search/?nothing&flow=this&car=bar",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
		"http://www.gogl.net:8080/?fuzz=baz&snow=cold&foo=bar",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/?zar=bar&atari=this",
		"http://2001:0db8:85a3:0000:0000:8a2e:0370:7334/path/tostuff",
		"https://gooooogle.com/search/?nothing&flow=this&car=bar",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
		"http://gogl.net:8080/?fuzz=baz&snow=cold&foo=bar",
	}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			NormalizeWWW(URL, false)
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("NormalizeWWW hide failed", checkURL,
					receivedURL, formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestRemoveDefaultQueryValues(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/?foo=bar&fuzz=234&atari=this",
		"https://www.gooooogle.com/search/?nothing&flow=this&car=bar&foo=1",
		"http://gogl.net/?fuzz=baz&snow=cold&foo=bar",
	}
	expectedValues := [...]string{
		"http://74.125.224.49/path/tostuff/?fuzz=234&atari=this",
		"https://www.gooooogle.com/search/?nothing&flow=this&car=bar&foo=1",
		"http://gogl.net/?snow=cold",
	}
	defaults := map[string]string{
		"foo":  "bar",
		"fuzz": "baz",
	}
	for i, checkURL := range urls {
		URL, _ := url.ParseWithReference(checkURL)
		formattedURL, _ := url.ParseWithReference(expectedValues[i])
		RemoveDefaultQueryValues(URL, defaults)
		NormalizeQueryVariableOrder(URL)
		NormalizeQueryVariableOrder(formattedURL)
		receivedURL := URL.String()
		formatted := formattedURL.String()
		if receivedURL != formatted {
			t.Error("RemoveDefaultQueryValues failed", checkURL,
				receivedURL, formatted)
		}
	}
}

func TestRemoveDirectoryIndex(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/index.html/?foo=bar&fuzz=234",
		"https://www.gooooogle.com/search/index.html?nothing&flow=index.html",
		"http://gogl.net/index.html#index.html",
		"http://gogl.net/index.html",
		"http://google.com/",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/index.html/?foo=bar&fuzz=234",
		"https://www.gooooogle.com/search/?nothing&flow=index.html",
		"http://gogl.net/#index.html",
		"http://gogl.net/",
		"http://google.com/",
	}
	index := "index.html"
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			RemoveDirectoryIndex(URL, index)
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("RemoveDirectoryIndex failed", checkURL,
					receivedURL, formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}

func TestRemoveFragment(t *testing.T) {
	urls := [...]string{
		"http://74.125.224.49/path/tostuff/index.html/?foo=bar&fuzz=234#moo",
		"https://www.google.com/search/index.html#?nothing&flow=index.html",
		"http://gogl.net/index.html#test#more#tests",
		"http://gogl.net/index.html#index.html",
	}
	formatted := [...]string{
		"http://74.125.224.49/path/tostuff/index.html/?foo=bar&fuzz=234",
		"https://www.google.com/search/index.html",
		"http://gogl.net/index.html",
		"http://gogl.net/index.html",
	}
	for i, checkURL := range urls {
		if URL, err := url.ParseWithReference(checkURL); err == nil {
			RemoveFragment(URL)
			receivedURL := URL.String()
			if receivedURL != formatted[i] {
				t.Error("RemoveFragment failed", checkURL,
					receivedURL, formatted[i])
			}
		} else {
			t.Error("Error parsing URL")
		}
	}
}
