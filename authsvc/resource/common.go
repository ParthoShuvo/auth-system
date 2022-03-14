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

func reqmuxb(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return ioutil.ReadAll(r.Body)
}

func unmarshall(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
