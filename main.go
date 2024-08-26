package main

import "fmt"

//ronald

// ro ron rona ronal
// ona onal onald
// na ald, ld
func main() {

	str := "ronald"

	for i := 0; i <= len(str); i++ {
		for j := 0; j <= i; j++ {
			fmt.Println(str[j:i])
		}
	}

}
