package main

import (
	"encoding/json"
	"fmt"
	"github.com/DoG-peer/gobou/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Task is interface with gobou
type Task struct {
	config Config
}

// Config is saved
type Config struct {
	Qiita Qiita
}

// Qiita is
type Qiita struct {
	SinceTime time.Time
}

func (c Config) String() string {
	s, _ := json.MarshalIndent(c, "", "  ")
	return string(s)
}

// Load loads config
func (c *Config) Load(fname string) error {
	file, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return err
	}
	return nil
}

// Save saves config
func (c *Config) Save(fname string) error {
	return ioutil.WriteFile(fname, []byte(c.String()), os.ModePerm)
}

// RequestURL returns url
func (q *Qiita) RequestURL(n int) string {
	return fmt.Sprintf("https://qiita.com/api/v2/items?page=%d&per_page=20", n)
}

// Configure is
// load config
// make wav files
func (p *Task) Configure(configFile string, e *error) error {
	// load config file
	var conf Config
	err := conf.Load(configFile)
	if err != nil {
		*e = err
		fmt.Println(err)
		return *e
	}

	// needs mkdir

	p.config = conf
	//p.bbs.baseurl, err = GetRawURL(conf.URL)
	//p.bbs.MoveTo(conf.Res)

	// voice
	//p.voiceMng = makeVoiceManager(conf.Cache)
	// add voice
	//for _, v := range conf.Voice {
	//	p.voiceMng.add(v)
	//}
	// time.Sleep(500 * time.Millisecond)

	return err
}

// Main task
func (p *Task) Main(configFile string, m *[]gobou.Message) error {
	n := 1
	/*
		q := Qiita{
			SinceTime: time.Now().Add(-10000 * time.Second),
		}
	*/
	q := p.config.Qiita
	t := q.SinceTime
	mes := []gobou.Message{}
L:
	for i := 1; i <= n; i++ {
		client, err := http.Get(q.RequestURL(i))
		defer client.Body.Close()

		arr, err := ioutil.ReadAll(client.Body)
		if err != nil {
			return err
		}
		var result QResult
		json.Unmarshal(arr, &result)
		for _, item := range result {
			title := item.Title
			at := item.CreatedAt
			s := fmt.Sprintf("---\ntitle: %s\ncreated_at: %s\n---\n", title, at)
			if at.After(t) {
				user := item.User
				// Notify(fmt.Sprintf("qiita by %s", user.ID), title)
				s = fmt.Sprintf("%sqiita by %s\n", s, user.ID)
			} else {
				break L
			}
			mes = append(mes, gobou.Print(s), gobou.Say(fmt.Sprint(i)))
		}
	}
	*m = mes
	//data, err := p.bbs.Read()

	return nil
}

// SaveData is loaded by gobou
func (p *Task) SaveData(configFile string, e *error) error {
	return nil
}

// SaveConfig is loaded by gobou
func (p *Task) SaveConfig(configFile string, e *error) error {
	p.config.Save(configFile)
	return nil
}

// Interval is loaded by gobou
func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 60 * time.Minute
	return nil
}

// QResult means result returned from Qiita
type QResult []QItem

// QItem means item from Qiita
type QItem struct {
	Title     string    `json:"title"`
	User      QUser     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/*
	RenderedBody interface{} `json:rendered_body`
	Tags         interface{} `json:tags`
	Body         interface{} `json:body`
	ID           interface{} `json:id`
	Private      interface{} `json:private`
	URL          interface{} `json:url`
*/

// QUser means Qiita user
type QUser struct {
	Name string `json:name`
	ID   string `json:id`
}

func main() {
	gobou.Register(&Task{})
}

/*
rendered_body
coediting
tags
title
updated_at
user
body
created_at
id
private
url
*/
/*arr := []byte(`[{
	"title": "ttttt",
	"user": {"name":"uuuuu","id":"ididid"},
	"created_at": "2014-08-25T00:00:00+09:00",
	"updated_at": "2014-08-25T00:00:00+09:00"
}]`)*/
