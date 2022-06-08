package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/royalcat/gotoschema"
)

func main() {
	gen := gotoschema.NewDocGenerator()
	err := gen.AddModel(Info{}, Info2{})
	if err != nil {
		log.Fatal(err)
	}
	doc, err := gen.EncodeYaml()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("doc.yml", []byte(doc), 0644)
	if err != nil {
		panic(err)
	}
}

type User struct {
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email,omitempty"`
	Document Document `json:"document"`
}

type Type string

type Info struct {
	Type     Type      `json:"type,omitempty"`
	Id       int64     `json:"id,omitempty"`
	IsActive bool      `json:"is_active,omitempty"`
	Date     time.Time `json:"sign_date,omitempty"`
	Users    []User    `json:"users,omitempty"`
	Document Document  `json:"document"`
}

type Info2 struct {
	Id       int64     `json:"id,omitempty"`
	Type     Type      `json:"type,omitempty"`
	IsActive bool      `json:"is_active,omitempty"`
	Date     time.Time `json:"sign_date,omitempty"`
	Users    []User    `json:"users,omitempty"`
	Document Document  `json:"document"`
}

type Document struct {
	Type    string  `json:"type,omitempty"`
	Series  string  `json:"series,omitempty"`
	Number  string  `json:"number,omitempty"`
	Issued  string  `json:"issued,omitempty"`
	Comment Comment `json:"comment,omitempty"`
}

type Comment struct {
	Text    string    `json:"text"`
	Options []Options `json:"options"`
}

type Options struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}
