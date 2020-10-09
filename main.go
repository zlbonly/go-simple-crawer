package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

/*
var (
	BaseUrl = "https://movie.douban.com/top250"
)

func Add(movies []parse.DoubanMovie) {
	for index, movie := range movies {
		if err := model.DB.Create(&movie).Error; err != nil {
			log.Printf("db.Create index:%s,err:%v", index, err)
		}
	}

}

func Start() {

	var movies []parse.DoubanMovie
	pages := parse.GetPages(BaseUrl)

	fmt.Println(len(pages))

	for _, page := range pages {
		doc, err := goquery.NewDocument(strings.Join([]string{BaseUrl, page.Url}, ""))
		if err != nil {
			log.Println(err)
		}

		movies = append(movies, parse.ParseMovies(doc)...)
	}


	fmt.Println(len(movies))
	Add(movies)
}*/

/*func main() {

	Start()
	defer model.DB.Close()
}
*/

func HttpGet(url string) (result string, err error) {
	fmt.Println("url->", url)

	client := &http.Client{}
	req, err1 := http.NewRequest("GET", url, nil)

	if err1 != nil {
		err = err1
		return
	}

	req.Header.Add("User-Agent", "test")

	resp, err2 := client.Do(req)
	if err2 != nil {
		err = err2
		return
	}

	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err3 := resp.Body.Read(buf)
		if n == 0 {
			break
		}

		if err3 != nil && err3 != io.EOF {
			err = err3
			return
		}
		result += string(buf[:n])
	}

	return
}

func Write2File(idx int, filmName, filmScore, filmRate [][]string) {
	f, err := os.Create("第 " + strconv.Itoa(idx) + "页.txt")
	if err != nil {
		fmt.Println("os.Create err", err)
		return
	}
	defer f.Close()
	n := len(filmName)
	f.WriteString("电影名字\t\t\t\t电影分数\t\t\t\t电影评论数\n")
	for i := 0; i < n; i++ {
		f.WriteString(filmName[i][1] + "\t\t\t\t" + filmScore[i][1] + "\t\t\t\t" + filmRate[i][1] + "\n")
	}
}

func SpiderPage(idx int, page chan int) {
	url := "https://movie.douban.com/top250?start=" + strconv.Itoa((idx-1)*25) + "&filter="
	result, err := HttpGet(url)
	if err != nil {
		fmt.Println("HttpGet err:", err)
		return
	}

	nameExp := regexp.MustCompile(`<img width="100" alt="(.*?)"`)
	filmName := nameExp.FindAllStringSubmatch(result, -1)

	scoreExp := regexp.MustCompile(`<span class="rating_num" property="v:average">(.*?)</span>`)
	filmScore := scoreExp.FindAllStringSubmatch(result, -1)

	rateExp := regexp.MustCompile(`<span>(.*?)人评价</span>`)
	filmRate := rateExp.FindAllStringSubmatch(result, -1)

	Write2File(idx, filmName, filmScore, filmRate)

	fmt.Println(filmName)
	// 与主线程同步，写入当前page数
	page <- idx
}

func toWork(start, end int) {
	fmt.Printf("正在爬取%d到%d页...\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d页读取完毕\n", <-page)
	}
}

func main() {
	var start, end int
	fmt.Print("请输入开始爬取的页 (>=1):")
	fmt.Scan(&start)
	fmt.Print("请输入结束爬取的页 (>=start):")
	fmt.Scan(&end)
	toWork(start, end)
}
