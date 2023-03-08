package booklib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/ilaziness/mdbook/internal/util"
)

// 单个书
type Book struct {
	Title   string `json:"title"`
	DirName string
	// 目录
	Catalogues CatalogueList `json:"catalogue"`
}

type CatalogueList []Catalogue

type Catalogue struct {
	Title    string      `json:"title"`
	Path     string      `json:"path"`
	Children []Catalogue `json:"children"`
	Parent   int         `json:"-"`
}

// 书库
var (
	Lib          = map[string]Book{}
	chapterMatch = regexp.MustCompile(`^( *)[\*\-].+?\[(.+)\]\((.+)\)`)
	defaultDir   = "book"
)

const (
	summaryFile = "SUMMARY.md"
)

// ScanBook 扫描目录，初始化书库
func ScanBook(dir string) error {
	util.Info("开始扫描书库内容...")
	if dir == "" {
		dir = defaultDir
	}
	dirinfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !dirinfo.IsDir() {
		return errors.New("BookDir is not a directory")
	}

	dfs := os.DirFS(dir)
	bookDir, err := fs.ReadDir(dfs, ".")
	if err != nil {
		return err
	}
	for _, bd := range bookDir {
		if !bd.IsDir() {
			continue
		}
		summaryCt, err := fs.ReadFile(dfs, bd.Name()+"/"+summaryFile)
		if err != nil {
			log.Println(err)
			continue
		}
		parseSummary(bd.Name(), summaryCt)
	}

	return nil
}

func parseSummary(bookDir string, content []byte) {
	scaner := bufio.NewScanner(bytes.NewBuffer(content))
	bk := Book{DirName: bookDir}
	catalogue := make(CatalogueList, 0)
	chapters := make([]string, 0)
	for scaner.Scan() {
		line := scaner.Text()
		if len(line) == 0 {
			continue
		}
		if bk.Title == "" && strings.HasPrefix(line, "#") {
			// 书名
			bk.Title = strings.TrimSpace(line[1:])
			continue
		}
		if strings.HasPrefix(line, "#") {
			// 分节名
			catalogue = append(catalogue, Catalogue{
				Title: "part",
				Path:  strings.TrimSpace(line[1:]),
			})
			if len(chapters) > 0 {
				catalogue = append(catalogue, parseChapter(chapters)...)
			}
			chapters = make([]string, 0)
			continue
		}
		chapters = append(chapters, line)
	}
	if len(chapters) > 0 {
		catalogue = append(catalogue, parseChapter(chapters)...)
	}
	bk.Catalogues = catalogue
	if bk.Title == "" {
		bk.Title = bookDir
	}
	Lib[bk.Title] = bk
	fmt.Println(Lib)
}

func parseChapter(lines []string) CatalogueList {
	catalogues := make(CatalogueList, 0)
	catalogueMap := make(map[int]*Catalogue)
	level := make([]int, 0)
	keys := make([]int, 0)
	for key, line := range lines {
		mt := chapterMatch.FindStringSubmatch(line)
		if len(mt) != 4 {
			continue
		}
		currentLevel := strings.Count(mt[1], " ")
		parent := 0
		if len(level) == 0 {
			level = append(level, currentLevel)
			keys = append(keys, key)
		} else {
			lastIndex := len(level) - 1
			if level[lastIndex] == currentLevel {
				// 同级
				level = level[:lastIndex]
				keys = keys[:lastIndex]
				level = append(level, currentLevel)
				keys = append(keys, key)
				if len(keys) > 0 {
					parent = keys[len(keys)-1]
				}
			} else if level[lastIndex] > currentLevel {
				// 当前比上一条低级
				parent = keys[lastIndex]
				level = append(level, currentLevel)
				keys = append(keys, key)
			} else if level[lastIndex] < currentLevel {
				// 当前比上一条高级
				for {
					level = level[:lastIndex]
					keys = keys[:lastIndex]
					lastIndex -= 1
					if len(level) == 0 {
						parent = 0
						break
					}
					if level[lastIndex] == currentLevel {
						level = level[:lastIndex]
						keys = keys[:lastIndex]
						level = append(level, currentLevel)
						keys = append(keys, key)
						if len(keys) > 0 {
							parent = keys[len(keys)-1]
						} else {
							parent = 0
						}
						break
					}
					if level[lastIndex] > currentLevel {
						parent = keys[lastIndex]
						level = append(level, currentLevel)
						keys = append(keys, key)
						break
					}
					lastIndex = len(level) - 1
				}
			}
		}

		catalogueMap[key] = &Catalogue{Title: mt[2], Path: mt[3], Parent: parent}
	}
	for _, v := range catalogueMap {
		if v.Parent == 0 {
			catalogues = append(catalogues, *v)
			continue
		}
		if _, ok := catalogueMap[v.Parent]; ok {
			catalogueMap[v.Parent].Children = append(catalogueMap[v.Parent].Children, *v)
			continue
		}
	}
	return catalogues
}

func BookExist(bookname string) error {
	err := errors.New("数据不存在")
	if bookname == "" {
		return err
	}
	if _, ok := Lib[bookname]; !ok {
		return err
	}
	return nil
}

func GetPageContent(bookname string, pagePath string) string {
	file := filepath.Join(defaultDir, bookname, pagePath)
	fileInfo, err := os.Stat(file)
	if err != nil {
		log.Println("get page content error:", err, "file:", file)
		return ""
	}
	if fileInfo.IsDir() {
		log.Println("get page content error:", "file is dir", "file:", file)
		return ""
	}
	content, err := os.ReadFile(file)
	if err != nil {
		log.Println("read page content error:", err, "file:", file)
		return ""
	}
	if filepath.Ext(file) == ".html" {
		return string(content)
	}
	return string(markdown.ToHTML(content, nil, nil))
}
