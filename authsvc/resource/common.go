package resource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	usrTable "github.com/parthoshuvo/authsvc/table/user"
)

type LoginUser struct {
	Email    usrTable.Email    `json:"email" validate:"required,email"`
	Password usrTable.Password `json:"password" validate:"required,validPwd,min=8,max=64"`
}

func (lusr *LoginUser) isAuthenticated(password usrTable.Password) bool {
	return lusr.Password.Hash().Equals(password)
}

type wrapper struct {
	req *http.Request
}

// requestWrapper creates a HTTP request wrapper
func requestWrapper(r *http.Request) *wrapper {
	return &wrapper{r}
}

func (w *wrapper) body() ([]byte, error) {
	return reqmuxb(w.req)
}

func (w *wrapper) host() string {
	return w.req.Host
}

func (w *wrapper) email() string {
	return reqmuxq(w.req, "email")
}

func (w *wrapper) verificationCode() string {
	return reqmuxq(w.req, "verification_code")
}

func (w *wrapper) loginUser() (*LoginUser, error) {
	data, err := w.body()
	if err != nil {
		return nil, err
	}
	lusr := LoginUser{}
	err = unmarshall(data, &lusr)
	if err != nil {
		return nil, err
	}
	return &lusr, nil
}

// func (w *wrapper) bearerAuth() (string, error) {
// 	header := w.req.Header.Get("Authorization")
// 	if header == "" {
// 		return "", fmt.Errorf("missing authorization header")
// 	}
// 	const authScheme = "Bearer"
// 	if !strings.HasPrefix(header, authScheme) {
// 		return "", fmt.Errorf("missing bearer auth scheme at authorization header")
// 	}
// 	return header[len(authScheme)+1:], nil
// }

func reqmuxq(r *http.Request, name string) string {
	return r.URL.Query().Get(name)
}

func reqmuxb(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func unmarshall(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
