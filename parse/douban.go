package parse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type DoubanMovie struct {
	Title    string
	Subtitle string
	Other    string
	Desc     string
	Year     string
	Area     string
	Tag      string
	Star     string
	Comment  string
	Quote    string
}
type Page struct {
	Page int
	Url  string
}

// 获取分页
func GetPages(url string) []Page {
	fmt.Println(url)
	doc, err := goquery.NewDocument(url)

	fmt.Println(doc.Html())
	if err != nil {
		log.Fatal(err)
	}
	return ParsePages(doc)
}

// 分析分页
func ParsePages(doc *goquery.Document) (pages []Page) {

	doc.Find("#content > h1").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(i)
		fmt.Println(selection.Text())
	})

	pages = append(pages, Page{Page: 1, Url: ""})

	doc.Find("#content").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(i)
		fmt.Println(selection.Text())
	})

	doc.Find("ol.grid_view li div.hd a").
		Each(func(i int, s *goquery.Selection) {
			// year
			fmt.Printf("%v", s)
		})
	doc.Find("#content > div > div.article > div.paginator > a").Each(func(i int, s *goquery.Selection) {

		page, _ := strconv.Atoi(s.Text())
		fmt.Println(page)
		url, _ := s.Attr("href")
		pages = append(pages, Page{
			Page: page,
			Url:  url,
		})
	})
	return pages
}

// 分析电影数据
func ParseMovies(doc *goquery.Document) (movies []DoubanMovie) {
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".hd a span").Eq(0).Text()
		subtitle := s.Find(".hd a span").Eq(1).Text()
		subtitle = strings.TrimLeft(subtitle, "  / ")
		other := s.Find(".hd a span").Eq(2).Text()
		other = strings.TrimLeft(other, "  / ")
		desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
		DescInfo := strings.Split(desc, "\n")
		desc = DescInfo[0]
		movieDesc := strings.Split(DescInfo[1], "/")
		year := strings.TrimSpace(movieDesc[0])
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])
		star := s.Find(".bd .star .rating_num").Text()
		comment := strings.TrimSpace(s.Find(".bd .star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")
		quote := s.Find(".quote .inq").Text()
		movie := DoubanMovie{
			Title:    title,
			Subtitle: subtitle,
			Other:    other,
			Desc:     desc,
			Year:     year,
			Area:     area,
			Tag:      tag,
			Star:     star,
			Comment:  comment,
			Quote:    quote,
		}
		log.Printf("i: %d, movie: %v", i, movie)
		movies = append(movies, movie)
	})
	return movies
}
