package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jessevdk/go-flags"
)

// 和暦を西暦に変換
func convertToChristianCal(date string, isLong bool) {

	var (
		isDigit     = false
		str         = ""
		magicNumber int
		y, m, d     int
		dateStr     []string
	)

	// どの年号か判別してmagicNumberを求める
	if strings.ContainsAny(date, "明治MmＭｍ") {
		magicNumber = 1867 // 明治元年は西暦1868年
	} else if strings.ContainsAny(date, "大正TtＴｔ") {
		magicNumber = 1911 // 大正元年は西暦1912年
	} else if strings.ContainsAny(date, "昭和SsＳｓ") {
		magicNumber = 1925 // 昭和元年は西暦1926年
	} else if strings.ContainsAny(date, "平成HhＨｈ") {
		magicNumber = 1988 // 平成元年は西暦1989年
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

	// 日付を正しく取得できてなかった場合(dateStrの長さが３でない)は終了
	if len(dateStr) != 3 {
		fmt.Println("日付の取得に失敗しました")
		os.Exit(1)
	}

	// 日付を取得
	y, _ = strconv.Atoi(dateStr[0])
	m, _ = strconv.Atoi(dateStr[1])
	d, _ = strconv.Atoi(dateStr[2])

	// 西暦の日付を返す
	if isLong {
		fmt.Printf("%d年%d月%d日\n", y+magicNumber, m, d)
	} else {
		fmt.Printf("%d/%d/%d\n", y+magicNumber, m, d)
	}

}

// 西暦を和暦に変換
func convertToJapaneseCal(date string, isLong, isGrif bool) {

	var (
		t           time.Time
		magicNumber int
		eraNameL    string
		eraNameS    string
		y, m, d     int
	)

	// dateから日付を生成する
	t, err := time.Parse("2006/01/02", date)
	if err != nil {
		t, err = time.Parse("2006/1/2", date)
		if err != nil {
			t, err = time.Parse("2006年01月02日", date)
			if err != nil {
				t, err = time.Parse("2006年1月2日", date)
				if err != nil {
					fmt.Println("日付の取得に失敗しました")
					os.Exit(1)
				}
			}
		}
	}

	// 元号を判定する
	if t.Before(time.Date(1912, 7, 30, 0, 0, 0, 0, time.Local)) {
		magicNumber = 1867
		eraNameL = "明治"
		eraNameS = "M"
	} else if t.Before(time.Date(1926, 12, 25, 0, 0, 0, 0, time.Local)) {
		magicNumber = 1911
		eraNameL = "大正"
		eraNameS = "T"
	} else if t.Before(time.Date(1989, 1, 8, 0, 0, 0, 0, time.Local)) {
		magicNumber = 1925
		eraNameL = "昭和"
		eraNameS = "S"
	} else if t.Before(time.Date(2019, 5, 1, 0, 0, 0, 0, time.Local)) {
		magicNumber = 1988
		eraNameL = "平成"
		eraNameS = "H"
	}

	// 和暦の年月日の取得
	y = t.Year() - magicNumber
	m = int(t.Month())
	d = t.Day()

	// 和暦の年月日の表示
	if isLong {
		if isGrif && y == 1 {
			fmt.Printf("%s元年%d月%d日\n", eraNameL, m, d)
		} else {
			fmt.Printf("%s%d年%d月%d日\n", eraNameL, y, m, d)
		}
	} else {
		fmt.Printf("%s%d.%d.%d\n", eraNameS, y, m, d)
	}

}

type options struct {
	Jpn  bool `short:"j" description:"西暦を和暦に変換します"`
	Chr  bool `short:"c" description:"和暦を西暦に変換します"`
	Long bool `short:"l" description:"日付をy年m月d日で表示します"`
	Grif bool `short:"g" description:"和暦の1年を元年で表示します(ｊオプションの時のみ)"`
}

var opt options

var parser = flags.NewParser(&opt, flags.Default)

func main() {

	//　コマンド名と使用法の設定
	parser.Name = "Era"
	parser.Usage = "[OPTIONS] Date"

	//　オプションのパース
	args, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	// コマンドライン引数がないときはヘルプを出して終了
	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	// オプションの解析とコマンドの実行
	if opt.Jpn && opt.Chr {
		fmt.Println("ｊオプションとｃオプションは同時に指定できません")
		os.Exit(1)
	} else if opt.Jpn {
		convertToJapaneseCal(args[0], opt.Long, opt.Grif)
	} else if opt.Chr {
		convertToChristianCal(args[0], opt.Long)
	}
}
