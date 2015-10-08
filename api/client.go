package api

import
(
	"net/http"
	"log"
	"io/ioutil"
)

type Client struct {
	Num int
}

type Request struct {
	req      *http.Request
	logger   *log.Logger
}

type Response struct {
	resp     *http.Response
	logger   *log.Logger
}

//PutRequest provides request with PUT
func (c *Client) PutRequest(url string)*Request {
	return &Request{req: request("PUT", url)}
}

func (c *Client) GetRequest(url string)*Request {
	return &Request{req: request("GET", url)}
}

//Set provides append params to request
func (c *Request) Set(item, value string){
	//c.req.params.Set(item, value)
}

//Request provides basic request to api
func (c *Request) Do(url string, args interface{})*Response {
	client := &http.Client{}
    resp, err := client.Do(c.req)
    if err != nil {
        panic(err)
    }

    return &Response{resp: resp, logger: c.logger}
}


//IsOK return true is current request has a status = 200
func (c *Response) IsOK()bool {
	return c.resp.Status == "200"
}


func (c *Response) Result() string {
	body, err := ioutil.ReadAll(c.resp.Body)
	if err != nil {
		c.logger.Fatal("api.client: %v", err)
	}
	return string(body)
}

func (c *Request) Close() {
	c.req.Body.Close()
}

func request(method, url string)*http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil
	}

	return req
}