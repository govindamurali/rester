package rester

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Requester struct {
	requestFunc RequestFunc
	url         string
	body        string
}

type RequestFunc func() (*http.Response, error)

var ErrRequestNotComplete = errors.New("request not complete")

func (r *Requester) Once() (*http.Response, error) {
	resp, err := r.requestFunc()
	return resp, err
}

func (r *Requester) WithConstantRetry(timeOut int) {
	retryFunc := func() error {
		res, err := r.requestFunc()
		if err != nil {
			return err
		}
		if err == ErrRequestNotComplete || HttpCodeForRetry(res.StatusCode) {
			return ErrRequestNotComplete
		}
		return nil
	}
	backoff.Retry(retryFunc, backoff.NewConstantBackOff(time.Duration(timeOut)*time.Second))
}

func (r *Requester) WithExponentialRetry(timeOut int) (respReturn *http.Response, err error) {
	reTryCount := -1
	retryFunc := func() error {
		reTryCount++
		res, err := r.requestFunc()
		if err == ErrRequestNotComplete {
			return ErrRequestNotComplete
		} else if HttpCodeForRetry(res.StatusCode) {
			err = errors.New(fmt.Sprintf("retriable resp status: %d", res.StatusCode))
			return ErrRequestNotComplete
		}

		respReturn, err = res, err
		return err
	}

	backOffPolicy := backoff.NewExponentialBackOff()
	backOffPolicy.MaxElapsedTime = time.Duration(timeOut) * time.Second
	backoff.Retry(retryFunc, backOffPolicy)
	return

}
