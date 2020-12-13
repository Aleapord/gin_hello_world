package main

import "github.com/gin-gonic/gin"

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}
var mutex = sync.Mutex{}
var num = 0
var has_pa = make(map[string]bool)

func get_body(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		my_panic(err)
	}
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	req.Header.Set("Referer", "https://www.nvshens.org/")
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		my_panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		my_panic(err)
	}
	return body
}
func save_image(path string, body []byte, num *int) bool {
	err := ioutil.WriteFile("./image/"+path+strconv.Itoa(*num)+".jpg", body, 0755)
	mutex.Lock()
	*num++
	mutex.Unlock()
	if err != nil {
		return false
	}
	return true
}
func gen_urls(url string, size int) []string {
	links := make([]string, 0)
	links = append(links, url)
	for i := 1; i < 10; i++ {
		u := url[:len(url)-4] + strconv.Itoa(0) + strconv.Itoa(i) + ".jpg"
		links = append(links, u)
	}
	for i := 10; i < size; i++ {
		u := url[:len(url)-4] + strconv.Itoa(i) + ".jpg"
		links = append(links, u)
	}
	return links

}
func paqu(url string, size int, path string, num *int) {
	image_links := gen_urls(url, size)
	var t_num = 16
	url_chan := make(chan string, size)
	//进程退出标志管道
	exit_chan := make(chan bool, t_num)
	//为本美女创建一个文件夹
	_ = os.Mkdir("./image/"+path, 0755)
	wg.Add(1)
	go func() {
		for _, i := range image_links {
			url_chan <- i
		}
		close(url_chan)
		wg.Done()
	}()
	wg.Add(t_num)
	for i := 0; i < t_num; i++ {
		go func() {
			for i := range url_chan {
				save_image(path+"/", get_body(i), num)
			}
			exit_chan <- true
			wg.Done()
		}()
	}
	wg.Add(1)
	go func() {
		for i := 0; i < t_num; i++ {
			<-exit_chan
		}
		close(exit_chan)
		wg.Done()
	}()
	wg.Wait()
}

func get_urlr(url string) (string, int, string) {
	var info string
	var ret string
	var size int
	var title string
	doc, err := goquery.NewDocument(url)
	if err != nil {
		my_panic(err)
	}
	title = doc.Find("#htilte").Text()
	doc.Find("#dinfo").Each(func(i int, selection *goquery.Selection) {
		info = selection.Text()
	})
	doc.Find("#hgallery").Each(func(i int, selection *goquery.Selection) {
		t := selection.Find("img")
		a, _ := t.Attr("src")
		ret = a
		fmt.Println(a)
	})
	r, _ := regexp.Compile("[0-9]+张")
	s := r.FindString(info)
	size, _ = strconv.Atoi(s[0:2])
	fmt.Println(title)
	fmt.Println(size)
	return ret, size, title
}
func girl_pachong(girl string, num *int) {
	url := "https://www.nvshens.org"
	links := make([]string, 0)
	doc, err := goquery.NewDocument(girl)
	if err != nil {
		my_panic(err)
	}
	t := doc.Find(".post_entry")
	t.Find(".igalleryli_link").Each(func(i int, selection *goquery.Selection) {
		a, _ := selection.Attr("href")
		links = append(links, url+a)
	})
	for _, i := range links {
		a, b, c := get_urlr(i)
		paqu(a, b, c, num)
	}
}
func dfs_pachong(url string) {
	url_l := "https://www.nvshens.org"
	doc, err := goquery.NewDocument(url)
	links := make([]string, 0)
	if err != nil {
		my_panic(err)

	}
	doc.Find(".suggestWrapper").Each(func(i int, selection *goquery.Selection) {
		selection.Find("li.galleryli").Each(func(i int, selection *goquery.Selection) {
			selection.Find("a.galleryli_link").Each(func(i int, selection *goquery.Selection) {
				a, _ := selection.Attr("href")
				links = append(links, url_l+a)
			})
		})
	})
	for _, link := range links {
		if has_pa[link] == false {
			has_pa[link] = true
			mjson, _ := json.Marshal(has_pa)
			_ = ioutil.WriteFile("./has_pa.json", mjson, 0755)
			a, b, c := get_urlr(link)
			paqu(a, b, c, nil)
			dfs_pachong(link)
		}
	}

}
func my_panic(err error) {
	mjson, _ := json.Marshal(has_pa)
	_ = ioutil.WriteFile("./has_pa.json", mjson, 0755)
	panic(err)
}

//func main() {
//	var m_type int
//
//START:
//	fmt.Println("输入爬虫模式：1自动，2输入专辑网址,3结束")
//
//	_, _ = fmt.Scan(&m_type)
//	if m_type == 1 {
//		var url string
//		fmt.Println("输入开始地址：")
//		_, _ = fmt.Scan(&url)
//		auto_pachong(url)
//	} else if m_type == 2 {
//		var url string
//		fmt.Println("请输入专辑网址：")
//		_, _ = fmt.Scan(&url)
//		girl_pachong(url)
//	} else if m_type != 3 {
//		goto START
//	}
//
//}
func auto_pachong(url string) {
	f, _ := ioutil.ReadFile("./has_pa.json")
	_ = json.Unmarshal(f, &has_pa)
	dfs_pachong(url)
}

func paqu2(url string, size int, path string, num *int) {
	image_links := gen_urls(url, size)
	var t_num = 8
	url_chan := make(chan string, size)
	//进程退出标志管道
	exit_chan := make(chan bool, t_num)
	//为本美女创建一个文件夹
	wg.Add(1)
	go func() {
		for _, i := range image_links {
			url_chan <- i
		}
		close(url_chan)
		wg.Done()
	}()
	wg.Add(t_num)
	for i := 0; i < t_num; i++ {
		go func() {
			for i := range url_chan {
				save_image("", get_body(i), num)
			}
			exit_chan <- true
			wg.Done()
		}()
	}
	wg.Add(1)
	go func() {
		for i := 0; i < t_num; i++ {
			<-exit_chan
		}
		close(exit_chan)
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	eng := gin.Default()
	eng.StaticFS("/static", http.Dir("./image"))
	eng.LoadHTMLGlob("./html/*")
	eng.GET("/", func(context *gin.Context) {
		_, _ = context.Writer.WriteString("hello world")
	})

	eng.GET("/pa", func(context *gin.Context) {
		exec.Command("rm -rf ./image/*")
		var num = 0
		var s []string
		var my_map = map[string]interface{}{}
		var url = context.DefaultQuery("url", "")
		a, b, c := get_urlr(url)
		paqu2(a, b, c, &num)
		for i := 0; i < num; i++ {
			s = append(s, "/static/"+strconv.Itoa(i)+".jpg")
		}
		my_map["title"] = c
		my_map["srcs"] = s
		context.HTML(http.StatusOK, "img.html", my_map)

	})
	_ = eng.Run()

}
