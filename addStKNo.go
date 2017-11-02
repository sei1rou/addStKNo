package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	//ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Print("Start\r\n")

	// ファイルを読み込んで二次元配列にいれる
	records := readFile(flag.Arg(0))

	// レコードの末尾にコンマを追加してファイルへ書き出す
	saveFile(flag.Arg(0), records)

	log.Print("Finesh !\r\n")

}

func readFile(fileName string) [][]string {
	//入力ファイル準備
	infile, err := os.Open(fileName)
	failOnError(err)
	defer infile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))

	//CSVファイルを２次元配列に展開
	readRecords := make([][]string, 0)
	for {
		record, err := reader.Read() //１行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		readRecords = append(readRecords, record)
	}

	return readRecords
}

func saveFile(fileName string, saveRecords [][]string) {
	//出力ファイル準備
	outDir, outFileName := filepath.Split(fileName)
	pos := strings.LastIndex(outFileName, ".")
	outfile, err := os.Create(outDir + outFileName[:pos] + ".csv")
	failOnError(err)
	defer outfile.Close()

	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	//writer.Comma = '\t'
	writer.UseCRLF = true

	for i, outRecord := range saveRecords {
		if i == 0 {
			outRecord = append(outRecord, "健診番号")
		} else {
			outRecord = append(outRecord, "")
		}
		writer.Write(outRecord)
	}

	writer.Flush()

}
