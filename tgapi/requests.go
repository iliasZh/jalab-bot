package tgapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"

	"jalabs.kz/bot/model"
)

func doRequest[Rq any, Rs any](
	ctx context.Context,
	t *TgAPI,
	endpoint string,
	rq Rq,
) (rs Rs, err error) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("doing '%s' request: %w", endpoint, err)
	}()

	botToken := fmt.Sprintf("bot%s", t.botToken)
	endpointURL := &url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   path.Join(botToken, endpoint),
	}

	rqBody, errNewRqBody := newRequestBody[Rq](rq)
	if errNewRqBody != nil {
		return rs, errNewRqBody
	}

	method := http.MethodPost
	if rqBody == nil {
		method = http.MethodGet
	}

	httpRq, errNewHTTPRq := http.NewRequestWithContext(ctx, method, endpointURL.String(), rqBody)
	if errNewHTTPRq != nil {
		return rs, fmt.Errorf("new http.Request: %w", errNewHTTPRq)
	}

	httpRq.Header.Set("Content-Type", "application/json")

	httpRs, errDo := t.cli.Do(httpRq)
	if errDo != nil {
		return rs, fmt.Errorf("doing request: %w", errDo)
	}

	rsBytes, errRead := readResponse(httpRs)
	if errRead != nil {
		return rs, errRead
	}

	fullRs := model.Response[Rs]{}
	if errUnmarshal := json.Unmarshal(rsBytes, &fullRs); errUnmarshal != nil {
		return rs, fmt.Errorf("unmarshaling response: %w", errUnmarshal)
	}

	if !fullRs.Ok {
		return rs, fmt.Errorf("failed with status %d: %s", fullRs.ErrorCode, fullRs.Description)
	}

	return fullRs.Result, nil
}

func readResponse(rs *http.Response) ([]byte, error) {
	contentLenStr := rs.Header.Get("Content-Length")

	var contentLen int
	if contentLenStr != "" {
		contentLen64, errParse := strconv.ParseInt(contentLenStr, 10, 32)
		if errParse != nil {
			return nil, fmt.Errorf("parsing content length: %w", errParse)
		}

		contentLen = int(contentLen64)
	}

	buf := &bytes.Buffer{}
	buf.Grow(bytes.MinRead + contentLen)
	_, errRead := buf.ReadFrom(rs.Body)
	if errRead != nil {
		return nil, fmt.Errorf("reading response body: %w", errRead)
	}

	return buf.Bytes(), nil
}

func newRequestBody[Rq any](rq Rq) (io.Reader, error) {
	if reflect.TypeOf(rq) == nil {
		return nil, nil
	}

	rqBytes, errMarshal := json.Marshal(rq)
	if errMarshal != nil {
		return nil, fmt.Errorf("marshaling request: %w", errMarshal)
	}

	return bytes.NewReader(rqBytes), nil
}
