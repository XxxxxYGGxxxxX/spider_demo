// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"regexp"
// 	"strconv"

// 	_ "github.com/go-sql-driver/mysql"

// 	"github.com/PuerkitoBio/goquery"
// )

// // 数据库信息定义
// const (
// 	username = "root"
// 	password = "123456"
// 	host     = "127.0.0.1"
// 	port     = "3306"
// 	dbname   = "douban_movie"
// )

// // 原生SQL
// var db *sql.DB

// // 定义一个结构体
// type MovieData struct {
// 	Title    string `json:"title"`
// 	Director string `json:"Director"`
// 	Picture  string `json:"Picture"`
// 	Actor    string `json:"Actor"`
// 	Year     string `json:"Year"`
// 	Score    string `json:"Score"`
// 	Quote    string `json:"Quote"`
// }

// func main() {
// 	err := InitDB()
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 	} else {
// 		fmt.Println("链接成功!")
// 	}
// 	for i := 1; i <= 10; i++ {
// 		fmt.Printf("正在爬取第 %d 页信息\n", i)
// 		//strconv.Itoa()将int转换为string
// 		Spider(strconv.Itoa(i * 25))
// 	}
// 	// //并发
// 	// ch := make(chan bool)
// 	// for i := 0; i < 10; i++ {
// 	// 	go Spider(strconv.Itoa(i * 25),ch)
// 	// }
// 	// for i := 0; i < 10; i++ {
// 	// 	<-ch
// 	// }
// }

// func Spider(page string, ch chan bool) {
// 	//1.发送请求
// 	//构造客户端
// 	client := http.Client{}
// 	// 构造GET请求
// 	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)
// 	if err != nil {
// 		fmt.Println("req err", err)
// 	}
// 	//添加请求头信息，模拟浏览器访问，防止服务器检测爬虫访问
// 	req.Header.Add("Connection", "keep-alive")
// 	req.Header.Add("Cache-Control", "max-age=0")
// 	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
// 	req.Header.Add("Sec-Ch-Ua-Platform", "Windows")
// 	req.Header.Add("Sec-Fetch-Dest", "document")
// 	req.Header.Add("Sec-Fetch-Mode", "navigate")
// 	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
// 	req.Header.Add("Upgrade-Insecure-Requests", "1")
// 	req.Header.Add("Sec-Fetch-User", "?1")
// 	req.Header.Add("Sec-Fetch-Site", "none")
// 	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
// 	resp, err := client.Do(req) // 发送请求
// 	if err != nil {
// 		fmt.Println("请求失败", err)
// 	}
// 	// bodyText, _ := ioutil.ReadAll(resp.Body) // 使用缓冲区读取网页内容
// 	defer resp.Body.Close()

// 	//2.解析网页
// 	//github.com/PuerkitoBio/goquery 提供了 .NewDocumentFromReader 方法进行网页的解析
// 	docDetial, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		fmt.Println("解析失败", err)
// 	}

// 	//3.获取节点信息
// 	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.hd > a > span:nth-child(1)
// 	//拿到上一步解析出来的 doc 之后，可以进行 css选择器语法，进行结点的选择(复制)
// 	// #content > div > div.article > ol > li
// 	// #content > div > div.article > ol > li:nth-child(1) > div > div.pic > a > img
// 	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p:nth-child(1)
// 	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > div > span.rating_num
// 	//#content > div > div.article > ol > li:nth-child(1) > div > div.info > div.bd > p.quote > span
// 	docDetial.Find("#content > div > div.article > ol > li"). //列表
// 									Each(func(i int, s *goquery.Selection) { //在列表中继续查找
// 			var data MovieData //将数据存储到结构体中
// 			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
// 			img := s.Find("div > div.pic > a > img") //拿到的是标签，内容在属性里
// 			imgTmp, ok := img.Attr("src")
// 			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
// 			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
// 			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()
// 			if ok {
// 				director, actor, year := InfoSpite(info)
// 				data.Title = title
// 				data.Director = director
// 				data.Picture = imgTmp
// 				data.Actor = actor
// 				data.Year = year
// 				data.Score = score
// 				data.Quote = quote

// 				//4.保存信息
// 				if InsertData(data) {
// 					// fmt.Println("插入成功")
// 				} else {
// 					fmt.Println("插入失败")
// 					return
// 				}
// 				// fmt.Println("data:", data)
// 			}
// 		})
// 	// fmt.Println("插入成功")
// 	return
// }

// // 正则表达式格式模板
// func InfoSpite(info string) (director, actor, year string) {
// 	directorRe, _ := regexp.Compile(`导演:(.*)`)
// 	director = string(directorRe.Find([]byte(info)))
// 	actorRe, _ := regexp.Compile(`主演:(.*)`)
// 	actor = string(actorRe.Find([]byte(info)))
// 	yearRe, _ := regexp.Compile(`(\d+)`)
// 	year = string(yearRe.Find([]byte(info)))
// 	return
// }

// // 数据库的初始化
// func InitDB() (err error) {
// 	dsn := "root:123456@tcp(127.0.0.1:3306)/douban_movie?charset=utf8mb4&parseTime=True"
// 	db, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		return err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // 插入数据
// func InsertData(movieData MovieData) bool {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		fmt.Println("begin err", err)
// 		return false
// 	}
// 	// stmt, err := tx.Prepare("INSERT INTO movie_data (`Title`,`Director`,`Picture`,`Actor`,`Year`,`Score`,`Quote`) VALUES (?,?,?,?,?,?,?)")
// 	// if err != nil {
// 	// 	fmt.Println("preare fail err", err)
// 	// 	return false
// 	// }
// 	// _, err = stmt.Exec(movieData.Title, movieData.Director, movieData.Picture, movieData.Year, movieData.Score, movieData.Quote)
// 	// if err != nil {
// 	// 	fmt.Println("exec fail", err)
// 	// 	return false
// 	// }
// 	stmt, err := tx.Prepare("INSERT INTO test_movie (`Title`,`Director`) VALUES (?,?)")
// 	if err != nil {
// 		fmt.Println("preare fail err", err)
// 		return false
// 	}
// 	_, err = stmt.Exec(movieData.Title, movieData.Director)
// 	if err != nil {
// 		fmt.Println("exec fail", err)
// 		return false
// 	}
// 	_ = tx.Commit()
// 	return true

// }
