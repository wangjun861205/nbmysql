package nbmysql

// import (
// 	"database/sql"
// 	"fmt"
// 	"strings"
// 	"time"
// )

// var BkDalianDB *sql.DB

// func init() {
// 	db, err := sql.Open("mysql", "wangjun:Wt20110523@tcp(127.0.0.1:12345)/bk_dalian")
// 	if err != nil {
// 		panic(err)
// 	}
// 	BkDalianDB = db
// }

// type Book struct {
// 	Id           *int64
// 	Title        *string
// 	Price        *int64
// 	Author       *string
// 	Publisher    *string
// 	Series       *string
// 	Isbn         *string
// 	PublishDate  *time.Time
// 	Binding      *string
// 	Format       *string
// 	Pages        *int64
// 	WordCount    *int64
// 	ContentIntro *string
// 	AuthorIntro  *string
// 	Menu         *string
// 	Volume       *int64
// 	Category     *string
// 	Count        *int64
// 	Code         *string
// }

// type BookToTag struct {
// 	All    func() ([]*Tag, error)
// 	Filter func(query string) ([]*Tag, error)
// }

// func NewBook(id *int64, title *string, price *int64, author *string, publisher *string, series *string, isbn *string, publishDate *time.Time,
// 	binding *string, format *string, pages *int64, wordCount *int64, contentIntro *string, authorIntro *string, menu *string, volume *int64,
// 	category *string, count *int64, code *string) *Book {
// 	book := &Book{id, title, price, author, publisher, series, isbn, publishDate, binding, format, pages, wordCount, contentIntro,
// 		authorIntro, menu, volume, category, count, code}
// 	return book
// }

// func AllBook() ([]*Book, error) {
// 	rows, err := BkDalianDB.Query("SELECT * FROM `book`")
// 	if err != nil {
// 		return nil, err
// 	}
// 	bookList := make([]*Book, 0, 256)
// 	for rows.Next() {
// 		book, err := bookFromRows(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		bookList = append(bookList, book)
// 	}
// 	return bookList, nil
// }

// func QueryBook(query string) ([]*Book, error) {
// 	rows, err := BkDalianDB.Query("SELECT * from book where ?", query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bookList := make([]*Book, 0, 256)
// 	for rows.Next() {
// 		book, err := bookFromRows(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		bookList = append(bookList, book)
// 	}
// 	return bookList, nil
// }

// func (b *Book) Tag() BookToTag {
// 	return BookToTag{
// 		All: func() ([]*Tag, error) {
// 			rows, err := BkDalianDB.Query("SELECT tags.* from book join book_to_tags on book.isbn = book_to_tags.isbn join tags on book_to_tags.tag_id = tags.id where book.isbn = ?;", b.Isbn)
// 			if err != nil {
// 				return nil, err
// 			}
// 			bookList := make([]*Tag, 0, 256)
// 			for rows.Next() {
// 				book, err := tagFromRows(rows)
// 				if err != nil {
// 					return nil, err
// 				}
// 				bookList = append(bookList, book)
// 			}
// 			return bookList, nil
// 		},
// 		Filter: func(query string) ([]*Tag, error) {
// 			rows, err := BkDalianDB.Query("SELECT tags.* from book JOIN book__tag ON book.isbn = book__tag.isbn JOIN tag ON book__tag.tag_id = tag.id where book.isbn = ? and ?", b.Isbn, query)
// 			if err != nil {
// 				return nil, err
// 			}
// 			tagList := make([]*Tag, 0, 256)
// 			for rows.Next() {
// 				tag, err := tagFromRows(rows)
// 				if err != nil {
// 					return nil, err
// 				}
// 				tagList = append(tagList, tag)
// 			}
// 			return tagList, nil
// 		},
// 	}
// }

// func (b *Book) Insert() error {
// 	colList := make([]string, 0, 32)
// 	valList := make([]string, 0, 32)
// 	if b.Title != nil {
// 		colList = append(colList, "title")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Title))
// 	}
// 	if b.Price != nil {
// 		colList = append(colList, "price")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Price))
// 	}
// 	if b.Author != nil {
// 		colList = append(colList, "author")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Author))
// 	}
// 	if b.Publisher != nil {
// 		colList = append(colList, "publisher")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Publisher))
// 	}
// 	if b.Series != nil {
// 		colList = append(colList, "series")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Series))
// 	}
// 	if b.Isbn != nil {
// 		colList = append(colList, "isbn")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Isbn))
// 	}
// 	if b.PublishDate != nil {
// 		colList = append(colList, "publish_date")
// 		valList = append(valList, fmt.Sprintf("%q", b.PublishDate.Format("2006-01-02 15:04:05")))
// 	}
// 	if b.Binding != nil {
// 		colList = append(colList, "binding")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Binding))
// 	}
// 	if b.Format != nil {
// 		colList = append(colList, "format")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Format))
// 	}
// 	if b.Pages != nil {
// 		colList = append(colList, "pages")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Pages))
// 	}
// 	if b.WordCount != nil {
// 		colList = append(colList, "word_count")
// 		valList = append(valList, fmt.Sprintf("%d", *b.WordCount))
// 	}
// 	if b.ContentIntro != nil {
// 		colList = append(colList, "content_intro")
// 		valList = append(valList, fmt.Sprintf("%q", *b.ContentIntro))
// 	}
// 	if b.AuthorIntro != nil {
// 		colList = append(colList, "author_intro")
// 		valList = append(valList, fmt.Sprintf("%q", *b.AuthorIntro))
// 	}
// 	if b.Menu != nil {
// 		colList = append(colList, "menu")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Menu))
// 	}
// 	if b.Volume != nil {
// 		colList = append(colList, "volume")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Volume))
// 	}
// 	if b.Category != nil {
// 		colList = append(colList, "category")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Category))
// 	}
// 	if b.Count != nil {
// 		colList = append(colList, "count")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Count))
// 	}
// 	if b.Code != nil {
// 		colList = append(colList, "code")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Code))
// 	}
// 	_, err := BkDalianDB.Exec("INSERT INTO `book` (?) VALUES (?)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 	if err != nil {
// 		return nil
// 	}
// 	lastInsertId := GetLastId(BkDalianDB)
// 	b.Id = &lastInsertId
// 	return nil
// }

// func (b *Book) Update() error {
// 	colList := make([]string, 0, 32)
// 	valList := make([]string, 0, 32)
// 	if b.Title != nil {
// 		colList = append(colList, "title")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Title))
// 	}
// 	if b.Price != nil {
// 		colList = append(colList, "price")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Price))
// 	}
// 	if b.Author != nil {
// 		colList = append(colList, "author")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Author))
// 	}
// 	if b.Publisher != nil {
// 		colList = append(colList, "publisher")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Publisher))
// 	}
// 	if b.Series != nil {
// 		colList = append(colList, "series")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Series))
// 	}
// 	if b.Isbn != nil {
// 		colList = append(colList, "isbn")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Isbn))
// 	}
// 	if b.PublishDate != nil {
// 		colList = append(colList, "publish_date")
// 		valList = append(valList, fmt.Sprintf("%q", b.PublishDate.Format("2006-01-02 15:04:05")))
// 	}
// 	if b.Binding != nil {
// 		colList = append(colList, "binding")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Binding))
// 	}
// 	if b.Format != nil {
// 		colList = append(colList, "format")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Format))
// 	}
// 	if b.Pages != nil {
// 		colList = append(colList, "pages")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Pages))
// 	}
// 	if b.WordCount != nil {
// 		colList = append(colList, "word_count")
// 		valList = append(valList, fmt.Sprintf("%d", *b.WordCount))
// 	}
// 	if b.ContentIntro != nil {
// 		colList = append(colList, "content_intro")
// 		valList = append(valList, fmt.Sprintf("%q", *b.ContentIntro))
// 	}
// 	if b.AuthorIntro != nil {
// 		colList = append(colList, "author_intro")
// 		valList = append(valList, fmt.Sprintf("%q", *b.AuthorIntro))
// 	}
// 	if b.Menu != nil {
// 		colList = append(colList, "menu")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Menu))
// 	}
// 	if b.Volume != nil {
// 		colList = append(colList, "volume")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Volume))
// 	}
// 	if b.Category != nil {
// 		colList = append(colList, "category")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Category))
// 	}
// 	if b.Count != nil {
// 		colList = append(colList, "count")
// 		valList = append(valList, fmt.Sprintf("%d", *b.Count))
// 	}
// 	if b.Code != nil {
// 		colList = append(colList, "code")
// 		valList = append(valList, fmt.Sprintf("%q", *b.Code))
// 	}
// 	updateList := make([]string, 0, 32)
// 	for i := 0; i < len(colList); i++ {
// 		updateList = append(updateList, fmt.Sprintf("%s=%s", colList[i], valList[i]))
// 	}
// 	_, err := BkDalianDB.Exec("UPDATE `book` SET ? WHERE `id` = ?", strings.Join(updateList, ", "), *b.Id)
// 	return err
// }

// func (b *Book) Delete() error {
// 	_, err := BkDalianDB.Exec("DELETE FROM `book` where `id` = ?", *b.Id)
// 	return err
// }

// func bookFromRows(rows *sql.Rows) (*Book, error) {
// 	_id := new(Int)
// 	_title := new(String)
// 	_price := new(Int)
// 	_author := new(String)
// 	_publisher := new(String)
// 	_series := new(String)
// 	_isbn := new(String)
// 	_publishDate := new(Time)
// 	_binding := new(String)
// 	_format := new(String)
// 	_pages := new(Int)
// 	_wordCount := new(Int)
// 	_contentIntro := new(String)
// 	_authorIntro := new(String)
// 	_menu := new(String)
// 	_volume := new(Int)
// 	_category := new(String)
// 	_count := new(Int)
// 	_code := new(String)
// 	err := rows.Scan(_id, _title, _price, _author, _publisher, _series, _isbn, _publishDate, _binding, _format, _pages, _wordCount, _contentIntro,
// 		_authorIntro, _menu, _volume, _category, _count, _code)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var (
// 		id           *int64
// 		title        *string
// 		price        *int64
// 		author       *string
// 		publisher    *string
// 		series       *string
// 		isbn         *string
// 		publishDate  *time.Time
// 		binding      *string
// 		format       *string
// 		pages        *int64
// 		wordCount    *int64
// 		contentIntro *string
// 		authorIntro  *string
// 		menu         *string
// 		volume       *int64
// 		category     *string
// 		count        *int64
// 		code         *string
// 	)
// 	if !_id.IsNull {
// 		id = &_id.Value
// 	}
// 	if !_title.IsNull {
// 		title = &_title.Value
// 	}
// 	if !_price.IsNull {
// 		price = &_price.Value
// 	}
// 	if !_author.IsNull {
// 		author = &_author.Value
// 	}
// 	if !_publisher.IsNull {
// 		publisher = &_publisher.Value
// 	}
// 	if !_series.IsNull {
// 		series = &_series.Value
// 	}
// 	if !_isbn.IsNull {
// 		isbn = &_isbn.Value
// 	}
// 	if !_publishDate.IsNull {
// 		publishDate = &_publishDate.Value
// 	}
// 	if !_binding.IsNull {
// 		binding = &_binding.Value
// 	}
// 	if !_format.IsNull {
// 		format = &_format.Value
// 	}
// 	if !_pages.IsNull {
// 		pages = &_pages.Value
// 	}
// 	if !_wordCount.IsNull {
// 		wordCount = &_wordCount.Value
// 	}
// 	if !_contentIntro.IsNull {
// 		contentIntro = &_contentIntro.Value
// 	}
// 	if !_authorIntro.IsNull {
// 		authorIntro = &_authorIntro.Value
// 	}
// 	if !_menu.IsNull {
// 		menu = &_menu.Value
// 	}
// 	if !_volume.IsNull {
// 		volume = &_volume.Value
// 	}
// 	if !_category.IsNull {
// 		category = &_category.Value
// 	}
// 	if !_count.IsNull {
// 		count = &_count.Value
// 	}
// 	if !_code.IsNull {
// 		code = &_code.Value
// 	}
// 	return NewBook(id, title, price, author, publisher, series, isbn, publishDate, binding, format, pages, wordCount, contentIntro, authorIntro,
// 		menu, volume, category, count, code), nil
// }

// type Tag struct {
// 	Id  *int64
// 	Tag *string
// }

// type TagToBook struct {
// 	All    func() ([]*Book, error)
// 	Filter func(query string) ([]*Book, error)
// }

// func NewTag(id *int64, tag *string) *Tag {
// 	return &Tag{id, tag}
// }

// func AllTag() ([]*Tag, error) {
// 	rows, err := BkDalianDB.Query("SELECT * FROM `tag`")
// 	if err != nil {
// 		return nil, err
// 	}
// 	tagList := make([]*Tag, 0, 256)
// 	for rows.Next() {
// 		tag, err := tagFromRows(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tagList = append(tagList, tag)
// 	}
// 	return tagList, nil
// }

// func QueryTag(query string) ([]*Tag, error) {
// 	rows, err := BkDalianDB.Query("SELECT * FROM `tag` where ?", query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tagList := make([]*Tag, 0, 256)
// 	for rows.Next() {
// 		tag, err := tagFromRows(rows)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tagList = append(tagList, tag)
// 	}
// 	return tagList, nil
// }

// func tagFromRows(rows *sql.Rows) (*Tag, error) {
// 	_id := new(Int)
// 	_tag := new(String)
// 	var id *int64
// 	var tag *string
// 	err := rows.Scan(_id, _tag)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !_id.IsNull {
// 		id = &_id.Value
// 	}
// 	if !_tag.IsNull {
// 		tag = &_tag.Value
// 	}
// 	return NewTag(id, tag), nil
// }

// func (t *Tag) Book() TagToBook {
// 	return TagToBook{
// 		All: func() ([]*Book, error) {
// 			rows, err := BkDalianDB.Query("SELECT book.* FROM book JOIN book__tag ON book.isbn = book__tag.isbn JOIN tag on book__tag.tag_id = tag.id where tag.id = ?", t.Id)
// 			if err != nil {
// 				return nil, err
// 			}
// 			bookList := make([]*Book, 0, 256)
// 			for rows.Next() {
// 				book, err := bookFromRows(rows)
// 				if err != nil {
// 					return nil, err
// 				}
// 				bookList = append(bookList, book)
// 			}
// 			return bookList, nil
// 		},
// 		Filter: func(query string) ([]*Book, error) {
// 			rows, err := BkDalianDB.Query("SELECT book.* FROM book JOIN book__tag ON book.isbn = book__tag.isbn JOIN tag ON book__tag.tag_id = tag.id where tag.id = ? and ?", t.Id, query)
// 			if err != nil {
// 				return nil, err
// 			}
// 			bookList := make([]*Book, 0, 256)
// 			for rows.Next() {
// 				book, err := bookFromRows(rows)
// 				if err != nil {
// 					return nil, err
// 				}
// 				bookList = append(bookList, book)
// 			}
// 			return bookList, nil
// 		},
// 	}
// }

// func (t *Tag) Insert() error {
// 	colList := make([]string, 0, 32)
// 	valList := make([]string, 0, 32)
// 	if t.Tag != nil {
// 		colList = append(colList, "tag")
// 		valList = append(valList, fmt.Sprintf("%q", *t.Tag))
// 	}
// 	_, err := BkDalianDB.Exec("INSERT INTO `tag` (?) VALUES (?)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 	if err != nil {
// 		return err
// 	}
// 	lastInsertId := GetLastId(BkDalianDB)
// 	t.Id = &lastInsertId
// 	return nil
// }

// func (t *Tag) Update() error {
// 	colList := make([]string, 0, 32)
// 	valList := make([]string, 0, 32)
// 	if t.Tag != nil {
// 		colList = append(colList, "tag")
// 		valList = append(valList, fmt.Sprintf("%q", *t.Tag))
// 	}
// 	updateList := make([]string, 0, 32)
// 	for i := 0; i < len(colList); i++ {
// 		updateList = append(updateList, fmt.Sprintf("%s=%s", colList[i], valList[i]))
// 	}
// 	_, err := BkDalianDB.Exec("UPDATE `tag` SET ? WHERE `id` = ?", strings.Join(updateList, ", "), *t.Id)
// 	return err
// }

// func (t *Tag) Delete() error {
// 	_, err := BkDalianDB.Exec("DELETE FROM `tag` WHERE `id` = ?", *t.Id)
// 	return err
// }
