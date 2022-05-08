package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Response interface {
	toString(word string) string
}

type CaiyunDictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type CaiyunDictResponse struct {
	Response
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type VolcDictRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

type VolcDictResponse struct {
	Response
	Words []struct {
		Source  int    `json:"source"`
		Text    string `json:"text"`
		PosList []struct {
			Type      int `json:"type"`
			Phonetics []struct {
				Type int    `json:"type"`
				Text string `json:"text"`
			} `json:"phonetics"`
			Explanations []struct {
				Text     string `json:"text"`
				Examples []struct {
					Type      int `json:"type"`
					Sentences []struct {
						Text      string `json:"text"`
						TransText string `json:"trans_text"`
					} `json:"sentences"`
				} `json:"examples"`
				Synonyms []interface{} `json:"synonyms"`
			} `json:"explanations"`
			Relevancys []interface{} `json:"relevancys"`
		} `json:"pos_list"`
	} `json:"words"`
	Phrases  []interface{} `json:"phrases"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}

func (c CaiyunDictResponse) toString(word string) string {
	res := ""
	res += fmt.Sprintln(word, "UK:", c.Dictionary.Prons.En, "US:", c.Dictionary.Prons.EnUs)
	for _, item := range c.Dictionary.Explanations {
		res += fmt.Sprintln(item)
	}
	return res
}

func (v VolcDictResponse) toString(word string) string {
	res := ""
	for _, item := range v.Words {
		res += fmt.Sprintln(word, "UK:", "["+item.PosList[0].Phonetics[0].Text+"]", "US:", "["+item.PosList[0].Phonetics[1].Text+"]")
		for _, exp := range item.PosList {
			if exp.Type == 13 {
				res += fmt.Sprint("adj. ")
			} else if exp.Type == 1 {
				res += fmt.Sprint("n. ")
			}
			res += fmt.Sprintln(exp.Explanations[0].Text)
		}
	}
	return res
}

func queryCaiyun(word string) Response {
	client := &http.Client{}
	request := CaiyunDictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse CaiyunDictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}

	return dictResponse
}

func queryVolc(word string) Response {
	client := &http.Client{}
	request := VolcDictRequest{
		Text:     word,
		Language: "en",
	}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	//var data = strings.NewReader(`{"text":"good","language":"en"}`)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzswVLQDc7EqHoSW8vEDcHHOL7&_signature=_02B4Z6wo00001.7kVtwAAIDAskiOivOOYq.-5FJAAJ3Yr3uJtSzulGGlFB0st7mU8IqLUULj-MoosmPjoskIvXuNJanvA9KHbHCWheYjRb68sQpFM0IKUUkKUdQes-PRgAMchQMrQrXTZubdf1", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16519790372301088; ttcid=e2395a6c2d0a4ccbb16860393f0e7d8e22; __tea_cookie_tokens_3569=%257B%2522web_id%2522%253A%25227095196007323895310%2522%252C%2522ssid%2522%253A%25225e0e0540-2693-4ef1-99ed-a871d21a879e%2522%252C%2522user_unique_id%2522%253A%25227095196007323895310%2522%252C%2522timestamp%2522%253A1651979063220%257D; isIntranet=-1; referrer_title=%E7%81%AB%E5%B1%B1%E5%BC%95%E6%93%8E-%E6%99%BA%E8%83%BD%E6%BF%80%E5%8F%91%E5%A2%9E%E9%95%BF; tt_scid=PYisN7WYzZ9MyoluLDcDIAxkkI1shfXDnpkY0T1HvH3-SsHuiIpbN-7xswgEztEAacce; i18next=translate")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=good")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}

	var dictResponse VolcDictResponse
	err = json.Unmarshal(bodyText, &dictResponse)

	if err != nil {
		log.Fatal(err)
	}

	return dictResponse
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	ch := make(chan Response)
	defer close(ch)
	go func() {
		ch <- queryCaiyun(word)
	}()
	go func() {
		ch <- queryVolc(word)
	}()
	res := <-ch
	fmt.Println(res.toString(word))
}
