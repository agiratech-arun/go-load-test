package main

import (
  "gopkg.in/alecthomas/kingpin.v2"
  "github.com/agiratech-arun/go-load-test/jurniapi_v2_client"
  "os"
  )


var environment = kingpin.Flag("e","Specify an environment by default it is development").Required().String()
var concoreny = kingpin.Flag("c","Specify an number of concoreny").Required().Int()
var post_per_user = kingpin.Flag("p","Specify an number of post per users").Required().Int()


func main() {
  kingpin.Parse()
  pwd,_ := os.Getwd()
  pwd += "/videos/SampleVideo_1280x720_1mb.mp4"
  jurniapi_v2_client.StepUp(*environment, *concoreny, *post_per_user,pwd)
}


// func testGorotine(){
//   for i:=0;i<5;i++ {
//     wg.Add(1)
//     go startRequest()
//   }
//   wg.Wait()
// }

// func startRequest() {
//   defer wg.Done()
//   fmt.Println("Start request",time.Now())
//   response, err :=http.Get("http://golang.org")
//   fmt.Println(response,err)
//   fmt.Println("end request")
// }

