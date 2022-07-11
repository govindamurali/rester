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

func (r *Requester) RequestWithExponentialRetry() error {
	_, err := RequestWithBackoff(r)
	return err
}

func (r *Requester) Exponential() (*http.Response, error) {
	return RequestWithBackoff(r)
}

func (r *Requester) RequestOnce() (*http.Response, error) {
	resp, err := r.requestFunc()
	return resp, err
}

var ErrRequestNotComplete = errors.New("request not complete")

func RequestWithBackoff(requestOp *Requester) (respReturn *http.Response, errReturn error) {
	reTryCount := -1
	retryFunc := func() error {
		reTryCount++
		res, err := requestOp.requestFunc()
		if err == ErrRequestNotComplete {
			return ErrRequestNotComplete
		} else if HttpCodeForRetry(res.StatusCode) {
			err = errors.New(fmt.Sprintf("retriable resp status: %d", res.StatusCode))
			return ErrRequestNotComplete
		}

		respReturn, errReturn = res, err
		return err
	}

	backOffPolicy := backoff.NewExponentialBackOff()
	backOffPolicy.MaxElapsedTime = 15 * time.Minute
	backoff.Retry(retryFunc, backOffPolicy)
	return
}

func RequestWithRetry(requestOp RequestFunc, timeSecond int) {
	retryFunc := func() error {
		res, err := requestOp()
		if err != nil {
			return err
		}
		if err == ErrRequestNotComplete || HttpCodeForRetry(res.StatusCode) {
			return ErrRequestNotComplete
		}
		return nil
	}
	backoff.Retry(retryFunc, backoff.NewConstantBackOff(time.Duration(timeSecond)*time.Second))
}
