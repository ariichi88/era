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

	// dateから年号を取り出す
	if strings.Contains(date, "MｍMm明") {
		era = 1
	} else if strings.Contains(date, "ＨｔTt大") {
		era = 2
	} else if strings.Contains(date, "ＳｓSs昭") {
		era = 3
	} else if strings.Contains(date, "ＨｈHh平") {
		era = 4
	} else if strings.Contains(date, "ＲｒRr令") {
		era = 5
	} else {
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

func toJP(y, m, d int,) string {

	var (
		nameK		string
		name		string
		magicNumber	int
	)

	date := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	
	meiji  := time.Date(1912, time.Month(7), 30, 0, 0, 0, 0, time.Local)
	taisho := time.Date(1926, time.Month(12), 25, 0, 0, 0, 0, time.Local)
	showa  := time.Date(1989, time.Month(1), 7, 0, 0, 0, 0, time.Local)
	heisei := time.Date(2019, time.Month(4), 30, 0, 0, 0, 0, time.Local)

	if date.Before(meiji) {
		nameK = "明治"
		name = "M"
		magicNumber = 1867
	} else if date.Before(taisho) {
		nameK = "大正"
		name = "T"
		magicNumber = 1911
	} else if date.Before(showa) {
		nameK = "昭和"
		name = "S"
		magicNumber = 1924
	} else if date.Before(heisei) {
		nameK = "平成"
		name = "H"
		magicNumber = 1988
	} else {
		nameK = "令和"
		name = "R"
		magicNumber = 2018
	}





}

func main() {
 
    // オプションの作成
    type option struct {
        Kanji bool `short:"k" long:"kanji" description:"日付をＸ年Ｘ月Ｘ日の形式で返します"`
    }

    var opt option

    // パース
    var parser = flags.NewParser(&opt, flags.Default)

    // コマンド名と使用法の指定
    parser.Name = "Era"
    parser.Usage = "[OPTIONS] date"

	// argsを取り出す（失敗したらhelpを出して終わる）
    args, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
}



