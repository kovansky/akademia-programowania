package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const urlFmt = "https://www.reddit.com/r/%s.json"

type MultiError struct {
	errors []error
}

func (m *MultiError) Error() string {
	var sb strings.Builder

	for errId, err := range m.errors {
		sb.WriteString(fmt.Sprintf("%d:\n%s\n", errId, err.Error()))
	}

	return sb.String()
}

func (m *MultiError) Append(err error) {
	m.errors = append(m.errors, err)
}

type RedditFetcherImpl struct {
	Subreddits []string
	data       map[string]response
}

func (reddit *RedditFetcherImpl) Fetch(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(len(reddit.Subreddits))

	reddit.data = make(map[string]response, len(reddit.Subreddits))

	var errs *MultiError

	for _, sub := range reddit.Subreddits {
		go func(sub string) {
			err := reddit.fetchSubreddit(ctx, &wg, sub)
			if err != nil {
				if errs == nil {
					errs = &MultiError{}
				}

				errs.Append(err)
			}
		}(sub)
	}

	wg.Wait()

	if errs != nil {
		return errs
	}
	return nil
}

func (reddit *RedditFetcherImpl) fetchSubreddit(ctx context.Context, wg *sync.WaitGroup, sub string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
		wg.Done()
	}()

	logger := ctx.Value(CtxLoggerKey).(*log.Logger)
	logPrefix := color.New(color.Bold, color.FgRed).Sprintf("[/r/%s]", sub)
	logger.Printf("%s %s", logPrefix, color.HiRedString("Starting fetch..."))

	req, err := http.NewRequestWithContext(timeoutCtx, http.MethodGet, fmt.Sprintf(urlFmt, sub), http.NoBody)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Reddit returned status code %d", res.StatusCode))
	}

	var parsed response
	err = json.NewDecoder(res.Body).Decode(&parsed)
	if err != nil {
		return err
	}

	reddit.data[sub] = parsed

	logger.Printf("%s %s", logPrefix, color.HiRedString("Data fetch finished"))
	wg.Done()

	return nil
}

func (reddit *RedditFetcherImpl) Save(writer io.Writer) error {
	for sub, data := range reddit.data {
		if _, err := writer.Write([]byte(fmt.Sprintf("--== /r/%s ==--\n", sub))); err != nil {
			return err
		}

		for _, post := range data.Data.Children {
			if _, err := writer.Write([]byte(fmt.Sprintf("\t%s\n%s\n", post.Data.Title, post.Data.URL))); err != nil {
				return err
			}
		}
	}

	return nil
}
