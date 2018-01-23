package main

import (
  "log"
  "net/http"
  "os"
  "fmt"
)

func main() {
  // get root directory, port and password header from command line arguments
  // https://gobyexample.com/command-line-arguments
  // ./main /var/data :3000 xyz
  root := os.Args[1]
  // TCP port
  // eg.: ':3000'
  port := os.Args[2]
  // authentication header
  pass := os.Args[3]

  if len(os.Args) >= 5 {
    // ./main /var/data :3000 xyz /var/log/go.log
    file := os.Args[4]
    // log to file
    f, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
    if err != nil {
      log.Fatalf("Error opening log file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)
  }

  log.Println("------")
  log.Println("Starting symlinks API")

  // root must end with bar
  lastRootChar := root[len(root) - 1:]
  if lastRootChar != "/" {
    root += "/"
  }
  log.Println("Server root")
  log.Println(root)

  http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
    // check authentication
    auth := r.Header.Get("X-Authentication")
    if auth != pass {
      unauthorized(w)
      return
    }

    _oldname, ok := r.URL.Query()["oldname"]
    if !ok || len(_oldname) < 1 {
      // no oldname query param
      clientError(w)
      return
    }

    _newname, ok := r.URL.Query()["newname"]
    if !ok || len(_newname) < 1 {
      // no newname query param
      clientError(w)
      return
    }

    oldname := fmt.Sprintf("%s%s", root, _oldname[0])
    newname := fmt.Sprintf("%s%s", root, _newname[0])
    // check if file already exists
    if _, err := os.Lstat(newname); err == nil {
      // remove file before symlink creation
      os.Remove(newname)
    }
    err := os.Symlink(oldname, newname)
    if err != nil {
      // permission ?
      log.Fatalf("Cannot creat symlink: %v", err)
      serverError(w)
      return
    }

    success(w)
  })

  http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
    // check authentication
    auth := r.Header.Get("X-Authentication")
    if auth != pass {
      unauthorized(w)
      return
    }

    _newname, ok := r.URL.Query()["newname"]
    if !ok || len(_newname) < 1 {
      // no newname query param
      clientError(w)
      return
    }

    newname := fmt.Sprintf("%s%s", root, _newname[0])
    // check if file exists, then remove
    if _, err := os.Lstat(newname); err == nil {
      os.Remove(newname)
    }

    success(w)
  })

  log.Println("Listening...")
  log.Println(port)
  log.Fatal(http.ListenAndServe(port, nil))
}

func unauthorized(w http.ResponseWriter) {
  // 401 response
  w.WriteHeader(http.StatusUnauthorized)
  w.Write([]byte("Unauthorized!\n"))
}

func success(w http.ResponseWriter) {
  // 200 response
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("OK!\n"))
}

func clientError(w http.ResponseWriter) {
  // 400 response
  w.WriteHeader(http.StatusBadRequest)
  w.Write([]byte("Bad Request!\n"))
}

func serverError(w http.ResponseWriter) {
  // 500 response
  w.WriteHeader(http.StatusInternalServerError)
  w.Write([]byte("Internal Server Error!\n"))
}
