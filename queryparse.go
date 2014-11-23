package main

import (
	"log"
	"regexp"
	"strings"

	redis "github.com/xuyu/goredis"
)

//Extract set expressions from the query (virtual index builders)
func extractSetExpressions(q string) []string {
	re := regexp.MustCompile("\\[([^]]+)\\]")

	return re.FindAllString(q, -1)
}

//Is token operator
func isOp(o string) bool {
	if o == "AND" || o == "OR" {
		return true
	}
	return false
}

//Used to pad the timestring components (timestring builder helper function)
func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

//Build the timestring (index suffix)
func buildTimestring(t []string) string {
	//Get the length of the time arguments
	l := len(t)

	//The return string
	rtn := ""

	//The different time granularities
	year := ""
	month := ""
	day := ""
	hour := ""
	minute := ""
	second := ""

	switch l {
	//Year only
	case 1:
		year = t[0]
		//Check year is properly formatted
		if len(year) == 4 {
			rtn += year
		}
	//year month
	case 2:
		year = t[0]
		month = leftPad2Len(t[1], "0", 2)
		//Check year is properly formatted
		if len(year) == 4 {
			rtn += year
		}
		//Check the month is properly formatted
		if len(month) == 2 {
			rtn += month
		}
	//year month day
	case 3:
		year = t[0]
		month = leftPad2Len(t[1], "0", 2)
		day = leftPad2Len(t[2], "0", 2)
		//Check year is properly formatted
		if len(year) == 4 {
			rtn += year
		}
		//Check the month is properly formatted
		if len(month) == 2 {
			rtn += month
		}
		//Check the day is properly formatted
		if len(day) == 2 {
			rtn += day
		}
	//year month day hour
	case 4:
		year = t[0]
		month = leftPad2Len(t[1], "0", 2)
		day = leftPad2Len(t[2], "0", 2)
		hour = leftPad2Len(t[3], "0", 2)
		//Check year is properly formatted
		if len(year) == 4 {
			rtn += year
		}
		//Check the month is properly formatted
		if len(month) == 2 {
			rtn += month
		}
		//Check the day is properly formatted
		if len(day) == 2 {
			rtn += day
		}
		//Check the hour is properly formatted
		if len(hour) == 2 {
			rtn += hour
		}
	//year month day hour minute
	case 5:
		year = t[0]
		month = leftPad2Len(t[1], "0", 2)
		day = leftPad2Len(t[2], "0", 2)
		hour = leftPad2Len(t[3], "0", 2)
		minute = leftPad2Len(t[4], "0", 2)
		//Check year is properly formatted
		if len(year) == 4 {
			rtn += year
		}
		//Check the month is properly formatted
		if len(month) == 2 {
			rtn += month
		}
		//Check the day is properly formatted
		if len(day) == 2 {
			rtn += day
		}
		//Check the hour is properly formatted
		if len(hour) == 2 {
			rtn += hour
		}
		//Check the minute is properly formatted
		if len(minute) == 2 {
			rtn += minute
		}
	//year month day hour minute second
	case 6:
		year = t[0]
		month = leftPad2Len(t[1], "0", 2)
		day = leftPad2Len(t[2], "0", 2)
		hour = leftPad2Len(t[3], "0", 2)
		minute = leftPad2Len(t[4], "0", 2)
		second = leftPad2Len(t[5], "0", 2)
		//Check year is properly formatted
		if len(year) == 4 {
			rtn += year
		}
		//Check the month is properly formatted
		if len(month) == 2 {
			rtn += month
		}
		//Check the day is properly formatted
		if len(day) == 2 {
			rtn += day
		}
		//Check the hour is properly formatted
		if len(hour) == 2 {
			rtn += hour
		}
		//Check the minute is properly formatted
		if len(minute) == 2 {
			rtn += minute
		}
		//Check the second is properly formatted
		if len(second) == 2 {
			rtn += second
		}
	}
	return rtn
}

//Build the Redis index using the operand supplied
//Return the correct Redis index and if it is negated or not
func buildIndex(oprnd string) (string, bool) {
	//Get the text from the inside of the parentheses
	re := regexp.MustCompile("\\(([^)]+)\\)")

	r := re.FindString(oprnd)
	//Remove the ( and ) from the string
	r = strings.Replace(r, "(", "", -1)
	r = strings.Replace(r, ")", "", -1)

	//Split the arguments of the operand (time arguments and index)
	args := strings.Split(r, ",")

	//Timestring (will remain null if the index does not have a time function attached)
	ts := ""
	//Length of arguments must be greater than 1 (time function) (not true if virtual index)
	if len(args) > 1 {
		//Pass the time arguments to build the index suffix (exclude the index istelf)
		//Time string is in the format YYYYMMDDHHMMSS
		ts = buildTimestring(args[0 : len(args)-1])
	}
	//Check to see if the index is negated or not
	index := args[len(args)-1]
	log.Println(index)
	if string(index[0]) == "~" {
		//Remove the tilde
		index := strings.Replace(index, "~", "", -1)
		//Remove the single quotes
		index = strings.Replace(index, "'", "", -1)
		//Return the index with the timestring and negation set to true
		if ts != "" {
			return index + ":" + ts, true
		} else {
			return index, true
		}
	} else { //Not negated
		//Remove the single quotes
		index := strings.Replace(index, "'", "", -1)
		//Return the index with the timestring and negation set to false
		if ts != "" {
			return index + ":" + ts, false
			} else {
				return index, false
			}
	}
}

//Perform the binary operation (interfaces with redis and will verify index existence and create new indexes)
// Returns the virtual index created
func performBinOp(optr string, oprnd1 string, oprnd2 string) string {

	//Connect to server (TESTING)
	client, _ := redis.Dial(&redis.DialConfig{Address: "127.0.0.1:6379"})

	log.Println("Evaluating: ", oprnd1, " ", optr, " ", oprnd2)

	i1, n1 := buildIndex(oprnd1)
	i2, n2 := buildIndex(oprnd2)

	//Verify the existence of both indices, log fatal if not exists
	i1exists, _ := client.Exists(i1)
	i2exists, _ := client.Exists(i2)

	if !i1exists {
		log.Fatal("KEY DOES NOT EXIST: ", i1)
	}
	if !i2exists {
		log.Fatal("KEY DOES NOT EXIST: ", i2)
	}

	//If the first operand is negated, create a temporary (virtual) idnex
	if n1 {
		client.BitOp("NOT:" + i1, i1)
		i1 = "NOT:" + i1
	}
	//If the second operand is negated, create a temporary (virtual) idnex
	if n2 {
		client.BitOp("NOT:" + i2, i2)
		i2 = "NOT:" + i2
	}

	//Perform the operation
	client.BitOp(optr, i1 + optr + i2, i1, i2)

	log.Println("('" + i1 + optr + i2 + "')")
	return "('" + i1 + optr + i2 + "')"
}

//Evaluate a rpn epxression of indices and bitset operators
// Returns any virtual indexes created
func eval(rpn string) {
	//The buffer stack
	buf := new(Stack)
	//List of virtual indexes created
	//virt := []string{}

	//Break the expression into fields
	e := strings.Fields(rpn)

	//Iterate over the stack and evaluate using the buffer stack
	for i := range e {
		//If the token is an operator, pop two items off the buffer and perform the binary operation
		if isOp(e[i]) {
			v := performBinOp(e[i], buf.Pop().(string), buf.Pop().(string))
			buf.Push(v)
		} else {
			//If the token is not an operator, push onto the buffer
			buf.Push(e[i])
		}
	}
}

func Test() {

	query := "COUNT[( MONTH(2014,1,'users:Active') AND MONTH(2014,1,'users:inactive') ) OR ( HOUR(2014,1,12,10,'users:Testing') AND HOUR(2014,2,15,5,'users:Testing') ) ]"

	setExps := extractSetExpressions(query)

	//Remove the square brackets
	q := strings.Replace(setExps[0], "[", "", -1)
	q = strings.Replace(q, "]", "", -1)

	//Convert the expression to rpn (reverse polish notation)
	rpn := ParseInfix(q)

	//fmt.Println(rpn)

	//Evaluate the rpn expression, and return the virtual index created
	eval(rpn)
}
