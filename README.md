# rester
Library for simple **REST** calls **GET**, **POST**, **PUT**, **DELETE** without worrying about nitty gritties. You could chose to make the calls with periodic retry, exponential backoff, or just once. 

# Install
go get github.com/govindamurali/rester

# Usage

### Initialize request and response structs
```
  type Request struct{
    SomeReqData string `json:"some_req_data"`
  }
  
  type Response struct{
    SomeRespData string `json:"some_resp_data"`
  }

```

### Initialize request and response variables
```
  request:= Request{SomeReqData: "some value"}
  response:= Response{}	
```

### Initialize request methods
```
  postRequest := rester.PostRequest(url, request, &response, customHeaders, customTransport)
  getRequest := rester.GetRequest(url, &response, nil, nil)    
  putRequest :=  rester.PutRequest(url, request, &response, nil, nil)   
  deleteRequest := rester.DeleteRequest(url, request, &response, customHeaders, customTransport)
```

### Make the call
```
  httpResponse,err := postRequest.Once()
  httpResponse,err := getRequest.WithConstantRetry(maxTimeOutInSeconds)
  httpResponse,err := putRequest.WithExponentialRetry(maxTimeOutInSeconds)
```

### Miscellaneous
* You'll be able to get the response in the **response** variable & httpStatus in **httpResponse** variable
* **customHeaders** are optional. Will use defaultHeaders when nil
* customTransport can be used for tracking. eg. Use newrelic agent transport to get metrics tracked on newrelic. 
