package main

import (
	"content-parser/pkg/mr"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func main() {
	fmt.Println("Parser Worker Start")
	mapf := func(filename string, contents string) []mr.KeyValue {
		// function to detect word separators.
		ff := func(r rune) bool { return !unicode.IsLetter(r) }

		// split contents into an array of words.
		words := strings.FieldsFunc(contents, ff)

		kva := []mr.KeyValue{}
		for _, w := range words {
			kv := mr.KeyValue{w, "1"}
			kva = append(kva, kv)
		}
		return kva
	}
	reducef := func(key string, values []string) string {
		// return the number of occurrences of this word.
		return strconv.Itoa(len(values))
	}

	w := mr.MakeWorker(mapf, reducef)

	for {
		w.Run()
		time.Sleep(2 * time.Second)
	}

	// mr.Worker(mapf, reducef)

}
