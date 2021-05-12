package urlx

import "net/url"

type URL struct {
	*url.URL
}

func Parse(rawURL string) (*URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &URL{
		URL: parsedURL,
	}, nil
}
