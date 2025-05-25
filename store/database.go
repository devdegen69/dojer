package store

import (
	"dojer/utils"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

const LIMIT = 40

type ListRequest struct {
	Page  int
	Seed  string
	Limit int
}

// Struct for the table
type Doujinshi struct {
	ID         string    `gorm:"primarykey;autoIncrement:false;index" json:"id"`
	Name       string    `json:"name"`
	Title      string    `json:"title"`
	Parodies   string    `json:"parodies"`
	Characters string    `json:"characters"`
	Tags       string    `json:"tags"`
	Artists    string    `json:"artists"`
	Groups     string    `json:"groups"`
	Languages  string    `json:"languages"`
	Categories string    `json:"categories"`
	Pages      int       `json:"pages"`
	Uploaded   string    `json:"uploaded"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (d Doujinshi) String() string {
	return fmt.Sprintf(
		"Doujinshi {\n"+
			"  ID: %s\n"+
			"  Name: %s\n"+
			"  Title: %s\n"+
			"  Parodies: %s\n"+
			"  Characters: %s\n"+
			"  Tags: %s\n"+
			"  Artists: %s\n"+
			"  Groups: %s\n"+
			"  Languages: %s\n"+
			"  Categories: %s\n"+
			"  Pages: %d\n"+
			"  Uploaded: %s\n"+
			"  CreatedAt: %s\n"+
			"}",
		d.ID,
		d.Name,
		d.Title,
		d.Parodies,
		d.Characters,
		d.Tags,
		d.Artists,
		d.Groups,
		d.Languages,
		d.Categories,
		d.Pages,
		d.Uploaded,
		d.CreatedAt.Format(time.RFC3339),
	)
}

type Pagination struct {
	CurrentPage  int   `json:"currentPage"`
	TotalPages   int   `json:"totalPages"`
	TotalResults int64 `json:"totalResults"`
	PageSize     int   `json:"pageSize"`
	Pages        []int `json:"pages"`
}

func gormInstance() *gorm.DB {
	databasePath := utils.GetDataPath("database.sqlite3")
	_ = utils.EnsureExists(databasePath)

	dialector := sqlite.Dialector{
		DriverName: "sqlite3",
		DSN:        databasePath,
	}

	store, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil
	}

	return store
}

func Insert(d Doujinshi) error {
	store := gormInstance()
	err := store.AutoMigrate(&Doujinshi{})
	if err != nil {
		return err
	}

	result := store.Create(&d)

	if result.Error != nil {
		return result.Error
	}

	err = Index(&d)
	if err != nil {
		return err
	}

	return nil
}

func List(req ListRequest) ([]Doujinshi, Pagination) {
	var d []Doujinshi
	var total int64
	store := gormInstance()

	offset := (req.Page - 1) * req.Limit

	// Workaround 'cause gorm's Count() seems to be tricky
	store.Table("doujinshis").Count(&total)

	if req.Seed != "0" && req.Seed != "" {
		store.Clauses(clause.OrderBy{
			Expression: clause.Expr{SQL: "SIN(id + ?)", Vars: []interface{}{req.Seed}},
		}).Limit(req.Limit).Offset(offset).Find(&d)
	} else {
		store.Model(&Doujinshi{}).Limit(req.Limit).Offset(offset).Order("created_at desc").Find(&d)
	}
	pagination := getPagination(total, req.Limit, req.Page)

	return d, pagination
}

func ListAll() []Doujinshi {
	var d []Doujinshi
	gormInstance().Table("doujinshis").Order("created_at desc").Find(&d)
	return d
}

func Get(id string) Doujinshi {
	var d Doujinshi
	gormInstance().First(&d, id)

	return d
}

func Search(text string, page int) ([]Doujinshi, Pagination) {
	var d []Doujinshi
	var total int
	var pagination Pagination
	offset := (page - 1) * LIMIT

	d, err := BleveSearch(text, offset, &total)
	if err != nil {
		return d, pagination
	}

	if d == nil {
		return nil, pagination
	}

	pagination = getPagination(int64(total), LIMIT, page)

	return d, pagination
}

func PickRandom() Doujinshi {
	var d Doujinshi
	store := gormInstance()
	store.Order("RANDOM()").First(&d)
	return d
}

func Delete(ids []string) error {
	for _, id := range ids {
		result := gormInstance().Delete(&Doujinshi{}, id)

		fmt.Printf("result.RowsAffected: %v\n", result.RowsAffected)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func Exists(id string) bool {

	d := Doujinshi{ID: id}
	gormInstance().First(&d)

	return d.Title != ""
}

func getPagination(total int64, pageSize int, page int) Pagination {
	totalPages := (int(total) + pageSize - 1) / pageSize

	startPage := max(1, page-5)
	endPage := min(totalPages, page+5)

	pages := make([]int, endPage-startPage+1)
	for i := range pages {
		pages[i] = startPage + i
	}

	pagination := Pagination{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalResults: total,
		PageSize:     pageSize,
		Pages:        pages,
	}

	return pagination
}

// Get the count of each field of the doujinshi
// this could just be a simple query but I wanted to
func (d *Doujinshi) Counters() map[string][]string {
	result := map[string][]string{}

	fields := map[string]string{
		"tags":       d.Tags,
		"artists":    d.Artists,
		"parodies":   d.Parodies,
		"characters": d.Characters,
		"groups":     d.Groups,
		"languages":  d.Languages,
		"categories": d.Categories,
	}

	for key, raw := range fields {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		result[key] = countItems("doujinshis", key, raw)
	}

	return result
}

func countItems(table, column, raw string) []string {
	values := []string{}
	items := strings.Split(raw, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue // pula item vazio
		}
		var count int64
		gormInstance().Table(table).Where(fmt.Sprintf("%s LIKE ?", column), "%"+item+"%").Count(&count)
		values = append(values, fmt.Sprintf("%d__%s", count, item))
	}
	return values
}
