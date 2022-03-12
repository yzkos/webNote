package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed static/*
var static embed.FS

//go:embed tpl/*
var tpl embed.FS

var templates, _ = template.ParseFS(tpl, "tpl/*.html")

type Info struct {
	Code     string
	Password string
	Title    string
	Content  string
}

//获取笔记数据
func getInfo(code string) *Info {
	body, _ := ioutil.ReadFile("./data/" + code + ".json")
	I := Info{}
	err := json.Unmarshal(body, &I)
	if err != nil {
		return &Info{Code: code}
	}
	return &Info{Code: I.Code, Password: I.Password, Title: I.Title, Content: I.Content}
}

// Save 保存笔记
func Save(w http.ResponseWriter, r *http.Request) {
	oldPassword := r.FormValue("old_password")
	code := r.FormValue("code")
	data := getInfo(code)
	if data.Password != "" && data.Password != oldPassword {
		goUrl(w, r, "/"+code)
	}
	password := r.FormValue("password")
	title := r.FormValue("title")
	content := r.FormValue("content")
	p := &Info{Code: code, Password: password, Title: title, Content: content}
	res, _ := json.Marshal(p)
	err := ioutil.WriteFile("./data/"+p.Code+".json", res, 0600)
	if err != nil {
		return
	}
	goUrl(w, r, "/"+code)
}

//跳转页面
func goUrl(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

// Index 首页
func Index(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/"):]
	rand.Seed(time.Now().UnixNano())
	a := rand.Int()
	if code == "" {
		goUrl(w, r, "/edit/"+strconv.Itoa(a))
	}
	data := getInfo(code)
	if data.Content == "" {
		goUrl(w, r, "/edit/"+code)
	}
	loadTpl(w, data, "view")
}

// Edit 编辑笔记
func Edit(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/edit/"):]
	password := r.FormValue("password")
	data := getInfo(code)
	if data.Password != "" && password != data.Password {
		goUrl(w, r, "/password/"+code)
	}
	loadTpl(w, data, "edit")
}

// Password 验证密码
func Password(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/password/"):]
	data := getInfo(code)
	loadTpl(w, data, "password")
}

// 加载模板
func loadTpl(w http.ResponseWriter, data *Info, tpl string) {
	err := templates.ExecuteTemplate(w, tpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dataPath() {
	path := "data"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func main() {
	dataPath()
	http.Handle("/static/", http.StripPrefix("/", http.FileServer(http.FS(static))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/edit/", Edit)
	http.HandleFunc("/save/", Save)
	http.HandleFunc("/password/", Password)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
