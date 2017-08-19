// https://github.com/niocs/SP500Addin/blob/master/getsp500.go
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"sort"
)

const URL string = "https://en.wikipedia.org/wiki/List_of_S%26P_500_companies"

type SP500Item struct {
	Ticker string
	Name   string
}

type SP500Items []SP500Item

func (slice SP500Items) Append(item SP500Item) SP500Items {
	return append(slice, item)
}

func (slice SP500Items) Len() int {
	return len(slice)
}

func (slice SP500Items) Less(i, j int) bool {
	return slice[i].Ticker < slice[j].Ticker
}

func (slice SP500Items) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func main() {

	doc, err := goquery.NewDocument(URL)
	if err != nil {
		fmt.Println("Error getting SP500 url, err =", err)
		return
	}

	items := SP500Items{}
	trs := doc.Find("table").First().Find("tbody").Find("tr")
	fmt.Println("trs len = ", trs.Length())
	trs.Each(func(idx int, s *goquery.Selection) {
		if idx == 0 {
			return
		}
		tds := s.Find("td")
		if (idx >= 0 && idx <= 20) || idx == 505 {
			fmt.Println("idx =", idx, "content = ", tds.Text())
		}
		td1 := tds.Get(0)
		td2 := tds.Get(1)
		ticker := td1.FirstChild.FirstChild.Data
		name := td2.FirstChild.FirstChild.Data
		items = items.Append(SP500Item{Ticker: ticker, Name: name})
	})

	sort.Sort(items)
	writeCSV(items)
	writeOUStringArray(items)
}

func writeCSV(items SP500Items) {
	fp, err := os.Create("sp500.csv")
	if err != nil {
		fmt.Println("Error writing to sp500.csv : err = ", err)
		return
	}

	for _, item := range items {
		fp.WriteString(fmt.Sprintf("%s,%s\n", item.Ticker, item.Name))
	}
	fp.Close()
}

func writeOUStringArray(items SP500Items) {
	fp, err := os.Create("data.hxx")
	if err != nil {
		fmt.Println("Error writing to data.hxx : err = ", err)
		return
	}

	fp.WriteString("#include <map>\n")
	fp.WriteString("#include <vector>\n")
	fp.WriteString("#include <rtl/ustring.hxx>\n\n")
	fp.WriteString("using rtl::OUString;\n\n\n")
	fp.WriteString("std::vector<OUString> aTickers = {\n")
	comma := ","
	for ii, item := range items {
		if ii+1 == len(items) {
			comma = ""
		}
		fp.WriteString(fmt.Sprintf(`    "%s"%s`+"\n", item.Ticker, comma))
	}
	fp.WriteString("};\n\n")
	fp.WriteString("std::map<OUString, OUString> aTickerToName = {\n")
	comma = ","
	for ii, item := range items {
		if ii+1 == len(items) {
			comma = ""
		}
		fp.WriteString(fmt.Sprintf(`    {"%s", "%s"}%s`+"\n", item.Ticker, item.Name, comma))
	}
	fp.WriteString("};\n")
	fp.Close()
}
