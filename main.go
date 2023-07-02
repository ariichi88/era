package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"golang.org/x/text/width"

    flags "github.com/jessevdk/go-flags"
)

// 文字列から日付及び元号を取り出す
func getDateAndEra(date string) (y, m, d, era int) {

	var (
		isDigit     = false
		str         = ""
		dateStr     []string
	)

	// dateの中から各元号に対応する文字を見つけてeraに格納
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

    // dateの前処理
	// 数字が半角でないとうまくいかないため
	// 最後が文字でないと次の処理で最後の数字が取り出せないため
	date = width.Fold.String(date)
	date = date + "."

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
			} else if string(char) == "元" { // 元年を1に変換
				str = "1"
				isDigit = true
			}
		}
	}

	// dateStrの長さが３でないときは失敗
	if len(dateStr) != 3 {
       fmt.Println("日付の取得に失敗しました")
       os.Exit(1)
    }

	// 年月日をstringからintにする
	y, _ = strconv.Atoi(dateStr[0])
	m, _ = strconv.Atoi(dateStr[1])
	d, _ = strconv.Atoi(dateStr[2])

    // 結果を返す
    return y, m, d, era
}


// 和暦を西暦に変換
func convertToChristianCal(y, m , d, era int, kanji bool) string {

    var magicNumber int

    // eraから各元号の計算に使うmagicnumberを求める
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

    // 年月日を文字列にする
    year := strconv.Itoa(y + magicNumber)
    month := strconv.Itoa(m)
    date := strconv.Itoa(d)

    // 日付を作成して返す
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

    // 作成した日付と各年号の最初の日を比較して計算用のmagicnumberと
    // 元号の文字列（アルファベットと漢字）を求める
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

    // 元号＋年の作成
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

    // 月と日をstringにする
    month := strconv.Itoa(m)
    date := strconv.Itoa(d)

    // 日付を作成して返す
    if kanji {
        return year + "年" + month + "月" + date + "日"
    } else {
        return year + "." + month + "." + date
    }
}

func main() {

    // オプションの作成
    type option struct {
        Kanji bool `short:"k" long:"kanji" description:"日付をＸ年Ｘ月Ｘ日の形式で返します"`
    }

    // convertToChristianCalとconvertToJapaneseCalから帰ってきた文字列の保持用
    var resalt string
 
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

	// argsが1でないときもhelpをだして終了
    if len(args) == 0 || len(args) > 1 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

    // 日付と元号の取得
    y, m, d, era := getDateAndEra(args[0])

    // eraに入ってる値に応じて和暦にするのか西暦にするのか判定
    if era == 0 {
        resalt = convertToJapaneseCal(y, m, d, opt.Kanji)
    } else {
        resalt = convertToChristianCal(y, m, d, era, opt.Kanji)
    }

    // 結果の出力
    fmt.Println(resalt)
}



