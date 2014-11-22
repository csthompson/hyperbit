package main

import (
	"fmt"
	_"log"
	"regexp"
	"strings"
)

//The query operation table
var optab = map[string]func(ideces []string){
	// Unary operator. Perform a bit-wise count on a single index.
	"COUNT": func(indices []string) {},
	// Binary operator. Perform a bit-wise AND operation on two indices.
	"AND": func(indices []string) {},
	// Binary operator. Perform a bit-wise OR operation on two indices.
	"OR": func(indices []string) {},
}

//Extract set expressions from the query (virtual index builders)
func extractSetExpressions(q string) []string {
	re := regexp.MustCompile("\\[([^]]+)\\]")

	return re.FindAllString(q, -1)
}

func Test() {

  query := "COUNT[( 'users:Active' AND 'users:inactive' ) OR ( 'users:Testing' AND 'users:Testing' ) ]"

  setExps := extractSetExpressions(query)

  //Remove the square brackets
  q := strings.Replace(setExps[0], "[", "", -1)
  q = strings.Replace(q, "]", "", -1)

  fmt.Println("postfix:", ParseInfix(q))
}
