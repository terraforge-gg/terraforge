package utils

import "net/url"

func ExtractOriginAndPathFromUrl(rawUrl string) (origin string, pathname string, err error) {
	u, err := url.Parse(rawUrl)

	if err != nil {
		return "", "", err
	}

	return u.Scheme + "://" + u.Host, u.Path, nil
}
