package jurniapi_v2_client

import (
        "fmt"
        "net/http"
       )

type Config struct {
  EnvVariable String
  Concurrency int
  EnvConvig EnvSetup
}


type EnvSetup struct {
  BaseUri string
  SSLCaFile string
}

type RequestSetup struct {
  Url string
  Params map[string]string
  Headers map[string]string
}



func (c *Config) ConfigSetup() {
  if c.EnvVariable == "staging" {
    c.EnvConvig.BaseUri = 'https://api-v2-staging.jurni.me/v2'
    c.EnvConvig.SSLCaFile = '/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt'
  }else if c.EnvVariable == "production" {
    c.EnvConvig.BaseUri = 'https://api-v2.jurni.me/v2'
    c.EnvConvig.SSLCaFile = '/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt'
  }else {
     c.EnvConvig.BaseUri = 'http://api-v2.jurni-dev.me:4000/v2'
  }
}

func (c *Config) Test(){
  fmt.Println("== Test:")
}


func (r *RequestSetup) DoGet(){
  fmt.Println("== -- url = %s",r.Url)
  // logr "== -- params: #{params.to_s}")
  client := &http.Client{}
  req, _ := http.NewRequest("GET", r.Url, nil)
  // req.Header.Add("Host", "example.com")
  res, _ := client.Do(req)
}