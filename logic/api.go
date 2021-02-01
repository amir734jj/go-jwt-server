package logic

import "net/http"

func AuthorizedApi(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Successfully validated auth"))
	if err != nil {
		panic("Failed retuning result")
	}
}
