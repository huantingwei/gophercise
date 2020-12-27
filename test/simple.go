package test

import (
	"io/ioutil"
	"net/http"
)

func first(num int) []int {

	return make([]int, num)
}

func second(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
