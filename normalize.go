package normalize

import (
	"os"
	"url"
)

//Naive normalization, normalizes those aspects of a URL it can
//without knowing much about it. Does not make changes that might
//change the location which the URL points to
func Normalize(url *url.URL) (err os.Error) {
	return nil
}

//Removes directory indexes when they point to the same place as
//the directory. For example if index.html points to / and
//index.html is given for the index parameter it will be removed
//from the URL
func RemoveDirectoryIndex(url *url.URL, index string) {
}

//Ordes query variables in alphabetic order. Order of variables
//in a query string should not matter, but some implementations
//may require an order, so this is in a separate emthod.
func NormalizeQueryVariableOrder(url *url.URL) {
}

//Remove query variables that have default values.  Provide a set of defaults
//(defaults[key] = value) wher key is the variable name and value is the string
//represenation of the default value.
func RemoveDefaultQueryValues(url *url.URL, defaults map[string] string) {
}

//Removes www. from a URL. Use if www. points to same resource as
//non-www address.
func NormalizeWWW(url *url.URL, showWWW bool) {
}

//Remove arbitary query variables. Include a slice of array variables
//to check against. If query variables are found not in the given slice,
//they are removed.
func NormalizeQuery(url *url.URL, params []string) {
}

//Normalize scheme or protocol. For example if valid scheme is url
//and not urls, url is changed to url if 'url' is given as scheme.
func NormalizeScheme(url *url.URL, scheme string) {
}

//Remove #fragment from a URL.
func RemoveFragment(url *url.URL) {
}

//Replaces domain or IP with given domain. Use to replace IP addresses with
//domain or domains that point to the same resource as prime domain.
func NormalizeDomain(url *url.URL, domain string) {
}

func NewNormalizeError(description string) *NormalizeError {
	err := new(NormalizeError)
	err.err = description
	return err
}

type NormalizeError struct {
	os.Error
	err string
}

func (err NormalizeError) String() string {
	return err.err
}

//Character values 0-31 need to be escaped in query strings:
var controlCharEnd int = 31
var reservedChars = map[int]byte{
	36: '$',
	38: '&',
	43: '+',
	44: ',',
	47: '/',
	58: ':',
	59: ';',
	61: '=',
	63: '?',
	64: '@',
	12: '.',
}
var unsafeChars = map[int]int{
	32:  ' ',
	34:  '"',
	35:  '#',
	37:  '%',
	60:  '<',
	62:  '>',
	91:  '[',
	92:  '\\',
	93:  ']',
	94:  '^',
	96:  '`',
	123: '{',
	124: '|',
	125: '}',
	126: '~',
}
//Character values 128-255 need to be escaped.
var nonASCIImin int = 128
var nonASCIImax int = 255
