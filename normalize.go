/* 
 * Package normalize normalizes URLs
 */
package normalize

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

var ipv6Regexp, _ = regexp.Compile("([0-9A-F]+:)+[0-9A-F]+")
var ipv4Regexp, _ = regexp.Compile("([0-9]+[.])+[0-9]+")

func Normalize(url *url.URL) {
	addSlash(url)
	removeDefaultPort(url)
	lowerCaseScheme(url)
	lowerCaseDomain(url)
	removeDoubleSlashes(url)
	removeDirectoryDots(url)
}

func removeDirectoryDots(url *url.URL) {
	url.Path = strings.Replace(url.Path, "/./", "/", -1)
	if strings.Contains(url.Path, "..") {
		seg := strings.Split(url.Path, "/")
		stk := NewStack(1000)
		for _, s := range seg {
			if len(s) == 0 {
				continue
			}
			if strings.EqualFold(s, string("..")) {
				stk.Pop()
			} else {
				stk.Push(s)
			}
		}
		if stk.Top > 0 {
			r := stk.Data[0:stk.Top]
			url.Path = fmt.Sprintf("/%s/", strings.Join(r, "/"))
		}
	}
}

func removeDoubleSlashes(url *url.URL) {
	url.Path = strings.Replace(url.Path, "//", "/", -1)
}

func lowerCaseDomain(url *url.URL) {
	url.Host = strings.ToLower(url.Host)
}

func lowerCaseScheme(url *url.URL) {
	url.Scheme = strings.ToLower(url.Scheme)
}

func removeDefaultPort(url *url.URL) {
	/*
	 * Have to ensure that not removing the last part of an ipv6
	 * address if it happens to be :80 as unlikely as that may be.
	 */
	host := url.Host
	if host[len(host)-3:] == ":80" {
		if found := ipv6Regexp.FindStringIndex(url.Host); found != nil {
			if strings.Count(host, ":") < 8 {
				return
			}
		}
		url.Host = host[:len(host)-3]
	}
}

func addSlash(url *url.URL) {
	if len(url.Path) == 0 {
		url.Path = "/"
	}
}

/*
 * Removes directory indexes when they point to the same place as
 * the directory. For example if index.html points to / and
 * index.html is given for the index parameter it will be removed
 * from the URL
 */
func RemoveDirectoryIndex(url *url.URL, index string) {
	pathLen := len(url.Path)
	indexLen := len(index)
	if pathLen >= indexLen {
		if url.Path[pathLen-indexLen:] == index {
			url.Path = url.Path[:pathLen-indexLen]
		}
	}
}

/*
 * Ordes query variables in alphabetic order. Order of variables
 * in a query string should not matter, but some implementations
 * may require an order, so this is in a separate emthod.
 */
func NormalizeQueryVariableOrder(url *url.URL) {
	keys := []string{}
	values := url.Query()
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	variables := []string{}
	for _, key := range keys {
		for _, value := range values[key] {
			if len(value) > 0 {
				variables = append(variables, key+"="+value)
			} else {
				variables = append(variables, key)
			}
		}
	}
	url.RawQuery = strings.Join(variables, "&")
}

/*
 * Remove query variables that have default values.  Provide a set of defaults
 * (defaults[key] = value) wher key is the variable name and value is the string
 * represenation of the default value.
 */
func RemoveDefaultQueryValues(url *url.URL, defaults map[string]string) {
	keys := []string{}
	values := url.Query()
	for k, _ := range values {
		keys = append(keys, k)
	}
	variables := []string{}
	for _, key := range keys {
		defaultValue, ok := defaults[key]
		for _, value := range values[key] {
			if len(value) > 0 {
				if !ok || value != defaultValue {
					variables = append(variables, key+"="+value)
				}
			} else {
				variables = append(variables, key)
			}
		}
	}
	url.RawQuery = strings.Join(variables, "&")
}

/*
 * Removes www. from a URL. Use if www. points to same resource as
 * non-www address.
 */
func NormalizeWWW(url *url.URL, showWWW bool) {
	var foundWWW bool
	if found := ipv6Regexp.FindStringIndex(url.Host); found != nil {
		return
	}
	if found := ipv4Regexp.FindStringIndex(url.Host); found != nil {
		return
	}
	if len(url.Host) <= 4 {
		foundWWW = false
	} else {
		foundWWW = url.Host[:4] == "www."
	}
	if showWWW && !foundWWW {
		url.Host = "www." + url.Host
	} else if !showWWW && foundWWW {
		url.Host = url.Host[4:]
	}
}

/*
 * Remove arbitary query variables. Include a slice of array variables
 * to check against. If query variables are found not in the given slice,
 * they are removed.
 */
func NormalizeQuery(url *url.URL, params []string) {
	keys := []string{}
	values := url.Query()
	for k, _ := range values {
		keys = append(keys, k)
	}
	variables := []string{}
	for _, key := range keys {
		for _, expectedParam := range params {
			if expectedParam == key {
				for _, value := range values[key] {
					if len(value) > 0 {
						variables = append(variables, key+"="+value)
					} else {
						variables = append(variables, key)
					}
				}
			}
		}
	}
	url.RawQuery = strings.Join(variables, "&")
}

/*
 * Normalize scheme or protocol. For example if valid scheme is url
 * and not urls, url is changed to url if 'url' is given as scheme.
 */
func NormalizeScheme(url *url.URL, scheme string) {
	url.Scheme = scheme
}

/*
 * Remove #fragment from a URL.
 */
func RemoveFragment(url *url.URL) {
	url.Fragment = ""
}

/*
 * Replaces domain or IP with given domain. Use to replace IP addresses with
 * domain or domains that point to the same resource as prime domain.
 */
func NormalizeDomain(url *url.URL, domain string) {
	url.Host = domain
}
