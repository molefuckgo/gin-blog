package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"` // gorm:"index"：用于声明这个字段为索引，使用自动迁移功能则会有影响？？？？	外键
	Tag   Tag `json:"tag"`                 // 嵌套struct, 利用TagID与Tag模型相互关联，执行查询时，可以达到Article、Tag关联查询的功能

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (article Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum * pageSize).Limit(pageSize).Find(&article)
	return
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	fmt.Println(article) // {{1 1570591493 0} 19 {{0 0 0}    0} test1 test-desc test-content   0}
	db.Model(&article).Related(&article.Tag)
	fmt.Println(article) // {{1 1570591493 0} 19 {{19 1569748326 0} 标签name 郭威  1} test1 test-desc test-content   0}

	return
}

func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Update(data)

	return true
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:   data["tag_id"].(int),
		Title:   data["title"].(string),
		Desc:    data["desc"].(string),
		Content: data["content"].(string),
		State:   data["state"].(int),
	})

	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)
	if article.ID > 0 {
		return true
	}
	return false
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
