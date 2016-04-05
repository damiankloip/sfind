package main

import (
  "fmt"
  "sync"
)

type Result interface {
  beforeResults()
  eachResult(matched string)
  afterResults()
}

type PrintResult struct {}

func (r *PrintResult) beforeResults() {

}

func (r *PrintResult) eachResult(matched string) {
  // Add a null byte to the end of each match string.
  fmt.Println(matched + "\x00")
}

func (r *PrintResult) afterResults() {

}

type CountResult struct {
  mutex *sync.Mutex
  count int
}

func (r *CountResult) beforeResults() {

}

func (r *CountResult) eachResult(matched string) {
  // Increment the matched results by one.
  r.mutex.Lock()
  r.count += 1
  r.mutex.Unlock()
}

func (r *CountResult) afterResults() {
  // Print the count.
  fmt.Println(r.count)
}
