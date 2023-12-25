package main

import (
	"os"
	"fmt"
	"unicode"
	"strings"
	"strconv"
	"time"
	"golang.org/x/text/width"

    flags "github.com/jessevdk/go-flags"
)

func getDateEra(date string) (y, m, d, era int) {

	var (
		isDigit     = false
		str         = ""
		dateStr     []string
	)

	switch {
	case strings.ContainsAny(date, "ＭｍMm明"):
		era = 1
	case strings.ContainsAny(date, "ＴｔTt大"):
		era = 2
	case strings.ContainsAny(date, "ＳｓSs昭"):
		era = 3
	case strings.ContainsAny(date, "ＨｈHh平"):
		era = 4
	case strings.ContainsAny(date, "ＲｒRr令"):
		era = 5
	default:
		era = 0
	}

    // dateの前処理
	date = width.Fold.String(date) + "."

    // dateから数字の部分だけ取り出す
	for _, char := range date {
		if unicode.IsDigit(char) {
			str = str + string(char)
			isDigit = true
		} else {
			if isDigit {
				dateStr = append(dateStr, str)
				str = ""
				isDigit = false
			} else if string(char) == "元" {
				str = "1"
				isDigit = true
			}
		}
	}

	if len(dateStr) != 3 {
		fmt.Println("日付の取り出しに失敗しました")
		os.Exit(1)
	}

	y, _ = strconv.Atoi(dateStr[0])
	m, _ = strconv.Atoi(dateStr[1])
	d, _ = strconv.Atoi(dateStr[2])

	return y, m, d, era
}

func toBC(y, m, d, era int, Kanji bool) string {

	var	magicNumber int

	switch era {
	case 1:
		magicNumber = 1867
	case 2:
		magicNumber = 1911
	case 3:
		magicNumber = 1924
	case 4:
		magicNumber = 1988
	case 5:
		magicNumber = 2018
	}

	if Kanji {
		return fmt.Sprintf("%d年%d月%d日", y + magicNumber, m, d)
	} else {
		return fmt.Sprintf("%d/%d/%d", y + magicNumber, m, d)
	}
}

func toJP(y, m, d int, Kanji bool) string {

	var (
		nameK		string
		name		string
		magicNumber	int
	)

	date := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

	switch {
	case date.Before(time.Date(1912, time.Month(7), 30, 0, 0, 0, 0, time.Local)):
		nameK = "明治"
		name = "M"
		magicNumber = 1867
	case date.Before(time.Date(1926, time.Month(12), 25, 0, 0, 0, 0, time.Local)):
		nameK = "大正"
		name = "T"
		magicNumber = 1911
	case date.Before(time.Date(1989, time.Month(1), 8, 0, 0, 0, 0, time.Local)):
		nameK = "昭和"
		name = "S"
		magicNumber = 1925
	case date.Before(time.Date(2019, time.Month(5), 1, 0, 0, 0, 0, time.Local)):
		nameK = "平成"
		name = "H"
		magicNumber = 1988
	default:
		nameK = "令和"
		name = "R"
		magicNumber = 2018
	}

	if Kanji {
		if y == 1 {
			return fmt.Sprintf("%s元年%d月%d日", nameK, m, d)
		} else {
			return fmt.Sprintf("%s%d年%d月%d日", nameK, y, m, d)
		}
	} else {
		return fmt.Sprintf("%s%d/%d/%d", name, y, m, d)
	}
}

func main() {
 
	type option struct {
		Kanji bool `short:"k" long:"kanji" description:"Returns the date in the format X年X月X日"`
    }

    var opt option

    var parser = flags.NewParser(&opt, flags.Default)

    parser.Name = "Era"
    parser.Usage = "[OPTIONS] date"

    args, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	if len(args) != 1 {
		fmt.Println("There is only one command line argument")
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	y, m, d, era := getDateEra(args[0])

	if era == 0 {
		fmt.Println(toJP(y, m, d, opt.Kanji))
	} else {
		fmt.Println(toBC(y, m, d, era, opt.Kanji))
	}
}



