package tentacle

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const retryTimes = 3

func (l *List) GetPages() (pages []string, err error) {
	pages = make([]string, 5)

	doc, err := goquery.NewDocument(l.Url)
	if err != nil {
		// need retry 3 times
		for i := 0; i < retryTimes; i++ {
			doc, err = goquery.NewDocument(l.Url)
			if err == nil {
				break
			} else if err != nil && i == retryTimes-1 {
				fmt.Println("GetPages NewDocument error. ", err)
				return pages, err
			}

			fmt.Println("GetPages retry ", i)
		}
	}

	doc.Find("#zg_page1 a").Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			pages[0] = path
		}
	})
	doc.Find("#zg_page2 a").Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			pages[1] = path
		}
	})
	doc.Find("#zg_page3 a").Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			pages[2] = path
		}
	})
	doc.Find("#zg_page4 a").Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			pages[3] = path
		}
	})
	doc.Find("#zg_page5 a").Each(func(i int, s *goquery.Selection) {
		path, exist := s.Attr("href")
		if exist {
			pages[4] = path
		}
	})

	return pages, nil
}

func (l *List) GetProducts(pages []string) (products []Product, err error) {
	products = make([]Product, 100)

	var count int = 0
	var rank int = 1

	for i := range pages {
		doc, err := goquery.NewDocument(pages[i])
		if err != nil {
			// need retry 3 times
			for i := 0; i < retryTimes; i++ {
				doc, err = goquery.NewDocument(pages[i])
				if err == nil {
					break
				} else if err != nil && i == retryTimes-1 {
					fmt.Println("GetProducts NewDocument error. ", err)
					return products, err
				}

				fmt.Println("GetProducts retry ", i)
			}
		}

		doc.Find(".zg_itemWrapper .zg_title a").Each(func(i int, s *goquery.Selection) {
			path, exist := s.Attr("href")
			if exist {
				// FIXME can I ignore error return in the below code snippet?
				path = strings.TrimSpace(path)
				path = strings.Trim(path, "\n")
				products[count], _ = getProductDetail(path, rank)
				count++
				rank++
			}
		})

	}

	return products, nil
}

func parseParameters(parameters string) string {
	var parsed string = ""

	if len(strings.TrimSpace(parameters)) == 0 {
		return parsed
	}

	regHTML := regexp.MustCompile("<[^>]*>")
	regScript := regexp.MustCompile("<script.*?>(.|\n)*</script>")
	regStyle := regexp.MustCompile("<style.*?>(.|\n)*</style>")

	// the following sequence is very important. step 1, remove script and its content; step 2, remove style and its content; step 3, finally remove all html directive, but not remove text.
	parameters = regScript.ReplaceAllString(parameters, "")
	parameters = regStyle.ReplaceAllString(parameters, "")
	parameters = regHTML.ReplaceAllString(parameters, "")

	parametersArray := strings.Split(parameters, "\n")

	var tempString string
	for i := range parametersArray {
		tempString = strings.TrimSpace(parametersArray[i])
		if len(tempString) != 0 {
			parsed += tempString + "|"
		}
	}

	parsed = strings.TrimRight(parsed, "|")

	return parsed
}

func getProductDetail(productUrl string, rank int) (product Product, err error) {
	fmt.Println(productUrl, rank)

	docDetail, err := goquery.NewDocument(productUrl)
	if err != nil {
		for i := 0; i < retryTimes; i++ {
			docDetail, err = goquery.NewDocument(productUrl)
			if err == nil {
				break
			} else if err != nil && i == retryTimes-1 {
				fmt.Println("getProductDetail NewDocument error. ", err)
				return product, err
			}

			fmt.Println("getProductDetail retry ", i)
		}
	}

	// get title
	product.Title = strings.TrimSpace(docDetail.Find("#productTitle").Text())
	// get manufacturer
	product.Manufacturer = strings.TrimSpace(docDetail.Find("#brand").Text())
	// get star
	product.Star, _ = docDetail.Find("#acrPopover").Attr("title")
	// get price
	var price string = strings.TrimSpace(docDetail.Find("#priceblock_ourprice").Text())
	if price == "" {
		price = strings.TrimSpace(docDetail.Find("#priceblock_saleprice").Text())
	}
	if price == "" {
		price = strings.TrimSpace(docDetail.Find("#priceblock_dealprice").Text())
	}
	if price == "" {
		price = strings.TrimSpace(docDetail.Find("#soldByThirdParty .price3P").Text())
	}
	product.Price = price
	// get rank
	product.Rank = rank
	// get url
	product.Url = productUrl
	// get image url
	imageUrl, _ := docDetail.Find("#imgTagWrapperId img").Attr("src")
	if strings.TrimSpace(imageUrl) != "" {
		imageUrl = strings.Trim(imageUrl, "\n")
	}
	product.ImageUrl = imageUrl
	// get technical information
	parameters, _ := docDetail.Find("#prodDetails").Html()
	if parameters == "" {
		parameters, _ = docDetail.Find("#detail-bullets").Html()
	}
	product.Parameters = parseParameters(parameters)

	return product, nil
}
