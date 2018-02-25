package main

import ("net/http"
        "log"
        )

//Authentication/page protection middleware
func authenticateFunc(loggedInFunc http.HandlerFunc,
                      loggedOutFunc http.HandlerFunc) http.HandlerFunc {
    return func(res http.ResponseWriter, req *http.Request) {
        if loggedIn(req) {
            loggedInFunc(res, req)
        } else {
            loggedOutFunc(res, req)
        }
    }
}

func authenticateHandler(loggedInHandler http.Handler,
                         loggedOutHandler http.Handler) http.Handler{
    return http.HandlerFunc(authenticateFunc(loggedInHandler.ServeHTTP,
                                             loggedOutHandler.ServeHTTP))
}

//Login checks
func loggedIn(req *http.Request) bool {
    log.Print("starting authentication")
    loggedin, err := GetCookie(req, "loggedin")
    if err != nil {
        log.Fatal("Failed to open cookie")
        return false
    }
    if loggedin == "false" || loggedin == "" { //return error/login page
        log.Print("authentication failed: not logged in")
        log.Print(loggedin)
        return false
    } else { 
        log.Print("authentication approved")
        return true
    }
}

//Login cookie management
func SetCookie(w http.ResponseWriter, r *http.Request, value map[string]string) error{
    if encoded, err := store.Encode(config.LoginCookieName, value); err == nil {
        cookie := &http.Cookie{
            Name: config.LoginCookieName,
            Value: encoded,
            Path: "/",
        }
        http.SetCookie(w, cookie)
        return nil
    } else {
        log.Print("Cookie write fail")
        log.Print(err)
        return err
    }
}

func GetCookie(r *http.Request, key string) (string, error) {
    if cookie, err := r.Cookie(config.LoginCookieName); err == nil {
        value := make(map[string]string)
        if err = store.Decode(config.LoginCookieName, cookie.Value, &value); err == nil {
            return value[key], err
        } else {
            log.Print("failed to decode cookie")
            return "", err
        }
    } else {
        log.Print("failed to retrieve cookie")
        return "", err
    }
}

