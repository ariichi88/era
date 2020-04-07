package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

    flags "github.com/jessevdk/go-flags"
)

// 文字列から日付及び元号を取り出す
func getDateAndEra(date string) (y, m, d, era int) {

	var (
		isDigit     = false
		str         = ""
		dateStr     []string
	)

	// どの年号か判別してera（元号)を求める
	if strings.ContainsAny(date, "明MmＭｍ") {
		era = 1 // 明治
	} else if strings.ContainsAny(date, "大TtＴｔ") {
		era = 2 // 大正
	} else if strings.ContainsAny(date, "昭SsＳｓ") {
		era = 3 // 昭和
	} else if strings.ContainsAny(date, "平HhＨｈ") {
		era = 4 // 平成
	} else if strings.ContainsAny(date, "令RrＲｒ") {
        era = 5 // 令和
    } else {
        era = 0 // 西暦
    }

    // dateの前処理(最後が文字でないと次の処理で最後の数字が取り出せないため)
	date = date + "."

    // 与えられた文字列から数字の部分だけ取り出す
	for _, char := range date {
		if unicode.IsDigit(char) {
			str = str + string(char)
			isDigit = true
		} else {
			if isDigit {
				dateStr = append(dateStr, str)
				str = ""
				isDigit = false
			} else if string(char) == "元" { // 元年を1に変換
				str = "1"
				isDigit = true
			}
		}
	}

	// 日付を正しく取得できてなかった場合(dateStrの長さが３でない)
	if len(dateStr) != 3 {
       fmt.Println("日付の取得に失敗しました")
       os.Exit(1)
    }

	// 日付を取得
	y, _ = strconv.Atoi(dateStr[0])
	m, _ = strconv.Atoi(dateStr[1])
	d, _ = strconv.Atoi(dateStr[2])

    return y, m, d, era
}


// 和暦を西暦に変換
func convertToChristianCal(y, m , d, era int, kanji bool) string {

    var magicNumber int

    switch era {
    case 1: // 明治
        magicNumber = 1868
    case 2: // 大正
        magicNumber = 1911
    case 3: // 昭和
        magicNumber = 1925
    case 4: // 平成
        magicNumber = 1988
    case 5: // 令和
        magicNumber = 2018
    }

    year := strconv.Itoa(y + magicNumber)
    month := strconv.Itoa(m)
    date := strconv.Itoa(d)

    if kanji {
        return year + "年" + month + "月" + date + "日"
    } else {
        return year + "/" + month + "/" + date
    }
}

// 西暦を和暦に変換
func convertToJapaneseCal(y, m, d int, kanji bool) string {

    var (
        magicNumber       int
        name, nameK       string
    )

    // 日付の作成
    t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

    // 元号を判定する
    if t.Before(time.Date(1912, 7, 30, 0, 0, 0, 0, time.Local)) {
        magicNumber = 1867
        name = "M"
        nameK = "明治"
    } else if t.Before(time.Date(1926, 12, 25, 0, 0, 0, 0, time.Local)) {
        magicNumber = 1911
        name = "T"
        nameK = "大正"
    } else if t.Before(time.Date(1989, 1, 8, 0, 0, 0, 0, time.Local)) {
        magicNumber = 1925
        name = "S"
        nameK = "昭和"
    } else if t.Before(time.Date(2019, 5, 1, 0, 0, 0, 0, time.Local)) {
        magicNumber = 1988
        name = "H"
        nameK = "平成"
    } else {
        magicNumber = 2018
        name = "R"
        nameK = "令和"
    }

    year := strconv.Itoa(y - magicNumber)
    if kanji {
        if year == "1" {
            year = nameK + "元"
        } else {
            year = nameK + year
        }
    } else {
        year = name + year
    }

    month := strconv.Itoa(m)
    date := strconv.Itoa(d)

    if kanji {
        return year + "年" + month + "月" + date + "日"
    } else {
        return year + "." + month + "." + date
    }
}

func main() {

    type option struct {
        Kanji bool `short:"k" long:"kanji" description:"日付をＸ年Ｘ月Ｘ日の形式で返します"`
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

	if len(args) == 0 || len(args) > 1 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

    y, m, d, era := getDateAndEra(args[0])

    if era == 0 {
        fmt.Println(convertToJapaneseCal(y, m, d, opt.Kanji))
    } else {
        fmt.Println(convertToChristianCal(y, m, d, era, opt.Kanji))
    }
}



