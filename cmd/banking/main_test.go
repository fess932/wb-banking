package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestState struct {
	account1, account2 string
	client             *http.Client
	server             *httptest.Server
}

func newTests() *TestState {
	return &TestState{client: &http.Client{Timeout: time.Second}, server: httptest.NewServer(rest())}
}

type JSON map[string]interface{}

func (j JSON) String() string {
	b, _ := json.Marshal(j)
	return string(b)
}

func Test_mainTest(t *testing.T) {
	ts := newTests()

	tests := []struct {
		name   string
		Body   interface{}
		Method string
		Path   func() string

		// response
		Status int
		Result fmt.Stringer
		Custom func(t *testing.T, r io.Reader) error
	}{
		{
			name: "регистрация",
			Body: JSON{
				"email":  "email1@email.ru",
				"amount": "123.45",
			},
			Method: http.MethodPost,
			Path:   ts.renderPath("/"),
			Custom: func(t *testing.T, r io.Reader) error {
				body := struct {
					Body struct {
						ID string `json:"id"`
					} `json:"body"`
				}{}

				if err := json.NewDecoder(r).Decode(&body); err != nil {
					require.NoError(t, err)
				}

				ts.account1 = body.Body.ID

				return nil
			},
		},
		{
			name: "регистрация 2",
			Body: map[string]interface{}{
				"email":  "email2@email.ru",
				"amount": "123.45",
			},
			Path: ts.renderPath("/"),
			Custom: func(t *testing.T, r io.Reader) error {
				body := struct {
					Body struct {
						ID string `json:"id"`
					} `json:"body"`
				}{}

				if err := json.NewDecoder(r).Decode(&body); err != nil {
					require.NoError(t, err)
				}

				ts.account2 = body.Body.ID

				return nil
			},
		},
		{
			name: "регистрация- ошибка аккаунт уже создан",
			Body: map[string]interface{}{
				"email":  "email2@email.ru",
				"amount": "123.45",
			},
			Status: http.StatusInternalServerError,
			Result: JSON{
				"error": "user existss",
			},
		},
		//{
		//	name: "amount",
		//	Body: map[string]interface{}{
		//		"email":  "email2@email.ru",
		//		"amount": "123.45",
		//	},
		//},
		//{
		//	name: "amount error: account not exists",
		//	Body: map[string]interface{}{
		//		"email":  "email2@email.ru",
		//		"amount": "123.45",
		//	},
		//},
		//
		//{
		//	name: "transfer error: account not exists",
		//},
		//{
		//	name: "transfer error: insufficient funds",
		//},
		//{
		//	name: "transfer ok",
		//},
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				err  error
				req  *http.Request
				resp *http.Response
			)

			switch bt := tt.Body.(type) {
			case JSON:
				data, err := json.Marshal(bt)
				require.NoError(t, err, "cant encode body to buffer")

				req, err = http.NewRequest(tt.Method, tt.Path(), bytes.NewBuffer(data))
				req.Header.Add("Content-Type", "application/json")
			case nil:
				req, err = http.NewRequest(tt.Method, tt.Path(), nil)
			}

			resp, err = ts.client.Do(req)
			require.NoError(t, err, "request error")
			defer resp.Body.Close()

			// result check
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			buf := bytes.NewBuffer(body)

			var result fmt.Stringer

			switch tt.Result.(type) {
			case JSON:
				var res JSON
				err = json.NewDecoder(buf).Decode(&res)
				result = res

			default:
				require.NoError(t, tt.Custom(t, buf))
				return
			}

			require.NoError(t, err, string(body))
			require.Equal(t, tt.Status, resp.StatusCode)

			exp := tt.Result.String()
			got := result.String()

			require.JSONEq(t, exp, got)
		})
	}
}

func (t *TestState) renderPath(staticPath string) func() string {
	return func() string {
		st := staticPath
		st = t.server.URL + st
		return st
	}
}
