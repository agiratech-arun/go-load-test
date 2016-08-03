package jurniapi_v2_client

import (
        "fmt"
        "net/http"
        "sync"
        "strconv"
        "time"
        "io/ioutil"
        "math/rand"
        "bytes"
        "encoding/json"
        "crypto/sha256"
        "encoding/base64"
        "os/exec"
        "net/url"
        _"strings"
       )

var base_wg sync.WaitGroup

var config Config

var app_config AppConfig

var commenter SignUpSession

var alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var key = "8YHsvw7fuylbLr5FevrFAsRC/v2sH5X8i9aWODH76908GxhIE/+jDj0cVJft+zTx2WkQmxiGM06KAnBtG1C7gg=="

var VideoUrl = ""

type Config struct {
  EnvVariable string
  Concurrency int
  PostConcurrency int
  EnvConvig EnvSetup
  RangeVaribales []string
  // AppConfig AppConfig
}


type SignUpSession struct {
  Status int
  UserId string `json:"user_id"`
  SessionId string `json:"session_id"`
  Error string `json:"error"`
  Post Post
}

type AppConfig struct {
  Status int
  ApiKey string `json:"api_key"`
  ApiSecret string   `json:"api_secret"`
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
  SessionId string
  // AppConfig AppConfig
}

type Post struct {
  Status int
  Error string `json:"error"`
  PostId string `json:"post_id"`
  IsPublic string `json:"is_public"`
  PostVideoUri string `json:"post_video_uri"`
  BackgroundOn bool `json:"background_on"`
  JobId string `json:"job_id"`
}

type TestRequest struct {
  Status int
  Error string `json:"error"`
}


// appliication starts here
func StepUp(env string, no int, post_no int,pwd string) {
  PrintSatement("Load Testing Steup")
  VideoUrl = pwd
  config.EnvVariable = env
  config.Concurrency = no
  config.PostConcurrency = post_no
  config.ConfigSetup()
  // config.TestGorotine()
  config.Register()
  config.TriggerLoadTest()
}

func (c *Config)TriggerLoadTest() {
  PrintSatement("Load Testing Started")
  fmt.Println(VideoUrl)
  commenter.CommenterSignUp()
  Concurrency()
}


func Concurrency() {
  for i:=0;i<config.Concurrency;i++ {
    var s SignUpSession
    base_wg.Add(1)
    go s.TrigegrConcurrency()
  }
  base_wg.Wait()
}

func (s *SignUpSession) TrigegrConcurrency() {
  defer base_wg.Done()
  // s.UserSignUp()
  s.Login("gkLz-20205","jurni123")
  s.PostTrigger()
}

func (c *Config) ConfigSetup() {
  // c.RangeVaribales =
  PrintSatement("Env Configuration")
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


func (r *RequestSetup) DoGet() (*http.Response,error){
  // defer base_wg.Done()
  PrintSatement("Get Request == -- url "+r.Url)
  client := &http.Client{}
  req, _ := http.NewRequest("GET", r.Url, nil)
  fmt.Println("%T",req)
  if r.SkipHeader != true{
    r.BuildHeader(req)
  }
  res, err := client.Do(req)
  return res,err
}


func (r *RequestSetup) DoPost() (*http.Response,error){
  client := &http.Client{}
  PrintSatement("Post Request == -- url "+r.Url)
  req, _ := http.NewRequest("POST", r.Url,  bytes.NewBufferString(r.Params))
  req.Header.Set("Content-Type", "application/json")
  if r.SkipHeader != true{
    r.BuildHeader(req)
  }
  res, err := client.Do(req)
  return res,err
}

// test the application to working
func (c *Config)TestGorotine(){
  PrintSatement("Api Test Goes")
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
  PrintSatement("Api Key Register")
  var r RequestSetup
  params := map[string]string{
    "app_secret": key}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.Url = config.EnvConvig.BaseUri + "/register"
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &app_config); err != nil {
          panic(err)
      }
    fmt.Println(app_config)
  }
}

// register a user based on appkey and appid
func (s *SignUpSession) UserSignUp(){
  PrintSatement("User SignUp")
  var r RequestSetup
  r.SignUpParams()
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &s); err != nil {
          panic(err)
      }
    fmt.Println(s)
  }
}

// Login user name
func (s *SignUpSession) Login(username string,pwd string) {
  PrintSatement(fmt.Sprintf("User %v Login",username))
  var r RequestSetup
  r.UserName = username
  r.Password = pwd
  r.SkipHeader = false
  r.Url = config.EnvConvig.BaseUri + "/users/login"
  params := map[string]string{
    "username": r.UserName,
    "password":   r.Password,
    "device_id": config.EnvConvig.DeviceId,
  }
  data,_ := json.Marshal(params)
  r.Params = string(data)
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &s); err != nil {
          panic(err)
      }
    fmt.Println(s)
  }
}

func (s *SignUpSession) PostTrigger() {
  var post_wg sync.WaitGroup
  for i:=0;i<config.PostConcurrency;i++ {
    post_wg.Add(1)
    go s.NewPost(&post_wg)
  }
  post_wg.Wait()
}

func (s *SignUpSession) CommenterSignUp(){
  PrintSatement("Commenter SignUp")
  var r RequestSetup
  r.SignUpParams()
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &commenter); err != nil {
          panic(err)
      }
    fmt.Println(commenter)
  }
}

//
func (s *SignUpSession) NewPost(post_wg *sync.WaitGroup ) {
  PrintSatement("Post Create")
  defer post_wg.Done()
  var r RequestSetup
  r.Url = fmt.Sprintf("%s/users/%v/posts/new",config.EnvConvig.BaseUri,s.UserId)
  params := map[string]int{"is_public": 1}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.SessionId = s.SessionId
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    var p Post
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &p); err != nil {
          panic(err)
      }
    fmt.Println(p)
    s.Post = p
    p.UploadVideo()
  }
}


func (s *SignUpSession) ShowPost(p *Post){
  PrintSatement("Show Post")
  var r RequestSetup
  r.Url = fmt.Sprintf("%v/users/%v/posts/%v",config.EnvConvig.BaseUri,s.UserId,p.PostId)
  params := map[string]string{}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    var test_data interface{}
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &test_data); err != nil {
          panic(err)
      }
    fmt.Println(test_data)
    // s.Fellow()
  }
}

func (s *SignUpSession) NewComment(p *Post) {
  // defer post_wg.Done()
  PrintSatement("Create Comment")
  var r RequestSetup
  r.Url = fmt.Sprintf("%s/users/%v/posts/%v/comments/new",config.EnvConvig.BaseUri,s.UserId,p.PostId)
  params := map[string]string{}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.SessionId = s.SessionId
  response,err := r.DoGet()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    var test_data interface{}
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &test_data); err != nil {
          panic(err)
      }
    fmt.Println(test_data)
    // s.Fellow()
  }
}

// func (s *SignUpSession) Fellow () {
//   var r RequestSetup
//   var commetor SignUpSession
//   commetor.SignUp()
//   params := map[string]string{
//     "topic_spec": "all"}
//   r.Url = fmt.Sprintf("%v/users/%v/follow/%v",config.EnvConvig.BaseUri,s.UserId,s.UserId)
//   data,_ := json.Marshal(params)
//   r.Params = string(data)
//   r.SessionId = s.SessionId
//   response,err := r.DoPost()
//   if err != nil {
//     fmt.Printf("Error %s \n",err)
//   }else {
//     body, err := ioutil.ReadAll(response.Body)
//     var s1 SignUpSession
//     if err = json.Unmarshal([]byte(body), &s1); err != nil {
//       panic(err)
//     }
//     fmt.Println(s1)
//   }
// }


func srand(size int) string {
  buf := make([]byte, size)
  rand.Seed(time.Now().UTC().UnixNano())
  for i := 0; i < size; i++ {
      buf[i] = alpha[rand.Intn(len(alpha))]
  }
  return string(buf)
}


func EncryptKey(digest_key string) string {
  h := sha256.Sum256([]byte(digest_key))
  return base64.StdEncoding.EncodeToString(h[:])
}

func (r *RequestSetup)BuildHeader(req *http.Request) {
  rand.Seed(time.Now().UTC().UnixNano())
  nonce := string(10000 + rand.Intn(89999))
  auth_key := EncryptKey(fmt.Sprintf("%v-%v-%v-%v", app_config.ApiKey, app_config.ApiSecret,nonce,r.Url))
  req.Header.Set("X-Api-Key", app_config.ApiKey)
  req.Header.Set("X-Api-Nonce", nonce)
  req.Header.Set("Authorization",auth_key)
  if r.SessionId != "" {
    req.Header.Set("X-Session-ID",r.SessionId)
  }
}

func (r *RequestSetup) SignUpParams() {
  rand.Seed(time.Now().UTC().UnixNano())
  rand_num := 10000 + rand.Intn(89999)
  r.UserName = fmt.Sprintf("%v-%v",srand(4),strconv.Itoa(rand_num))
  r.Password = "jurni123"
  r.Email = fmt.Sprintf("vivek+%v-%v@jurni.me",srand(8),strconv.Itoa(rand_num))
  r.SkipHeader = false
  r.Url = config.EnvConvig.BaseUri + "/signup"
  params := map[string]string{
    "username": r.UserName,
    "password":   r.Password,
    "email": r.Email,
    "device_id": config.EnvConvig.DeviceId,
  }
  data,_ := json.Marshal(params)
  r.Params = string(data)
}

func (p *Post) UploadVideo(){
  PrintSatement("Upload Post")
  uri,_ := url.Parse(p.PostVideoUri)
  comment := fmt.Sprintf("curl -X PUT -T %v '%v'",VideoUrl,uri.String())
  fmt.Println(exec.Command(comment).Run())
}

func PrintSatement(val string) {
  fmt.Printf("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
  fmt.Printf("    %v    ",val)
  fmt.Printf("\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
}