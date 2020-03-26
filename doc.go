/*
Package pinky implements JavaScript style promises to Golang

Installation

To download the the latest, run:

	go get github.com/zainkai/pinky@latest

Import it in your program as:

	import "github.com/zainkai/pinky"


Usage

Promises can be chained together

	customErr := errors.New("my promise rejected")

	sum := make(chan int, 1)
	go NewPromise(0).Then(func(value interface{}, resolve ResolveFunc, reject RejectFunc) {
		i, _ := value.(int)
		fmt.Println("Adding one.")

		resolve(i + 1)
	}).Then(func(value interface{}) (interface{}, error) {
		i, _ := value.(int)
		fmt.Println("Adding two.")

		return i + 2, nil
	}).Then(func(value interface{}) (interface{}, error) {
		fmt.Println("Sending three to sum channel.")
		i, ok := value.(int)
		if !ok {
			return nil, errors.New("could no type change to int")
		}
		return i, nil
	}).CatchCase(customErr, func(err error) {
		fmt.Println("err: ", err)
	}).Catch(func(err error) {
		fmt.Println("unknown error: ", err)
	}).Finally(func(value interface{}, _ error) {
		i, _ := value.(int)
		sum <- i
	})

	fmt.Println("Recieved Value: ", <-sum)
*/
package pinky
