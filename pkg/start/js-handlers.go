package start

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func configJs(env *envContext) (http.Handler, error) {
	buf := bytes.NewBufferString("")
	e := json.NewEncoder(buf)
	e.SetIndent("", "")
	if err := e.Encode(env.PublicEnv()); err != nil {
		return nil, err
	}
	jsonValue := strings.Replace(buf.String(), "\n", "", -1)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		_, _ = w.Write([]byte(fmt.Sprintf("window.Config = %s;", jsonValue)))
	}), nil
}
