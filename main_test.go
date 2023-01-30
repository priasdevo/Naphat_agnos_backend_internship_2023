package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestMinStep(t *testing.T) {
	router := SetUpRouter()
	router.POST("/minStep", func(c *gin.Context) {
		// receive params from request
		var password postPass
		c.BindJSON(&password)
		// Check is password is in desired format (according to requirement)
		if !isValid(password.Pass) {
			c.JSON(400, gin.H{
				"status":  "fail",
				"message": "password is in a wrong format",
			})
		} else { // password is in correct format calculate min step to Strong password
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "The request was successful",
				"data": gin.H{
					"minStep": minStep(password.Pass),
				},
			})
		}
	})

	// Test 1 in strong-range
	Popassword := postPass{
		Pass: "qropkw",
	}
	jsonValue, _ := json.Marshal(Popassword)
	req, _ := http.NewRequest("POST", "http://localhost:8080/minStep", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	mockResponse := `{"data":{"minStep":2},"message":"The request was successful","status":"success"}`
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)

	// Test 2 shorter than strong pass
	Popassword = postPass{
		Pass: "qro",
	}
	jsonValue, _ = json.Marshal(Popassword)
	req, _ = http.NewRequest("POST", "http://localhost:8080/minStep", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	mockResponse = `{"data":{"minStep":3},"message":"The request was successful","status":"success"}`
	responseData, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)

	// Test 3 longer than strong pass and contain repeat
	Popassword = postPass{
		Pass: "qqqwlfpqoqk5781f2261",
	}
	jsonValue, _ = json.Marshal(Popassword)
	req, _ = http.NewRequest("POST", "http://localhost:8080/minStep", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	mockResponse = `{"data":{"minStep":2},"message":"The request was successful","status":"success"}`
	responseData, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)

	// Test 4 Wrong format
	Popassword = postPass{
		Pass: "qqq-lfpqoqk5781f2261",
	}
	jsonValue, _ = json.Marshal(Popassword)
	req, _ = http.NewRequest("POST", "http://localhost:8080/minStep", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	mockResponse = `{"message":"password is in a wrong format","status":"fail"}`
	responseData, _ = ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
