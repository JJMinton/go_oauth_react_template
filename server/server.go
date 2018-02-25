package main


import ("net/http"
        "io"
        "io/ioutil"
        "log"
        "os"

        "encoding/json"

        //"github.com/gorilla/sessions"
        "github.com/gorilla/securecookie"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        )

//Config
//var store *sessions.CookieStore
var store *securecookie.SecureCookie
type Config struct {
    RootURL         string `json:"rootURL"`
    Port            string `json:"port"`
    CookieHashKey   []byte `json:"cookieHashKey"`
    CookieBlockKey  []byte `json:"cookieBlockKey"`
    LoginCookieName      string `json:"loginCookieName"`
}
var config Config

// Construction config with credentials using init
type Credentials struct {
    Cid     string `json:"cid"`
    Csecret string `json:"csecret"`
    Callback string `json:"callbackURL"`
}
var googleCred Credentials
var googleConf *oauth2.Config

func init() {
    //Generic config
    file, err := ioutil.ReadFile("./server/config.json")
    if err != nil {
        log.Printf("File error: %v\n", err)
        os.Exit(1)
    }
    json.Unmarshal(file, &config)
    log.Printf("Starting server at %s on port %s\n", config.RootURL, config.Port)
    //store = securecookie.New([]byte(config.CookieHashKey), []byte(config.CookieBlockKey))
    store = securecookie.New(config.CookieHashKey, config.CookieBlockKey)

    //Google oauth config
    file, err = ioutil.ReadFile("./server/google_creds.json")
    if err != nil {
        log.Printf("File error: %v\n", err)
        os.Exit(1)
    }
    json.Unmarshal(file, &googleCred)

    googleConf = &oauth2.Config{
        ClientID:       googleCred.Cid,
        ClientSecret:   googleCred.Csecret,
        RedirectURL:    config.RootURL+googleCred.Callback,
        Scopes: []string{
            "email",
        },
        Endpoint: google.Endpoint,
    }
}

//Router
func main() {
    http.HandleFunc("/root", rootHandler)
    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./react/dist/public"))))
    http.Handle("/private/", authenticateHandler(
        http.StripPrefix("/private/", http.FileServer(http.Dir("./react/dist/private"))),
        http.HandlerFunc(failedLoginHandler)))
    http.HandleFunc("/auth/google/login", googleLoginHandler)
    http.HandleFunc(googleCred.Callback, googleCallbackHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/protectedpage", authenticateFunc(protectedPageHandler, failedLoginHandler))
    http.ListenAndServe(config.Port, nil)
}


//Handlers
func rootHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <p><a href="/auth/google/login">LOGIN</a></p>
    <p><a href="/logout">LOGOUT</a></p>
    <p><a href="/protectedpage">Test login</a></p>
  </body>
</html>`)
}

func failedLoginHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <p>Failed login, try again.</p>
    <a href="/">Login Page</a>
  </body>
</html>`)
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
    logout(res, req)
    http.Redirect(res, req, "/", http.StatusSeeOther)
}

func protectedPageHandler(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <p>This page shouldn't be accessible without logging in</p>
    <a href="/logout">Logout</a>
  </body>
</html>`)

}

//Google login
var state string
func googleLoginHandler(res http.ResponseWriter, req *http.Request) {
    state = alphaNum(32)
    http.Redirect(res, req, googleConf.AuthCodeURL(state), http.StatusSeeOther)
}

func googleCallbackHandler(w http.ResponseWriter, req *http.Request) {
    //Check state
    if state != req.FormValue("state") {
        log.Fatal("Response state doesn't match set state: auth forgery attack?!")
    }

    //Retrieve access token
    authcode := req.FormValue("code")
    tok, err := googleConf.Exchange(oauth2.NoContext, authcode)
    if err != nil {
        log.Fatal("err is ", err)
    }
    log.Print(string(tok.AccessToken))

    //Retrieve google info
    response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
      log.Fatal("failed to read google respones body")
    }
    log.Print(string(contents))
    //TODO: do a check and login here.
    log.Print("logging in from google callback handler")
    login(w, req)
    io.WriteString(w, `
<!DOCTYPE html>
<html>
  <head></head>
  <body>` +
   `<p>You are now logged in.</p>` +
   `<a href="/protectedpage">Test on this protected page</a>` +
`  </body>
</html>`)
}


//Login/logout
func login(res http.ResponseWriter, req *http.Request) {
    log.Print("saving logged in cookies")
    err := SetCookie(res, req, map[string]string{"loggedin": "true",})
    if err != nil {
        log.Fatal("Failed to save cookie")
    }
}

func logout(res http.ResponseWriter, req *http.Request) {
    log.Print("saving logged out cookies")
    err := SetCookie(res, req, map[string]string{"loggedin": "false",})
    if err != nil {
        log.Fatal("Failed to save cookie")
    }
}

