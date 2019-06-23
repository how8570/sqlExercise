package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/", index) // 设置访问的路由
	r.HandleFunc("/login/action", loginAction)
	r.HandleFunc("/query", query)
	r.HandleFunc("/store/{id:[0-9]+}", modifyStore)
	r.HandleFunc("/register", register)
	r.HandleFunc("/register/action", registerAction)
	err := http.ListenAndServe(":80", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func modifyStore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	page := `<!DOCTYPE html>
	<head></head>
	`
	vars := mux.Vars(r)
	id := vars["id"]
	//fmt.Fprintf(w, id)

	if r.Form.Get("a") == "upd" {

		db, err := sql.Open("sqlite3", "./food.db")
		checkErr(err)
		defer db.Close()
		stmt, err := db.Prepare("UPDATE stores SET name = ?, open_begin = ?, open_end = ?, location = ?, comment = ?  WHERE id = ?;")
		checkErr(err)
		defer stmt.Close()
		res, err := stmt.Exec(r.Form.Get("name"), r.Form.Get("open_begin"), r.Form.Get("open_end"), r.Form.Get("location"), r.Form.Get("comment"), id)
		for err != nil {
			time.Sleep(30 * time.Millisecond)
			res, err = stmt.Exec(r.Form.Get("name"), r.Form.Get("open_begin"), r.Form.Get("open_end"), r.Form.Get("location"), r.Form.Get("comment"), id)
		}

		lastID, err := res.LastInsertId()
		checkErr(err)
		rowCnt, err := res.RowsAffected()
		checkErr(err)
		log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)

		page += "<p>更新成功!</p>"
	} else if r.Form.Get("a") == "del" {
		page += `<p>真的要刪除麻?</p>
		<form name="confirm" action="./` + id + `?a=delCONF"  method="POST">
		<input type="submit" value="確定">
		</form>
		`
	} else if r.Form.Get("a") == "delCONF" {

		page += `<script>
		alert("刪除成功")
		window.location.replace("/query");
		</script>
		`
	}

	page += `
	<body>
		<table width="80%" border="1">
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

	rows, err := db.Query("SELECT * FROM stores WHERE id = " + id)
	checkErr(err)
	var name string
	var open_begin string
	var open_end string
	var location string
	var comment string

	if rows.Next() {
		err = rows.Scan(&id, &name, &open_begin, &open_end, &location, &comment)
		checkErr(err)
		page += `
		<form name="myForm" method="POST" action="./` + id + `" >
			<tr>
			<td>` + id + `</td>
			<td><input type="text" name="name" value="` + name + `" required></td>
			<td><input type="text" name="open_begin" value="` + open_begin + `" required></td>
			<td><input type="text" name="open_end" value="` + open_end + `" required></td>
			<td><input type="text" name="location" value="` + location + `" required></td>
			<td><input type="message" name="comment" value="` + comment + `" ></td>
			<td>
			<input type="hidden" name="a" value="upd">
			<input type="submit" value="修改">
			<input type="submit" value="刪除" formaction = "./` + id + `" onsubmit="return del()">
			<input type="reset" value="回復">
			</td>
			</tr>
		</form>
		<script>
		function del() {
			document.forms["myForm"]["a"].value = 'del';
			alert(document.forms["myForm"]["a"].value);
			return false;
		  }
		</script>
		`
	}
	rows.Close()

	page += `
		</table>
	</body>`

	fmt.Fprintf(w, page)
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
			<input type="button" value="修改"  onclick="location.href='./store/` + strconv.Itoa(id) + `'">
			</td>
			</tr>
			`
		}
		rows.Close()

		page += `
			</table>
		</body>`
	} else if q == "random" {
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
		var id int
		var name string
		var open_begin string
		var open_end string
		var location string
		var comment string

		var storesCount int
		rows, err := db.Query("SELECT COUNT(*) FROM stores")
		checkErr(err)
		if rows.Next() {
			err = rows.Scan(&storesCount)
			checkErr(err)
		}
		rows.Close()

		rows, err = db.Query("SELECT * FROM stores WHERE id = " + strconv.Itoa(rand.Intn(storesCount)+1))
		if rows.Next() {
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
			<input type="button" value="修改"  onclick="location.href='./store/` + strconv.Itoa(id) + `'">
			</td>
			</tr>
			`
		}
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
		log.Fatal(err)
	}
}
