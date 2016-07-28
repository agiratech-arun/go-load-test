package main

import (
  "fmt"
  // "runtime"
  "sync"
  "time"
  "net/http"
  "gopkg.in/alecthomas/kingpin.v2"
  "github.com/agiratech-arun/go-load-test/jurniapi_v2_client"
  )

var wg sync.WaitGroup

var environment = kingpin.Flag("e","Specify an environment by default it is development").Required().String()
var concoreny = kingpin.Flag("c","Specify an number of concoreny").Required().Int()


func main() {
  kingpin.Parse()
  testGorotine()
}


func testGorotine(){
  for i:=0;i<5;i++ {
    wg.Add(1)
    go startRequest()
  }
  wg.Wait()
}

func startRequest() {
  defer wg.Done()
  fmt.Println("Start request",time.Now())
  response, err :=http.Get("http://golang.org")
  fmt.Println(response,err)
  fmt.Println("end request")
}

