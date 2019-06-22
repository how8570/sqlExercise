package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", index) // 设置访问的路由
	http.HandleFunc("/login/action", loginAction)
	http.HandleFunc("/query", query)
	http.HandleFunc("/register", register)
	http.HandleFunc("/register/action", registerAction)
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
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }

	pageByte, err := ioutil.ReadFile("./pages/index.html")
	checkErr(err)

	page := string(pageByte[:])
	fmt.Fprintf(w, page) // 这个写入到 w 的是输出到客户端的

}

func loginAction(w http.ResponseWriter, r *http.Request) {
	var page string

	r.ParseForm()

	username := r.Form.Get("username")
	pass := r.Form.Get("pass")

	db, err := sql.Open("sqlite3", "./food.db")
	checkErr(err)
	defer db.Close()
	stmt, err := db.Prepare("SELECT * FROM users WHERE username = ? AND pass = ?;")
	defer stmt.Close()
	checkErr(err)
	q, err := stmt.Query(username, pass)
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

	fmt.Fprintf(w, page)
}

func query(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pageByte, err := ioutil.ReadFile("./pages/query.html")
	checkErr(err)

	page := string(pageByte[:])

	q := r.Form.Get("q")
	if q == "stores" {
		page += `
		<body>
			<table width="50%" border="1">
				<tr>
				<td>ID</td>
				<td>店名</td>
				<td>開始營業時間</td>
				<td>結束營業時間</td>
				<td>地址</td>
				<td>評論</td>
				<td>操作選項</td>
				</tr>
			
		`

		db, err := sql.Open("sqlite3", "./food.db")
		checkErr(err)
		defer db.Close()

		rows, err := db.Query("SELECT * FROM stores")
		checkErr(err)
		var id int
		var name string
		var open_begin string
		var open_end string
		var location string
		var comment string

		for rows.Next() {
			err = rows.Scan(&id, &name, &open_begin, &open_end, &location, &comment)
			checkErr(err)
			page += `<tr>
			<td>` + strconv.Itoa(id) + `</td>
			<td>` + name + `</td>
			<td>` + open_begin + `</td>
			<td>` + open_end + `</td>
			<td>` + location + `</td>
			<td>` + comment + `</td>
			<td>
			<input type="button" value="修改"  onclick="location.href='http://google.com'">
			<input type="button" value="刪除"  onclick="location.href='http://google.com'">
			</td>
			</tr>
			`
		}
		rows.Close()

		page += `
			</table>
		</body>`
	}

	fmt.Fprintf(w, page)
}

func register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pageByte, err := ioutil.ReadFile("./pages/register/register.html")
	checkErr(err)

	page := string(pageByte[:])

	fmt.Fprintf(w, page)
}

func registerAction(w http.ResponseWriter, r *http.Request) {
	var page string

	r.ParseForm()

	db, err := sql.Open("sqlite3", "./food.db")
	checkErr(err)
	defer db.Close()

	// insert
	stmt, err := db.Prepare("INSERT INTO users(username, pass, email) values(?,?,?)")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(r.Form.Get("username"), r.Form.Get("pass"), r.Form.Get("email"))
	if err != nil {
		pageByte, err := ioutil.ReadFile("./pages/register/register_fail.html")
		checkErr(err)
		page = string(pageByte[:])
	} else {
		rows, err := res.RowsAffected()
		checkErr(err)
		if rows != 1 {
			log.Fatalf("expected to affect 1 row, affected %d", rows)
		}

		pageByte, err := ioutil.ReadFile("./pages/register/register_success.html")
		checkErr(err)
		page = string(pageByte[:])
	}

	fmt.Fprintf(w, page)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
