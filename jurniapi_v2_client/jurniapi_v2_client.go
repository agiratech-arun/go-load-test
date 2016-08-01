package jurniapi_v2_client

import (
        "fmt"
        "net/http"
        "sync"
        _"time"
        "io/ioutil"
        "math/rand"
        "bytes"
        "encoding/json"
        "crypto/sha256"
        "encoding/base64"
       )

var base_wg sync.WaitGroup

var config Config


var alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type Config struct {
  EnvVariable string
  Concurrency int
  EnvConvig EnvSetup
  RangeVaribales []string
  AppConfig AppConfig
}


type AppConfig struct {
  Status int
  ApiKey string `json:"api_key"`
  ApiSecret string  `json:"api_secret"`
  SessionId *string
}

type EnvSetup struct {
  BaseUri string
  SSLCaFile string
  DeviceId string
}

type RequestSetup struct {
  Url string
  Params  string
  Headers map[string]string
  Password string
  UserName string
  Email string
  SkipHeader bool
  AppConfig AppConfig
}

type TestRequest struct {
  Args map[string]string
  Data string
  Files map[string]string
  Header map[string]string
  Json string
  Origin string
  Url string
}



func (c *Config) ConfigSetup() {
  // c.RangeVaribales =
  c.EnvConvig.DeviceId = srand(64)
  if c.EnvVariable == "staging" {
    c.EnvConvig.BaseUri = "https://api-v2-staging.jurni.me/v2"
    c.EnvConvig.SSLCaFile = "/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt"
  }else if c.EnvVariable == "production" {
    c.EnvConvig.BaseUri = "https://api-v2.jurni.me/v2"
    c.EnvConvig.SSLCaFile = "/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt"
  }else {
     c.EnvConvig.BaseUri = "http://vpc-api-v2.jurni.me/v2"
  }
}

func (c *Config) Test(){
  fmt.Println("== Test:")
}


func (r *RequestSetup) DoGet(){
  defer base_wg.Done()
  fmt.Printf("== -- url = %s\n",r.Url)
  // logr "== -- params: #{params.to_s}")
  client := &http.Client{}
  req, _ := http.NewRequest("GET", r.Url, nil)
  fmt.Println("%T",req)
  // req.Header.Add("Host", "example.com")
  res, err := client.Do(req)
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else{
   fmt.Printf("\nresponse html %s\n",res)
  }
}


func (r *RequestSetup) DoPost() (*http.Response,error){
  client := &http.Client{}
  req, _ := http.NewRequest("POST", r.Url,  bytes.NewBufferString(r.Params))
  if r.SkipHeader != true{

  }
  res, err := client.Do(req)

  return res,err
}

// test the application to working
func (c *Config)TestGorotine(){
  fmt.Println("-- test starting")
  var r RequestSetup
  r.Url = config.EnvConvig.BaseUri + "/test"
  for i:=0;i<c.Concurrency;i++ {
    base_wg.Add(1)
    go r.DoGet()
  }
  base_wg.Wait()
}
// registeration for app to get appkey and appid
func (c *Config)Register() {
  var r RequestSetup
  params := map[string]string{
    "app_secret": "8YHsvw7fuylbLr5FevrFAsRC/v2sH5X8i9aWODH76908GxhIE/+jDj0cVJft+zTx2WkQmxiGM06KAnBtG1C7gg=="}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  // var params = []byte(`{"app_secret":"8YHsvw7fuylbLr5FevrFAsRC/v2sH5X8i9aWODH76908GxhIE/+jDj0cVJft+zTx2WkQmxiGM06KAnBtG1C7gg=="}`)
  // r.Params = params
  r.Url = config.EnvConvig.BaseUri + "/register"
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    var app AppConfig
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &app); err != nil {
          panic(err)
      }
    c.AppConfig = app
    fmt.Println(app)
  }
}
// register a user based on appkey and appid
func (c *Config)SignUp(){
  var r RequestSetup
  rand_num := 10000 + rand.Intn(89999)
  r.UserName = fmt.Sprintf("%s-%s",srand(4),rand_num)
  r.Password = "jurni123"
  r.Email = fmt.Sprintf("vivek+%s-%s@jurni.me",srand(8),rand_num)
  r.SkipHeader = false
  r.Url = config.EnvConvig.BaseUri + "/signup"
  r.AppConfig = c.AppConfig
  // params := map[string]string{
  //   "username": r.UserName,
  //   "password":   r.Password,
  //   "email": r.Email,
  //   "device_id": c.EnvConvig.DeviceId,
  // }
  // data,_ := json.Marshal(params)
  // r.Params = string(data)
  response,err := r.DoPost()
  fmt.Println(response,err)


}
// appliication starts here
func StepUp(env string, no int) {
  config.EnvVariable = env
  config.Concurrency = no
  config.ConfigSetup()
  // config.TestGorotine()
  config.Register()
  fmt.Println(config.AppConfig)
  // config.SignUp()
}

// func startRequest() {
//   defer base_wg.Done()
//   fmt.Println("Start request",time.Now())
//   response, err :=http.Get("http://golang.org")
//   fmt.Println(response,err)
//   fmt.Println("end request")
// }


func srand(size int) string {
  buf := make([]byte, size)
  for i := 0; i < size; i++ {
      buf[i] = alpha[rand.Intn(len(alpha))]
  }
  return string(buf)
}


func EncryptPassword2(digest_key string) string {
  h := sha256.Sum256([]byte(digest_key))
  return base64.StdEncoding.EncodeToString(h[:])
}

// func (r *RequestSetup)BuildHeader(req *http.Request ) {
//   nonce := string(10000 + rand.Intn(89999))
//   auth_key := EncryptPassword2(fmt.Sprintf("%v-%v-%v-%v", r.ApiKey, r.ApiSecret,nonce,r.Url))
//   req.Header.Set("X-Api-Key", r.ApiKey)
//   req.Header.Set("X-Api-Nonce", nonce)
//   req.Header.Set("Authorization",auth_key)
//   req.Header.Set("X-Session-ID",r.SessionId)

// }
