package main

import (
	"database/sql"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// check if password is valid according to given format
// 1<= len(password) <=40
// character is letter,digit,".","!"
func isValid(password string) bool {
	if len(password) == 0 || len(password) > 40 {
		return false
	}
	for _, char := range password {
		if !(('a' < char && char < 'z') || ('A' < char && char < 'Z') || ('0' < char && char < '9') || (char == '!') || (char == '.')) {
			return false
		}
	}
	return true
}
func minStep(password string) int {
	var lenght int = len(password)
	// Considered 3 case seperately 1. lenght < 6, 2. lenght between 6-20, 3. lenght >= 20
	// 3 requirement
	var isUpper, isLower, isDigit int = 1, 1, 1
	var totalRequire int = 0
	// variable use to deal with repeated in a row char
	var current_repeated rune = '?'
	var repeat_count int = 0
	var action_require int = 0
	// variable use to deal with delete
	var repeat_cdel []int
	var over bool = false
	// number of delete require
	var n int = lenght - 19
	// Process that is need for all case
	for _, char := range password {
		// Check isUpper,lower,Digit
		if char >= 'A' && char <= 'Z' {
			isUpper = 0
		} else if char >= 'a' && char <= 'z' {
			isLower = 0
		} else if char >= '0' && char <= '9' {
			isDigit = 0
		}
		// Calculate value for further processing on repeated in a row
		// Check if current char is the same with previous ine
		if char == current_repeated {
			// if yes +1 to count
			repeat_count += 1
			// if repeat = 3 change the character to something else
			// add action count and reset repeat count
			if repeat_count == 3 {
				action_require += 1
				repeat_count = 0
				// set over repeat = true => must be delete from the right-hand of changed character
				over = true
			}
		} else if char != current_repeated {
			if over {
				over = false
				repeat_cdel = append(repeat_cdel, repeat_count)
			}
			repeat_count = 1
			current_repeated = char
		}
	}
	totalRequire = isDigit + isLower + isUpper
	// Case 1 need to fill to 6 character
	// Consider 1 lower, 1 upper, 1 digit requirement
	//  '.' and '!' is not required
	// Step need is max(sum of requirement lack, Character lack)
	// we can alway find the character to fill in so than no 3 repeating character in a row so not consider this requirement
	if lenght < 6 {
		if 6-lenght < totalRequire {
			return totalRequire
		} else {
			return 6 - lenght
		}
	} else if lenght < 20 {
		// Case 2 doesn't need to fill or delete any character but might need to change character to smooth out 3 repeated in row
		// Character change can also use to deal with no Digit/lower/upper requirement
		// So this 2 requirement can be use as overlap requirement
		if action_require < totalRequire {
			return totalRequire
		} else {
			return action_require
		}
	} else { // lenght >= 20
		// Case 3 need to delete some character
		// Problem : delete may overlap / not overlap with other action
		// delete might make some unrepeated in the row become repeated in the row
		// Fix : delete only right-most/left-most of the string
		// Fix : delete only the Character which is the source of repeated already
		// if delete the source of repeated we can be sure that Digit/lower/upper will still be in the password
		// if delete right-most/left-most it might make 3 condition not meet anymore

		// Calculate step in-case no delete need
		if action_require < totalRequire {
			action_require = totalRequire
		}
		// sort repeat_cdel (delete changed character of repeated in row without same character in the right use less action)
		sort.Sort(sort.IntSlice(repeat_cdel))
		// repeat through all character that has been change due to repeat
		for ite, num := range repeat_cdel {
			if num+1 <= n {
				// in-case the change is overlap with islower/upper/digit requirement you can't delete it
				if totalRequire != 0 && len(repeat_cdel)-ite <= totalRequire {
					totalRequire -= 1
					n += 1
				}
				n = n - num - 1
				action_require += num
			} else {
				// there are right repeat character more than need to delete, delete only need
				action_require += n
				n = 0
				break
			}
		}
		if n == 0 {
			// already strong
			return action_require
		} else {
			// delete more character (no more overlap, can alway find character to delete that can prevent 3 in a row)
			return action_require + n
		}
	}
}

// JSON POST receive
// password/dbname need to change
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "agnos_assign"
)

func main() {
	// connect database
	db, err := sql.Open("postgres", "user=postgres password=qduZ9Cg2D.gWU.m dbname=agnos_assign sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// init router
	router := gin.Default()

	// Log to psql
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		status := param.StatusCode
		latency := param.Latency.String()
		clientIP := param.ClientIP
		method := param.Method
		path := param.Path

		_, err := db.Exec("INSERT INTO logs (time, status, latency, client_ip, method, path) VALUES ($1, $2, $3, $4, $5, $6)",
			time.Now(), status, latency, clientIP, method, path)
		if err != nil {
			panic(err)
		}
		return ""
	}))

	// GET requset for calculate minStep
	router.POST("/minStep", func(c *gin.Context) {
		// receive params from request
		password := c.PostForm("pass")
		// Check is password is in desired format (according to requirement)
		if !isValid(password) {
			c.JSON(400, gin.H{
				"status":  "fail",
				"message": "password is in a wrong format",
			})
		} else { // password is in correct format calculate min step to Strong password
			c.JSON(200, gin.H{
				"status":  "success",
				"message": "The request was successful",
				"data": gin.H{
					"minStep": minStep(password),
				},
			})
		}
	})

	// start router
	router.Run("localhost:8080")
}
