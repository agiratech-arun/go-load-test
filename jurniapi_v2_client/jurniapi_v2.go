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
        _"os/exec"
        "net/url"
        "github.com/codeskyblue/go-sh"
        _"strings"
       )

var base_wg sync.WaitGroup
var config Config
var app_config AppConfig
var commenter UserSession
var login_users []User
var alpha = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
var key = "8YHsvw7fuylbLr5FevrFAsRC/v2sH5X8i9aWODH76908GxhIE/+jDj0cVJft+zTx2WkQmxiGM06KAnBtG1C7gg=="
var VideoUrl = ""
var username = "jurni_test0"
var password = "jurni123"
var post_concurrency = 1
var comment_concurrency = 1

type Config struct {
  EnvVariable string
  Concurrency int
  EnvConvig EnvSetup
  RangeVaribales []string
}

type User struct {
  Status int
  Error string `json:"error"`
  UserId string `json:"user_id"`
  UserName string `json:"username"`
}

type UserSession struct {
  Status int
  UserId string `json:"user_id"`
  SessionId string `json:"session_id"`
  Error string `json:"error"`
  Post Post
  Comment Comment
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

type Comment struct {
  Status int
  Error string `json:"error"`
  CommentId string `json:"comment_id"`
  IsPublic string `json:"is_public"`
  CommentVideoUri string `json:"comment_video_uri"`
  // Created string `json:"created_at"`
  BackgroundOn bool `json:"background_on"`
  JobId string `json:"job_id"`
  Commenter map[string]string `json:"commenter"`
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

type UserSearch struct {
  Status int
  Error string `json:"error"`
  UserCount int `json:"num_users"`
  Users []User `json:"users"`
}

type PostSearch struct {
  Status int `json:"status"`
  Error string `json:"error"`
  PostCount int `json:"num_posts"`
  Posts []Post `json:"posts"`
}

// appliication starts here
func StepUp(env string, concurrency int, video_path string,method_name string) {
  PrintSatement("Load Testing Setup")
  if concurrency == 0 {
    concurrency = 1
  }
  VideoUrl = video_path
  config.EnvVariable = env
  config.Concurrency = concurrency
  config.ConfigSetup()
  config.Register()
  switch method_name {
    case "scenario_1":
      ScenarioOne(concurrency,strconv.Itoa(concurrency))
      break;
    case "scenario_2":
      ScenarioTwo(concurrency,strconv.Itoa(concurrency))
      break;
    default:
      fmt.Println("\n=============================================================================")
      fmt.Println("Usange go run main.go --e #Environment --c #concurrency --method #method_name")
      fmt.Println("=============================================================================\n")
  }
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


func ScenarioOne(n int, limit string) {
  PrintSatement("Scenario One")
  var scenario_1_wg sync.WaitGroup
  var s UserSession
  s.Login(username,password)
  s.SearchUser(limit)
  posts := s.SearchPublishedPost()
  for _,user := range login_users {
    scenario_1_wg.Add(1)
    rand.Seed(time.Now().UTC().UnixNano())
    post := posts[rand.Intn(len(posts))]
    new_user := user
    go new_user.ScenarioOneFlow(&scenario_1_wg, post)
  }
  scenario_1_wg.Wait()

}

func ScenarioTwo(n int, limit string) {
  PrintSatement("Scenario Two")
  var feed_metric_wg sync.WaitGroup
  var s UserSession
  s.Login(username,password)
  s.SearchUser(limit)
  fmt.Println("Login Users",len(login_users))
  for _,user := range login_users {
    feed_metric_wg.Add(1)
    new_user := user
    go new_user.ScenarioTwoFlow(&feed_metric_wg)
  }
  feed_metric_wg.Wait()
}

func (user *User)ScenarioOneFlow(scenario_1_wg *sync.WaitGroup, post Post) {
  defer scenario_1_wg.Done()
  var post_wg sync.WaitGroup
  var session UserSession
  session.Post = post
  session.Login(user.UserName, password)
  if len(session.Error) == 0 {
    post_wg.Add(2)
    go session.PostTrigger(&post_wg)
    go session.CommentTrigger(&post_wg)
    post_wg.Wait()
  }else {
    fmt.Println("Error",session.Error)
  }
}

func (user *User)ScenarioTwoFlow(scenario_2_wg *sync.WaitGroup) {
  defer scenario_2_wg.Done()
  var post_wg sync.WaitGroup
  var session UserSession
  session.Login(user.UserName, password)
  if len(session.Error) == 0 {
    post_wg.Add(1)
    go session.SearchFeed(&post_wg)
    post_wg.Wait()
  }else{
    fmt.Println("Error",session.Error)
  }
}

// Login user name
func (s *UserSession) Login(username string,pwd string) {
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

func (s *UserSession) PostTrigger(go_wg *sync.WaitGroup) {
  defer go_wg.Done()
  var post_wg sync.WaitGroup
  for i:=0;i<post_concurrency;i++ {
    post_wg.Add(1)
    go s.NewPost(&post_wg)
  }
  post_wg.Wait()
}

func (s *UserSession) CommentTrigger(go_wg *sync.WaitGroup) {
  defer go_wg.Done()
  var comment_wg sync.WaitGroup
  for i:=0;i<comment_concurrency;i++ {
    comment_wg.Add(1)
    go s.NewComment(&comment_wg)
  }
  comment_wg.Wait()
}

//
func (s *UserSession) NewPost(post_wg *sync.WaitGroup ) {
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
    // s.Post = p
    UploadVideo(VideoUrl,p.PostVideoUri)
  }
}

func (s *UserSession) NewComment(comment_wg *sync.WaitGroup) {
  defer comment_wg.Done()
  PrintSatement("Create Comment")
  var r RequestSetup
  r.Url = fmt.Sprintf("%s/users/%v/posts/%v/comments/new",config.EnvConvig.BaseUri,s.UserId,s.Post.PostId)
  params := map[string]string{}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.SessionId = s.SessionId
  response,err := r.DoPost()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    var c Comment
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &c); err != nil {
      panic(err)
    }
    fmt.Println(c)
    s.Comment = c
    UploadVideo(VideoUrl,c.CommentVideoUri)
  }
}

// Search published posts
func (s *UserSession) SearchPublishedPost()([]Post){
  PrintSatement("Search Post")
  var p PostSearch
  var r RequestSetup
  r.Url = fmt.Sprintf("%v/posts/published_post_search",config.EnvConvig.BaseUri)
  fmt.Println(r.Url)
  params := map[string]string{}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.SessionId = s.SessionId
  response,err := r.DoGet()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &p); err != nil {
      panic(err)
    }
    fmt.Println(p)
  }
  return p.Posts
}

// Search published posts
func (s *UserSession) SearchFeed(go_wg *sync.WaitGroup) {
  defer go_wg.Done()
  PrintSatement("Search Post")
  var p PostSearch
  var r RequestSetup
  r.Url = fmt.Sprintf("%v/users/%v/public",config.EnvConvig.BaseUri,s.UserId)
  fmt.Println(r.Url)
  params := map[string]string{}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.SessionId = s.SessionId
  response,err := r.DoGet()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    if err = json.Unmarshal([]byte(body), &p); err != nil {
      panic(err)
    }
    fmt.Println(p)
  }
}

// Search users
func (s *UserSession) SearchUser(limit string){
  PrintSatement("User Search")
  var r RequestSetup
  r.Url = fmt.Sprintf("%v/users/public/user_search?limit=%v",config.EnvConvig.BaseUri,limit)
  params := map[string]string{}
  data,_ := json.Marshal(params)
  r.Params = string(data)
  r.SessionId = s.SessionId
  response,err := r.DoGet()
  if err != nil {
    fmt.Printf("Error %s \n",err)
  }else {
    body, err := ioutil.ReadAll(response.Body)
    var u UserSearch
    if err = json.Unmarshal([]byte(body), &u); err != nil {
      panic(err)
    }
    login_users =  u.Users
    fmt.Println(u)
  }
}


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
  nonce := strconv.Itoa(10000 + rand.Intn(89999))
  auth_key := EncryptKey(fmt.Sprintf("%v-%v-%v-%v", app_config.ApiKey, app_config.ApiSecret,nonce,r.Url))
  req.Header.Set("X-Api-Key", app_config.ApiKey)
  req.Header.Set("X-Api-Nonce", nonce)
  req.Header.Set("Authorization",auth_key)
  fmt.Println("~~~~~~~~~~~~~~Headers~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
  fmt.Println("Authorization",auth_key)
  fmt.Println("X-Api-Nonce",nonce)
  fmt.Println("X-Api-Key",app_config.ApiKey)
  fmt.Println("X-Session-ID",r.SessionId)
  fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
  if r.SessionId != "" {
    req.Header.Set("X-Session-ID",r.SessionId)
  }
}


func UploadVideo(file_path string, url_string string){
  PrintSatement("Uploading")
  uri,_ := url.Parse(url_string)
  fmt.Println(sh.Command("curl", "-X", "PUT", "-T",  file_path, uri.String()).Run())
  PrintSatement("Upload Finished")

}

func PrintSatement(val string) {
  fmt.Printf("\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
  fmt.Printf("    %v    ",val)
  fmt.Printf("\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
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

func (c *Config) ConfigSetup() {
  // c.RangeVaribales =
  PrintSatement("Env Configuration")
  c.EnvConvig.DeviceId = srand(64)
  if c.EnvVariable == "staging" {
    c.EnvConvig.BaseUri = "http://api-v2-staging.jurni.me/v2"
    // c.EnvConvig.SSLCaFile = "/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt"
  }else if c.EnvVariable == "production" {
    c.EnvConvig.BaseUri = "http://api-v2-staging.jurni.me/v2"
    // c.EnvConvig.SSLCaFile = "/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt"
  }else {
     c.EnvConvig.BaseUri = "http://api-v2-staging.jurni.me/v2"
     // c.EnvConvig.SSLCaFile = "/home/ubuntu/jurni_devops/conf/ssl/new/gd_bundle-g2-g1.crt"
  }
}