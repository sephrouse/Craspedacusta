package tentacle

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	var list = new(List)

	list.Url = "https://www.amazon.com/Best-Sellers-Camera-Photo-Digital-Picture-Frames/zgbs/photo/525460/ref=zg_bs_nav_p_1_p"

	var pages []string
	pages, err := list.GetPages()
	if err != nil {
		t.Errorf("get pages error, %s", err)
		return
	}

	var products []Product
	products, err = list.GetProducts(pages)
	if err != nil {
		t.Errorf("get products error, %s", err)
		return
	}

	for i := range products {
		fmt.Println(products[i].Rank, products[i].Title, products[i].Star, products[i].Price, products[i].Manufacturer, products[i].Parameters, products[i].Url, products[i].ImageUrl)
	}

	return
}
