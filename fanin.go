//this programme is a fanin exapmle which is going to read data from two external files
//and merge the output into a single channel and then
//print the output

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func readData(file string) <-chan string {
	f, err := os.Open(file) //opens the file for reading
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string) //channel declared

	//returns a scanner to read from f
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines) //scanning it line-by-line token

	//loop through the fileScanner based on our token split
	go func() {
		for fileScanner.Scan() {
			val := fileScanner.Text() //returns the recent token
			out <- val                //passed the token value to our channel
		}

		close(out) //closed the channel when all content of file is read

		//closed the file
		err := f.Close()
		if err != nil {
			fmt.Printf("Unable to close an opened file: %v\n", err.Error())
			return
		}
	}()

	return out
}

func fanInMergeData(inputs ...<-chan string) chan string {
	chRes := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	for _, in := range inputs {
		go func(ch <-chan string) {
			for n := range ch {
				chRes <- n
			}
			wg.Done()

		}(in)
	}

	go func() {
		wg.Wait()    //waits till the goroutines are completed and wg marked Done
		close(chRes) //close the result channel
	}()

	return chRes
}

func main() {
	ch1 := readData("D:\\assignment_Saisandeep\\text1.txt")
	ch2 := readData("D:\\assignment_Saisandeep\\text2.txt")

	chRes := fanInMergeData(ch1, ch2)


	for val := range chRes {
		fmt.Println(val)
	}
}