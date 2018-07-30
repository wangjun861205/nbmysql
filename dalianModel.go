package nbmysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var BkDalian *sql.DB

func init() {
	db, err := sql.Open("mysql", "wangjun:Wt20110523@tcp(127.0.0.1:12345)/bk_dalian")
	if err != nil {
		panic(err)
	}
	BkDalian = db
}

var BookMap = map[string]string{
	"@Id":           "id",
	"@Title":        "title",
	"@Price":        "price",
	"@Author":       "author",
	"@Publisher":    "publisher",
	"@Series":       "series",
	"@Isbn":         "isbn",
	"@PublishDate":  "publish_date",
	"@Binding":      "binding",
	"@Format":       "format",
	"@Pages":        "pages",
	"@WordCount":    "word_count",
	"@ContentIntro": "content_intro",
	"@AuthorIntro":  "author_intro",
	"@Menu":         "menu",
	"@Volume":       "volume",
	"@Category":     "category",
	"@Count":        "count",
	"@Code":         "code",
}

type Book struct {
	Id           *int64
	Title        *string
	Price        *int64
	Author       *string
	Publisher    *string
	Series       *string
	Isbn         *string
	PublishDate  *time.Time
	Binding      *string
	Format       *string
	Pages        *int64
	WordCount    *int64
	ContentIntro *string
	AuthorIntro  *string
	Menu         *string
	Volume       *int64
	Category     *string
	Count        *int64
	Code         *string
}

type BookToTags struct {
	All    func() ([]*Tags, error)
	Filter func(query string) ([]*Tags, error)
}

func (m *Book) Tags() BookToTags {
	return BookToTags{
		All: func() ([]*Tags, error) {
			rows, err := BkDalian.Query("SELECT `tags`.* FROM `book` JOIN `book__tags` ON `book`.`isbn`=`book__tags`.`book__isbn` JOIN `tags` on `book__tags`.`tags__id` = `tags`.`id` WHERE `book`.`isbn` = ?", *m.Isbn)
			if err != nil {
				return nil, err
			}
			list := make([]*Tags, 0, 256)
			for rows.Next() {
				model, err := TagsFromRows(rows)
				if err != nil {
					return nil, err
				}
				list = append(list, model)
			}
			return list, nil
		},
		Filter: func(query string) ([]*Tags, error) {
			for k, v := range TagsMap {
				query = strings.Replace(query, k, v, -1)
			}
			rows, err := BkDalian.Query("SELECT `tags`.* FROM `book` JOIN `book__tags` ON `book`.`isbn`=`book__tags`.`book__isbn` JOIN `tags` on `book__tags`.`tags__id` = `tags`.`id` WHERE `book`.`isbn` = ? AND ?", *m.Isbn, query)
			if err != nil {
				return nil, err
			}
			list := make([]*Tags, 0, 256)
			for rows.Next() {
				model, err := TagsFromRows(rows)
				if err != nil {
					return nil, err
				}
				list = append(list, model)
			}
			return list, nil
		},
	}
}
func NewBook(id *int64, title *string, price *int64, author *string, publisher *string, series *string, isbn *string, publishDate *time.Time, binding *string, format *string, pages *int64, wordCount *int64, contentIntro *string, authorIntro *string, menu *string, volume *int64, category *string, count *int64, code *string) *Book {
	book := &Book{id, title, price, author, publisher, series, isbn, publishDate, binding, format, pages, wordCount, contentIntro, authorIntro, menu, volume, category, count, code}
	return book
}
func AllBook() ([]*Book, error) {
	rows, err := BkDalian.Query("SELECT * FROM `book`")
	if err != nil {
		return nil, err
	}
	list := make([]*Book, 0, 256)
	for rows.Next() {
		model, err := BookFromRows(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, model)
	}
	return list, nil
}
func QueryBook(query string) ([]*Book, error) {
	for k, v := range BookMap {
		query = strings.Replace(query, k, v, -1)
	}
	rows, err := BkDalian.Query("SELECT * FROM `book` WHERE ?", query)
	if err != nil {
		return nil, err
	}
	list := make([]*Book, 0, 256)
	for rows.Next() {
		model, err := BookFromRows(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, model)
	}
	return list, nil
}
func (m *Book) Insert() error {
	colList := make([]string, 0, 32)
	valList := make([]string, 0, 32)
	if m.Id != nil {
		colList = append(colList, "id")
		valList = append(valList, fmt.Sprintf("%d", *m.Id))
	}
	if m.Title != nil {
		colList = append(colList, "title")
		valList = append(valList, fmt.Sprintf("%q", *m.Title))
	}
	if m.Price != nil {
		colList = append(colList, "price")
		valList = append(valList, fmt.Sprintf("%d", *m.Price))
	}
	if m.Author != nil {
		colList = append(colList, "author")
		valList = append(valList, fmt.Sprintf("%q", *m.Author))
	}
	if m.Publisher != nil {
		colList = append(colList, "publisher")
		valList = append(valList, fmt.Sprintf("%q", *m.Publisher))
	}
	if m.Series != nil {
		colList = append(colList, "series")
		valList = append(valList, fmt.Sprintf("%q", *m.Series))
	}
	if m.Isbn != nil {
		colList = append(colList, "isbn")
		valList = append(valList, fmt.Sprintf("%q", *m.Isbn))
	}
	if m.PublishDate != nil {
		colList = append(colList, "publish_date")
		valList = append(valList, fmt.Sprintf("%q", m.PublishDate.Format("2006-01-02 15:04:05")))
	}
	if m.Binding != nil {
		colList = append(colList, "binding")
		valList = append(valList, fmt.Sprintf("%q", *m.Binding))
	}
	if m.Format != nil {
		colList = append(colList, "format")
		valList = append(valList, fmt.Sprintf("%q", *m.Format))
	}
	if m.Pages != nil {
		colList = append(colList, "pages")
		valList = append(valList, fmt.Sprintf("%d", *m.Pages))
	}
	if m.WordCount != nil {
		colList = append(colList, "word_count")
		valList = append(valList, fmt.Sprintf("%d", *m.WordCount))
	}
	if m.ContentIntro != nil {
		colList = append(colList, "content_intro")
		valList = append(valList, fmt.Sprintf("%q", *m.ContentIntro))
	}
	if m.AuthorIntro != nil {
		colList = append(colList, "author_intro")
		valList = append(valList, fmt.Sprintf("%q", *m.AuthorIntro))
	}
	if m.Menu != nil {
		colList = append(colList, "menu")
		valList = append(valList, fmt.Sprintf("%q", *m.Menu))
	}
	if m.Volume != nil {
		colList = append(colList, "volume")
		valList = append(valList, fmt.Sprintf("%d", *m.Volume))
	}
	if m.Category != nil {
		colList = append(colList, "category")
		valList = append(valList, fmt.Sprintf("%q", *m.Category))
	}
	if m.Count != nil {
		colList = append(colList, "count")
		valList = append(valList, fmt.Sprintf("%d", *m.Count))
	}
	if m.Code != nil {
		colList = append(colList, "code")
		valList = append(valList, fmt.Sprintf("%q", *m.Code))
	}
	_, err := BkDalian.Exec("INSERT INTO `book` (?) VALUES (?)", strings.Join(colList, ", "), strings.Join(valList, ", "))
	if err != nil {
		return nil
	}
	lastInsertId := GetLastId(BkDalian)
	m.Id = &lastInsertId
	return nil
}
func (m *Book) Update() error {
	colList := make([]string, 0, 32)
	valList := make([]string, 0, 32)
	if m.Id != nil {
		colList = append(colList, "id")
		valList = append(valList, fmt.Sprintf("%d", *m.Id))
	}
	if m.Title != nil {
		colList = append(colList, "title")
		valList = append(valList, fmt.Sprintf("%q", *m.Title))
	}
	if m.Price != nil {
		colList = append(colList, "price")
		valList = append(valList, fmt.Sprintf("%d", *m.Price))
	}
	if m.Author != nil {
		colList = append(colList, "author")
		valList = append(valList, fmt.Sprintf("%q", *m.Author))
	}
	if m.Publisher != nil {
		colList = append(colList, "publisher")
		valList = append(valList, fmt.Sprintf("%q", *m.Publisher))
	}
	if m.Series != nil {
		colList = append(colList, "series")
		valList = append(valList, fmt.Sprintf("%q", *m.Series))
	}
	if m.Isbn != nil {
		colList = append(colList, "isbn")
		valList = append(valList, fmt.Sprintf("%q", *m.Isbn))
	}
	if m.PublishDate != nil {
		colList = append(colList, "publish_date")
		valList = append(valList, fmt.Sprintf("%q", m.PublishDate.Format("2006-01-02 15:04:05")))
	}
	if m.Binding != nil {
		colList = append(colList, "binding")
		valList = append(valList, fmt.Sprintf("%q", *m.Binding))
	}
	if m.Format != nil {
		colList = append(colList, "format")
		valList = append(valList, fmt.Sprintf("%q", *m.Format))
	}
	if m.Pages != nil {
		colList = append(colList, "pages")
		valList = append(valList, fmt.Sprintf("%d", *m.Pages))
	}
	if m.WordCount != nil {
		colList = append(colList, "word_count")
		valList = append(valList, fmt.Sprintf("%d", *m.WordCount))
	}
	if m.ContentIntro != nil {
		colList = append(colList, "content_intro")
		valList = append(valList, fmt.Sprintf("%q", *m.ContentIntro))
	}
	if m.AuthorIntro != nil {
		colList = append(colList, "author_intro")
		valList = append(valList, fmt.Sprintf("%q", *m.AuthorIntro))
	}
	if m.Menu != nil {
		colList = append(colList, "menu")
		valList = append(valList, fmt.Sprintf("%q", *m.Menu))
	}
	if m.Volume != nil {
		colList = append(colList, "volume")
		valList = append(valList, fmt.Sprintf("%d", *m.Volume))
	}
	if m.Category != nil {
		colList = append(colList, "category")
		valList = append(valList, fmt.Sprintf("%q", *m.Category))
	}
	if m.Count != nil {
		colList = append(colList, "count")
		valList = append(valList, fmt.Sprintf("%d", *m.Count))
	}
	if m.Code != nil {
		colList = append(colList, "code")
		valList = append(valList, fmt.Sprintf("%q", *m.Code))
	}
	updateList := make([]string, 0, 32)
	for i := 0; i < len(colList); i++ {
		updateList = append(updateList, fmt.Sprintf("%s=%s", colList[i], valList[i]))
	}
	_, err := BkDalian.Exec("UPDATE `book` SET ? WHERE id = ?", strings.Join(updateList, ", "), *m.Id)
	return err
}
func (m *Book) Delete() error {
	_, err := BkDalian.Exec("DELETE FROM `book` where 'id' = ?", *m.Id)
	return err
}
func BookFromRows(rows *sql.Rows) (*Book, error) {
	_id := new(Int)
	_title := new(String)
	_price := new(Int)
	_author := new(String)
	_publisher := new(String)
	_series := new(String)
	_isbn := new(String)
	_publishDate := new(Time)
	_binding := new(String)
	_format := new(String)
	_pages := new(Int)
	_wordCount := new(Int)
	_contentIntro := new(String)
	_authorIntro := new(String)
	_menu := new(String)
	_volume := new(Int)
	_category := new(String)
	_count := new(Int)
	_code := new(String)
	err := rows.Scan(_id, _title, _price, _author, _publisher, _series, _isbn, _publishDate, _binding, _format, _pages, _wordCount, _contentIntro, _authorIntro, _menu, _volume, _category, _count, _code)
	if err != nil {
		return nil, err
	}
	var (
		id           *int64
		title        *string
		price        *int64
		author       *string
		publisher    *string
		series       *string
		isbn         *string
		publishDate  *time.Time
		binding      *string
		format       *string
		pages        *int64
		wordCount    *int64
		contentIntro *string
		authorIntro  *string
		menu         *string
		volume       *int64
		category     *string
		count        *int64
		code         *string
	)
	if !_id.IsNull {
		id = &_id.Value
	}
	if !_title.IsNull {
		title = &_title.Value
	}
	if !_price.IsNull {
		price = &_price.Value
	}
	if !_author.IsNull {
		author = &_author.Value
	}
	if !_publisher.IsNull {
		publisher = &_publisher.Value
	}
	if !_series.IsNull {
		series = &_series.Value
	}
	if !_isbn.IsNull {
		isbn = &_isbn.Value
	}
	if !_publishDate.IsNull {
		publishDate = &_publishDate.Value
	}
	if !_binding.IsNull {
		binding = &_binding.Value
	}
	if !_format.IsNull {
		format = &_format.Value
	}
	if !_pages.IsNull {
		pages = &_pages.Value
	}
	if !_wordCount.IsNull {
		wordCount = &_wordCount.Value
	}
	if !_contentIntro.IsNull {
		contentIntro = &_contentIntro.Value
	}
	if !_authorIntro.IsNull {
		authorIntro = &_authorIntro.Value
	}
	if !_menu.IsNull {
		menu = &_menu.Value
	}
	if !_volume.IsNull {
		volume = &_volume.Value
	}
	if !_category.IsNull {
		category = &_category.Value
	}
	if !_count.IsNull {
		count = &_count.Value
	}
	if !_code.IsNull {
		code = &_code.Value
	}
	return NewBook(id, title, price, author, publisher, series, isbn, publishDate, binding, format, pages, wordCount, contentIntro, authorIntro, menu, volume, category, count, code), nil
}

var TagsMap = map[string]string{
	"@Id":  "id",
	"@Tag": "tag",
}

type Tags struct {
	Id  *int64
	Tag *string
}
type TagsToBook struct {
	All    func() ([]*Book, error)
	Filter func(query string) ([]*Book, error)
}

func (m *Tags) Book() TagsToBook {
	return TagsToBook{
		All: func() ([]*Book, error) {
			rows, err := BkDalian.Query("SELECT `book`.* FROM `tags` JOIN `book__tags` ON `tags`.`id`=`book__tags`.`tags__id` JOIN `book` on `book__tags`.`book__isbn` = `book`.`isbn` WHERE `tags`.`id` = ?", *m.Id)
			if err != nil {
				return nil, err
			}
			list := make([]*Book, 0, 256)
			for rows.Next() {
				model, err := BookFromRows(rows)
				if err != nil {
					return nil, err
				}
				list = append(list, model)
			}
			return list, nil
		},
		Filter: func(query string) ([]*Book, error) {
			for k, v := range BookMap {
				query = strings.Replace(query, k, v, -1)
			}
			rows, err := BkDalian.Query("SELECT `book`.* FROM `tags` JOIN `book__tags` ON `tags`.`id`=`book__tags`.`tags__id` JOIN `book` on `book__tags`.`book__isbn` = `book`.`isbn` WHERE `tags`.`id` = ? AND ?", *m.Id, query)
			if err != nil {
				return nil, err
			}
			list := make([]*Book, 0, 256)
			for rows.Next() {
				model, err := BookFromRows(rows)
				if err != nil {
					return nil, err
				}
				list = append(list, model)
			}
			return list, nil
		},
	}
}
func NewTags(id *int64, tag *string) *Tags {
	tags := &Tags{id, tag}
	return tags
}
func AllTags() ([]*Tags, error) {
	rows, err := BkDalian.Query("SELECT * FROM `tags`")
	if err != nil {
		return nil, err
	}
	list := make([]*Tags, 0, 256)
	for rows.Next() {
		model, err := TagsFromRows(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, model)
	}
	return list, nil
}
func QueryTags(query string) ([]*Tags, error) {
	for k, v := range TagsMap {
		query = strings.Replace(query, k, v, -1)
	}
	rows, err := BkDalian.Query("SELECT * FROM `tags` WHERE ?", query)
	if err != nil {
		return nil, err
	}
	list := make([]*Tags, 0, 256)
	for rows.Next() {
		model, err := TagsFromRows(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, model)
	}
	return list, nil
}
func (m *Tags) Insert() error {
	colList := make([]string, 0, 32)
	valList := make([]string, 0, 32)
	if m.Id != nil {
		colList = append(colList, "id")
		valList = append(valList, fmt.Sprintf("%d", *m.Id))
	}
	if m.Tag != nil {
		colList = append(colList, "tag")
		valList = append(valList, fmt.Sprintf("%q", *m.Tag))
	}
	_, err := BkDalian.Exec("INSERT INTO `tags` (?) VALUES (?)", strings.Join(colList, ", "), strings.Join(valList, ", "))
	if err != nil {
		return nil
	}
	lastInsertId := GetLastId(BkDalian)
	m.Id = &lastInsertId
	return nil
}
func (m *Tags) Update() error {
	colList := make([]string, 0, 32)
	valList := make([]string, 0, 32)
	if m.Id != nil {
		colList = append(colList, "id")
		valList = append(valList, fmt.Sprintf("%d", *m.Id))
	}
	if m.Tag != nil {
		colList = append(colList, "tag")
		valList = append(valList, fmt.Sprintf("%q", *m.Tag))
	}
	updateList := make([]string, 0, 32)
	for i := 0; i < len(colList); i++ {
		updateList = append(updateList, fmt.Sprintf("%s=%s", colList[i], valList[i]))
	}
	_, err := BkDalian.Exec("UPDATE `tags` SET ? WHERE id = ?", strings.Join(updateList, ", "), *m.Id)
	return err
}
func (m *Tags) Delete() error {
	_, err := BkDalian.Exec("DELETE FROM `tags` where 'id' = ?", *m.Id)
	return err
}
func TagsFromRows(rows *sql.Rows) (*Tags, error) {
	_id := new(Int)
	_tag := new(String)
	err := rows.Scan(_id, _tag)
	if err != nil {
		return nil, err
	}
	var (
		id  *int64
		tag *string
	)
	if !_id.IsNull {
		id = &_id.Value
	}
	if !_tag.IsNull {
		tag = &_tag.Value
	}
	return NewTags(id, tag), nil
}
