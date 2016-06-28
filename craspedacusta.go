// author syney

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

// TODO use goroutine to retrieve each level I catagory concurrently.
// TODO store farther name of each catagory.
// TODO use 1 webdriver to get all URLs.
// TODO send the result of the excution when the program is done.

var baseUrl string

const retryTimes = 3

var baseCatagory []string = []string{"Camera & Photo", "Cell Phones & Accessories", "Computers & Accessories", "Electronics", "Office Products", "Sports & Outdoors", "Toys & Games"}
var allCatagory []Catagory

func isInTaskList(catagory string) bool {
	for i := range baseCatagory {
		if baseCatagory[i] == catagory {
			return true
		}
	}

	return false
}

func getLinkMap(catagory string) map[string]string {
	for i := range allCatagory {
		if allCatagory[i].Name == catagory {
			return allCatagory[i].Links
		}
	}

	return nil
}

func getLinksFromMenu(wd selenium.WebDriver, level int, fatherMenu string, url string) (err error) {
	err = nil

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
	subMenuElement = "#zg_browseRoot ul"

	for i := 0; i < level; i++ {
		subMenuElement += " ul"
	}
	subMenuElement += " li a"

	doc.Find(subMenuElement).Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			if level == 0 && isInTaskList(s.Text()) == false {
				fmt.Println("getLinksFromMenu: find each ", s.Text(), " is not what we want.")
				return
			}

			var menuName string
			if level == 0 {
				menuName = s.Text()
			} else if level == 1 {
				menuName = fatherMenu
			} else {
				fmt.Println("getLinksFromMenu: level ", level, " is more than expected.")
				return
			}

			links := getLinkMap(menuName)
			if nil == links {
				fmt.Println("getLinksFromMenu: getLinkMap gets nothing when menu name is ", menuName)
				return
			}

			links[s.Text()] = strconv.Itoa(level) + "|" + path

			// only gather catagories of both level I and II.
			if level == 0 {
				err := getLinksFromMenu(wd, level+1, menuName, path)
				if err != nil {
					fmt.Println("find each ", i, " getLinksFromMenu error. ", err)
					// FIXME: need retry?
				}
			}

		}

	})

	return nil
}

func showAllLinks() {
	fmt.Println("showAllLinks start:")

	for i := range allCatagory {
		fmt.Println("---", allCatagory[i].Name, " start:")

		for k, v := range allCatagory[i].Links {
			fmt.Println(k, v)
		}

		fmt.Println("---", allCatagory[i].Name, " end.")
	}

	fmt.Println("showAllLinks end.")
}

func initCatagory() {
	allCatagory = make([]Catagory, len(baseCatagory))

	for i := range allCatagory {
		allCatagory[i].Name = baseCatagory[i]
		allCatagory[i].Links = make(map[string]string)
	}
}

// TODO: use multi thread to call getLinksFromMenu.
// TODO: use goquery to analyze every item page and get useful information.
func main() {
	if len(os.Args) != 2 {
		// at present, only one parameter can be executed. the parameter means the root URL.
		fmt.Println("main: invalid number of parameters than expected. ", len(os.Args))
		return
	}

	baseUrl = os.Args[1]
	// won't check the real url to avoid the search from git/internet.
	if strings.Contains(baseUrl, "https://") == false {
		fmt.Println("main: seems that the input URL is not a correct URL. ", baseUrl)
		return
	}

	// there are more than 1 catagory, we need to make the maps as same amount as catagories.
	initCatagory()

	// create a new webdriver.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		//retry 3 times, otherwise print error and return it.
		for i := 0; i < retryTimes; i++ {
			wd, err = selenium.NewRemote(caps, "")
			if err == nil {
				break
			} else if err != nil && i == retryTimes-1 {
				fmt.Println("main: NewRemote error. ", err)
				return
			} else {
				fmt.Println("main: NewRemote retry ", i, " times.")
			}
		}
	}
	defer wd.Quit()

	// start to gather information from the root URL.
	err = getLinksFromMenu(wd, 0, "Any Department", baseUrl)
	if err != nil {
		fmt.Println("main: getLinksFromMenu error", err)
		return
	}

	//showAllLinks()

	// TODO: need to output information stored in catagory array into database.

	fmt.Println("end of program.")
}
