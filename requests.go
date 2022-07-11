package rester

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func PostRequest(url string, request interface{}, conversionType interface{}, headers map[string]string, customTransport http.RoundTripper) Requester {
	bty, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	r := func() (*http.Response, error) {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bty))
		req.Header.Add("content-type", "application/json")

		for key, value := range headers {
			req.Header.Add(key, value)
		}

		client := getHttpClient(customTransport)

		res, err := client.Do(req)
		if err != nil {
			return res, ErrRequestNotComplete
		}

		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&conversionType)

		return res, err
	}
	return Requester{r, url, string(bty)}
}

func PutRequest(url string, request interface{}, conversionType interface{}, headers map[string]string, customTransport http.RoundTripper) Requester {
	return GenericRequest(url, request, conversionType, headers, "PUT", customTransport)
}

func DeleteRequest(url string, request interface{}, conversionType interface{}, headers map[string]string, customTransport http.RoundTripper) Requester {
	return GenericRequest(url, request, conversionType, headers, "DELETE", customTransport)
}

func GetRequest(url string, conversionType interface{}, headers map[string]string, customTransport http.RoundTripper) Requester {

	r := func() (*http.Response, error) {
		req, _ := http.NewRequest("GET", url, nil)
		for key, value := range headers {
			req.Header.Add(key, value)
		}

		client := getHttpClient(customTransport)

		res, err := client.Do(req)
		if err != nil {
			return res, ErrRequestNotComplete
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&conversionType)

		return res, err
	}
	return Requester{r, url, ""}
}

func GenericRequest(url string, request interface{}, conversionType interface{}, headers map[string]string, requestMethod string, customTransport http.RoundTripper) Requester {

	bty, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	r := func() (*http.Response, error) {
		req, _ := http.NewRequest(requestMethod, url, bytes.NewBuffer(bty))
		req.Header.Add("content-type", "application/json")

		for key, value := range headers {
			req.Header.Add(key, value)
		}

		client := getHttpClient(customTransport)

		res, err := client.Do(req)
		if err != nil {
			return res, ErrRequestNotComplete
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&conversionType)

		return res, err
	}
	return Requester{r, url, string(bty)}
}

func getHttpClient(customTransport http.RoundTripper) (client *http.Client) {
	client = &http.Client{}
	if customTransport != nil {
		client.Transport = customTransport
	} else {
		client.Transport = http.DefaultTransport
	}
	return
}
