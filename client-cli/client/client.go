package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func postWrapper(url string, body string) (string, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func New() Client {
	var client Client
	client.URL = "http://localhost:8080/api"

	return client
}

type Client struct {
	URL string
}

func (c *Client) StepInitialize() (string, error) {
	result, err := postWrapper(c.URL+"/step/initialize", fmt.Sprintf("{}"))
	return result, err
}

func (c *Client) StepIdempotentApply(transactionName string) (string, error) {
	result, err := postWrapper(c.URL+"/step/idempotent-apply", fmt.Sprintf("{\"transaction\": \"%s\"}", transactionName))
	return result, err
}

func (c *Client) StepDiff(transactionName string) (string, error) {
	result, err := postWrapper(c.URL+"/step/diff", fmt.Sprintf("{\"transaction\": \"%s\"}", transactionName))
	return result, err
}

func (c *Client) TransformerAdd(transactionName string, typ string, transformer string, content string) (string, error) {
	result, err := postWrapper(c.URL+"/transformer/add", fmt.Sprintf("{\"transaction\": \"%s\", \"type\": \"%s\", \"transformer\": \"%s\", \"content\": \"%s\"}", transactionName, typ, transformer, content))
	return result, err
}

func (c *Client) TransformerRemove(transactionName string, transformer string) (string, error) {
	result, err := postWrapper(c.URL+"/transformer/remove", fmt.Sprintf("{\"transaction\": \"%s\", \"transformer\": \"%s\"}", transactionName, transformer))
	return result, err
}

func (c *Client) TransformerEdit(transactionName string, transformer string, newContent string) (string, error) {
	result, err := postWrapper(c.URL+"/transformer/edit", fmt.Sprintf("{\"transaction\": \"%s\", \"transformer\": \"%s\", \"newContent\": \"%s\"}", transactionName, transformer, newContent))
	return result, err
}

func (c *Client) TransformerOrder(transactionName string, order string) (string, error) {
	result, err := postWrapper(c.URL+"/transformer/order", fmt.Sprintf("{\"transaction\": \"%s\", \"order\": \"%s\"}", transactionName, order))
	return result, err
}

func (c *Client) RepoAdd(transaction string, repo string) (string, error) {
	result, err := postWrapper(c.URL+"/repo/add", fmt.Sprintf("{\"transaction\": \"%s\", \"repo\": \"%s\"}", transaction, repo))
	return result, err
}

func (c *Client) RepoRemove(transaction string, repo string) (string, error) {
	result, err := postWrapper(c.URL+"/repo/remove", fmt.Sprintf("{\"transaction\": \"%s\", \"repo\": \"%s\"}", transaction, repo))
	return result, err
}

func (c *Client) TransactionGet(name string) (string, error) {
	result, err := postWrapper(c.URL+"/transaction/get", fmt.Sprintf("{ \"name\": \"%s\" }", name))
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Client) TransactionAdd(name string) error {
	_, err := postWrapper(c.URL+"/transaction/add", fmt.Sprintf("{ \"name\": \"%s\" }", name))
	return err
}

func (c *Client) TransactionRemove(name string) error {
	_, err := postWrapper(c.URL+"/transaction/remove", fmt.Sprintf("{ \"name\": \"%s\" }", name))
	return err
}

func (c *Client) TransactionRename(oldName string, newName string) error {
	_, err := postWrapper(c.URL+"/transaction/rename", fmt.Sprintf("{ \"oldName\": \"%s\", \"newName\": \"%s\" }", oldName, newName))
	return err
}

func (c *Client) TransactionList() (string, error) {
	result, err := postWrapper(c.URL+"/transaction/list", "{}")
	return result, err
}
