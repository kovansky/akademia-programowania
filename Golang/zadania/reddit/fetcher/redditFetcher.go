package fetcher

import (
	"context"
	"io"
)

const CtxLoggerKey = "RDF_LOGGER"

type response struct {
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditFetcher interface {
	Fetch(context.Context) error
	Save(io.Writer) error
}
