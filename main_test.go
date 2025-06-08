package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil) // функция принимает HTTP метод запроса, URL запроса и тело запроса

		handler.ServeHTTP(response, req) // Вызывается обработчик mainHandle(), в который передаётся response и req

		assert.Equal(t, http.StatusOK, response.Code) // Идёт проверка полученного в response ответа.
		fmt.Println(response.Body.String())
	}
}

func BadRoad(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	request := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}

	for _, v := range request {
		response := httptest.NewRecorder() // могу получить код ответа из поля Code (response.Code)
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	moscowCity := len(cafeList["moscow"])
	tulaCity := len(cafeList["tula"])

	requests := []struct {
		count int
		want  int
		city  string
	}{
		{0, 0, "moscow"},
		{1, 1, "moscow"},
		{2, 2, "moscow"},
		{100, moscowCity, "moscow"},
		{0, 0, "tula"},
		{1, 1, "tula"},
		{2, 2, "tula"},
		{100, tulaCity, "tula"},
	}

	for _, v := range requests {

		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?count=%d&city=%s", v.count, v.city)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		if response.Code != http.StatusOK {
			t.Errorf("не верный статус код, нынешний код %d", response.Code)
			continue
		}

		// обрабатываем ответ

		body := response.Body.String()
		var actualCount int

		if body == "" {
			actualCount = 0
		} else {
			cafes := strings.Split(body, ",")
			actualCount = len(cafes)
		}

		if actualCount != v.want {
			t.Errorf("колличество кафе %d, ожидаемый результат %d", actualCount, v.want)
		}

	}

}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range request {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		if response.Code != http.StatusOK {
			t.Errorf("не верный статус код, нынешний код %d", response.Code)
			continue
		}

		body := response.Body.String()
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",") // слайс строк
		}

		if len(cafes) != v.wantCount {
			t.Errorf("Не совпадение колличества кафе с искомым")
			continue
		}

		lowerSearch := strings.ToLower(v.search)
		for _, cafe := range cafes {
			lowerCafe := strings.ToLower(cafe)
			if !strings.Contains(lowerCafe, lowerSearch) {
				t.Errorf("не совпадения параметра search c искомым")
				continue
			}
		}
	}

}
