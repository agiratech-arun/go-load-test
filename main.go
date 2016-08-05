package main

import (
  "gopkg.in/alecthomas/kingpin.v2"
  "github.com/agiratech-arun/go-load-test/jurniapi_v2_client"
  "os"
  )


var environment = kingpin.Flag("e","Specify an environment by default it is development").String()
var concurrency = kingpin.Flag("c","Specify an number of concurrency").Int()
var method = kingpin.Flag("method","Specify an Method Goging to be use").String()
// var username = kingpin.Flag("u","Specify username").String()


func main() {
  kingpin.Parse()
  video_path,_ := os.Getwd()
  video_path += "/videos/SampleVideo_1280x720_1mb.mp4"
  jurniapi_v2_client.StepUp(*environment, *concurrency,video_path,*method)
}


