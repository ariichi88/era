package main

import (
	"fmt"
	"os"
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
			} else if string(char) == "元" { // 元年を1に変換
				str = "1"
				isDigit = true
			}
		}
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

	// argsが1でないときもhelpをだして終了
    if len(args) == 0 || len(args) > 1 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

}



