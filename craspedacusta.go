// author syney

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

var baseUrl string

const retryTimes = 3

//var baseCatagory []string = {"Camera & Photo", "Cell Phones & Accessories", "Computers & Accessories", "Electronics", "Office Products", "Sports & Outdoors", "Toys & Games"}
var baseCatagory []string = []string{"Camera & Photo"} //example for test.

// allCataListLinks map[menu string] link string contain all top 100 items of each catagory.
var allCataListLinks map[string]string

func getLinksFromMenu(level int, fatherMenu string, url string) (err error) {
	err = nil

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		//retry 3 times, otherwise print error and return it.
		for i := 0; i < retryTimes; i++ {
			wd, err = selenium.NewRemote(caps, "")
			if err == nil {
				break
			} else if err != nil && i == retryTimes-1 {
				fmt.Println("getLinksFromMenu: ", level, fatherMenu, url, err)
				return err
			} else {
				fmt.Println("getLinksFromMenu: retry ", i, " times.")
			}
		}
	}
	defer wd.Quit()

	wd.Get(url)

	realHtml, err := wd.PageSource()
	if err != nil {
		fmt.Println("getLinksFromMenu: wd.PageSource error. ", err)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(realHtml))
	if err != nil {
		fmt.Println("getLinksFromMenu: goquery.NewDocumentFromReader error. ", err)
		return err
	}

	var subMenuElement string
	var subMenuCount int = 0
	var subMenuExisted bool = false

	subMenuElement = "#zg_browseRoot ul"
	for i := 0; i < level; i++ {
		subMenuElement += " ul"
	}
	subMenuElement += " li a"

	doc.Find(subMenuElement).Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			if level == 0 && s.Text() != baseCatagory[0] {
				fmt.Println("find each : ", s.Text(), " is not ", baseCatagory[0])
				return
			}
			//menuName := fatherMenu + "|" + s.Text()
			menuName := s.Text()
			err := getLinksFromMenu(level+1, menuName, path)
			if err != nil {
				fmt.Println("find each ", i, " getLinksFromMenu error. ", err)
				// fixme need retry?
			}
		}

		subMenuCount = i
		subMenuExisted = true
	})

	if subMenuCount == 0 && subMenuExisted == false {
		allCataListLinks[fatherMenu] = url
		fmt.Println("find a leaf url. level ", level, fatherMenu, url)
	}

	return nil
}

func showAllLinks() {
	fmt.Println("showAllLinks start:")

	for k, v := range allCataListLinks {
		fmt.Println(k, v)
	}

	fmt.Println("showAllLinks end.")
}

// todo: use multi thread to call getLinksFromMenu
// todo: use goquery to analyze every item page and get useful information.
func main() {
	baseUrl = os.Args[1]

	// initial links map
	allCataListLinks = make(map[string]string)

	err := getLinksFromMenu(0, "Any Department", baseUrl)
	if err != nil {
		fmt.Println("main getLinksFromMenu error", err)
		return
	}

	fmt.Println("end of program.")
}
