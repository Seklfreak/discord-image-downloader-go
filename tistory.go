package main

import (
	"net/url"
	"strings"
)

// getTistoryUrls downloads tistory URLs
// http://t1.daumcdn.net/cfile/tistory/[…] => http://t1.daumcdn.net/cfile/tistory/[…]
// http://t1.daumcdn.net/cfile/tistory/[…]?original => as is
func getTistoryUrls(link string) (map[string]string, error) {
	if !strings.HasSuffix(link, "?original") {
		link += "?original"
	}
	return map[string]string{link: ""}, nil
}

func getLegacyTistoryUrls(link string) (map[string]string, error) {
	link = strings.Replace(link, "/image/", "/original/", -1)
	return map[string]string{link: ""}, nil
}

func getTistoryWithCDNUrls(urlI string) (map[string]string, error) {
	parameters, _ := url.ParseQuery(urlI)
	if val, ok := parameters["fname"]; ok {
		if len(val) > 0 {
			if RegexpUrlTistoryLegacy.MatchString(val[0]) {
				return getLegacyTistoryUrls(val[0])
			}
		}
	}
	return nil, nil
}
