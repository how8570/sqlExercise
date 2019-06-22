package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", index) // 设置访问的路由
	http.HandleFunc("/login/action", loginAction)
	http.HandleFunc("/query", query)
	http.HandleFunc("/register", register)
	http.HandleFunc("/register/action", registerAction)
	http.HandleFunc("/register/result", registerResult)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	// fmt.Println(r.Form) // 这些信息是输出到服务器端的打印信息
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	pageByte, err := ioutil.ReadFile("./pages/index.html")
	checkErr(err)

	page := string(pageByte[:])
	fmt.Fprintf(w, page) // 这个写入到 w 的是输出到客户端的

}

func loginAction(w http.ResponseWriter, r *http.Request) {
	var page string

	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	account := r.Form.Get("account")
	pass := r.Form.Get("pass")

	if account == "" || pass == "" { // 回到 index
		page = `<!DOCTYPE html>
		<script>
		window.location.replace("/"); 
		</script>
		`
	} else {
		db, err := sql.Open("sqlite3", "./food.db")
		checkErr(err)
		defer db.Close()
		sql := "SELECT * FROM users WHERE account = '" + account + "' AND " + "pass = '" + pass + "';"
		fmt.Println(sql)
		q, err := db.Query(sql)
		checkErr(err)
		if q.Next() {
			page = `
			<!DOCTYPE html>
			<script>
			window.location.replace("/query"); 
			</script>
			`
		} else { // 回到 index
			page = `<!DOCTYPE html>
			<body>
			<p>沒這個 user 或是 密碼錯了喔</p>
			<a id="確定" href="/">
			<input type="submit" value="確定">
			</a>
			</body>
			</html>
			`
		}
		q.Close()
	}

	fmt.Fprintf(w, page)
}

func query(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	page := `<!DOCTYPE html>
	<html>
	<body>

	<p>:D</p>

	</body>
	</html>
	`

	fmt.Fprintf(w, page)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	page := `
	`

	fmt.Fprintf(w, page)
}

func registerAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	for k, v := range r.Form {

		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	page := `
		註冊成功
	`

	fmt.Fprintf(w, page)
}

func registerResult(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	for k, v := range r.Form {

		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	page := `
		註冊成功
	`

	fmt.Fprintf(w, page)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
