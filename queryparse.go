package main

import (
	_ "fmt"
	"log"
	"regexp"
	"strings"
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

//Build the timestring (index suffix)
func buildTimestring(t []string) string {
  //Get the length of the time arguments
  log.Println(t)

  return "TS"
}

//Build the Redis index using the operand supplied
//Return the correct Redis index
func buildIndex(oprnd string) {
	//Get the text from the inside of the parentheses
	re := regexp.MustCompile("\\(([^)]+)\\)")

	r := re.FindString(oprnd)
  //Remove the ( and ) from the string
  r = strings.Replace(r, "(", "", -1)
  r = strings.Replace(r, ")", "", -1)

  //Split the arguments of the operand (time arguments and index)
  args := strings.Split(r, ",")

  //Pass the time arguments to build the index suffix (exclude the index istelf)
  buildTimestring(args[0:len(args)-1])
}

//Perform the operation (interfaces with redis and will verify index existence and create new indexes)
// Returns the virtual index created
func performOp(optr string, oprnd1 string, oprnd2 string) string {
	log.Println("Evaluating: ", oprnd1, " ", optr, " ", oprnd2)

	buildIndex(oprnd1)

	return "XXXXX"
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
		//If the token is an operator, pop two items off the buffer and perform the operation
		if isOp(e[i]) {
			v := performOp(e[i], buf.Pop().(string), buf.Pop().(string))
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
