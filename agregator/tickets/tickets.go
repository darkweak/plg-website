package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const (
	artist              = "ben+plg"
	artistWithoutSpaces = "benplg"
)

type (
	provider interface {
		scrap()
		getEndpoint() string
		getName() string
	}

	storage struct {
		*sync.Mutex
		values map[string]*exportedProvider
	}

	exportedProvider struct {
		Date    string            `yaml:"date"`
		Address string            `yaml:"address"`
		Tickets map[string]string `yaml:"tickets"`
	}

	tmProvider   struct{}
	stProvider   struct{}
	fnacProvider struct{}
)

var (
	_ (provider) = (*tmProvider)(nil)
	_ (provider) = (*stProvider)(nil)
	_ (provider) = (*fnacProvider)(nil)
)

func (s *storage) store(e exportedProvider, p provider) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if s.values[e.Date] == nil {
		s.values[e.Date] = &e
	} else {
		s.values[e.Date].Tickets[p.getName()] = e.Tickets[p.getName()]
	}
}

func (t *tmProvider) scrap() {
	c := colly.NewCollector()
	c.OnHTML("#resultsListZone", func(e *colly.HTMLElement) {
		e.DOM.Children().Find(".bloc-result-content").Each(func(_ int, s *goquery.Selection) {
			url := s.Find("#urlToConcertHallLabel")
			e := exportedProvider{
				Tickets: make(map[string]string),
			}
			e.Address = url.Text()
			e.Tickets[t.getName()], _ = url.Attr("href")

			d, _ := s.Find("time").Attr("content")

			e.Date = d

			store.store(e, t)
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(t.getEndpoint())
}

func (t *tmProvider) getEndpoint() string {
	return "https://www.ticketmaster.fr/fr/resultat?ipSearch=" + artist
}

func (t *tmProvider) getName() string {
	return "ticketmaster"
}

func (st *stProvider) scrap() {
	c := colly.NewCollector()
	c.OnHTML("#search-results-wrapper", func(e *colly.HTMLElement) {
		e.DOM.Children().Find(".g-blocklist-link").Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			e := exportedProvider{
				Tickets: make(map[string]string),
			}
			e.Tickets[st.getName()] = "https://www.seetickets.com" + href
			values := strings.Split(s.Find(".g-blocklist-sub-text").Text(), "\n")
			e.Address = strings.TrimSpace(values[4])

			d := ""
			s.Find("time").Each(func(x int, se *goquery.Selection) {
				if x != 1 {
					return
				}

				d, _ = se.Attr("datetime")
			})

			date, _ := time.Parse("02 01 2006", d)

			e.Date = date.Format("2006-01-02")

			store.store(e, st)
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(st.getEndpoint())
}

func (s *stProvider) getEndpoint() string {
	return "https://www.seetickets.com/fr/search?q=" + artistWithoutSpaces
}

func (s *stProvider) getName() string {
	return "seetickets"
}

type dateFnac struct {
	Id             string `json:"productGroupId"`
	Name           string `json:"name"`
	TypeAttributes struct {
		LiveEntertainment struct {
			StartDate string `json:"startDate"`
			Location  struct {
				City string `json:"city"`
				Name string `json:"name"`
			} `json:"location"`
		} `json:"liveEntertainment"`
	} `json:"typeAttributes"`
}

func (fnac *fnacProvider) scrap() {
	time.Sleep(5 * time.Second)
	req, _ := http.NewRequest("GET", "https://public-api.eventim.com/websearch/search/api/exploration/v1/products?webId=web__eventim-fr&search_term=ben%20plg&language=FR&page=1&retail_partner=FS8&sort=Recommendation", nil)
	req.Host = "public-api.eventim.com"
	req.Header.Set("Host", "public-api.eventim.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6.1 Safari/605.1.15")
	req.Header.Set("oidc-client-id", "web__eventim-fr")
	req.Header.Set("Cookie", "_abck=62FBB8FF4E6582CAE3CEFBB0771BFD17~-1~YAAQEsARYAtoI8OFAQAACkaS6AkwLzaBmU3tr7IqD0RfAGcRMIZ2IqBatQeSgn7oBanYMcb+sWooQfiOMZGvFZdAmOEDDJttomfltxqkEDwhKD5ngzt8pUy12s3YkrTkakmd8F7hOMVFDs3XCB1FUxbru74aCx6i1LZTocgeWlSZJuuOWdu67qcaRY49GZrhYL0EemW4IpXtDQAlgJ129a2kdJFvhWw96NXxpKZrUFpZC9ODxaGnW74hZugPh3kU4OeuoSdI+2McG5AWUTJ/+1NdDMO6D1relU3de2F1HFZ4vIsp4J3XhyMr1gUzoOMRM4dOoPwqBgSVns/3N5I/MXpfMBF0Fp9ijhYdl0UHs7BGnhrJabLaCeHBGoDR+A==~-1~-1~1674647715; ak_bmsc=809A3456DA83F39C20BA43F599105EC7~000000000000000000000000000000~YAAQEsARYAxoI8OFAQAACkaS6BJ48bu70cn3m3zRrZrl+ZYQCt/5lWTkhtpyJ+hZcFe2bj2syXtGUiWlE4/pCa4f11McKxUtUuypJTdJ7uB07DBXD4HEutyETGi62cIDhG4oHJ3lsbl8E1su5c4UlgUJyK5j7YwWwd2oN7vOuCdN7wGEFzCAWFCdPLy2JmB6VtPc3krqFmNL/BcTrnHttYb/6RgG3VjxMCUAxjHI39oOfuewyTcCrd3QlySDFrwdHHQtGxnC9v9VZyX+n/CKRKYpcgudQxC14pV6izhUwuPE92ckqgHes1Cszo2qBMyMRnkiszHXqVtxgtPfzMBYnO3m7fzMNE8R3dg6vQMem57yxjFmNkg1z8qIKvI=; bm_sv=1A290704215DB7165EB0C3CBD9D3933D~YAAQxXRZaNrnvsqFAQAAMFbB6BLGPr+z6J6z4de4WnnHMh6kIuu+gb10/naywrzyc26n9Edk7MSzz/zLtYWZTbOFuk05Ce2jpopiwkYhTa9YQuIcLzk1u/F5oRVbvn1QMzyFMmBKg62baH8ygR+cWkPnk+imRnylTk6xjFPk+CFs7HpUVfumgapb5slkTYHPwXDzPhFcpzEDbUb3+REMbxtAKBypkq7yn/1Nln3FR4TuOpBWeH+9yusa2D4OmBsb9NE=~1; bm_sz=D30ED54275D27DF43F0B0AFCF9D3BA2E~YAAQEsARYA1oI8OFAQAACkaS6BKsMWeemkT/GfoocUZD2Lmom9TN1tCUEylsH+BNdoAEPDw8T9FpdKm3HHqgvTIq4+QjajItyciPeX7cXQYfRBVB9NUvtK9zHTJSB25eeIqSxwQXNwwq/kognIyiTREcK948onyWoE4sW2sonBQOzdLGI/qLY3QgkeVFCb6nqeRCLdo5Pen6qk4tr3EpjAWZ+OjlDAqgyzDopluzifg1D7sMEtXppIvf/MyMyCqqZNgB0HogjF0gK1KNmIxR7eLLO8n/IgXaRq81sq50tpxFsVkA~3162691~3355958")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res == nil || res.Body == nil {
		fmt.Println("The response is nil")
	}

	var payload struct {
		Products []dateFnac `json:"products"`
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		fmt.Printf("Impossible to decode the payload %v", err)
	}

	for _, ticket := range payload.Products {
		e := exportedProvider{
			Tickets: make(map[string]string),
		}
		e.Tickets[fnac.getName()] = "https://www.fnacspectacles.com/artist/ben-plg/" + ticket.Id
		e.Date = strings.Split(ticket.TypeAttributes.LiveEntertainment.StartDate, "T")[0]

		city := cases.Title(language.Und).String(ticket.TypeAttributes.LiveEntertainment.Location.City)
		rg := regexp.MustCompile(" " + city + ".*")
		address := cases.Title(language.Und).String(ticket.TypeAttributes.LiveEntertainment.Location.Name)

		e.Address = rg.ReplaceAllString(address, "") + ", " + city

		store.store(e, fnac)
	}
}

func (s *fnacProvider) getEndpoint() string {
	return ""
}

func (s *fnacProvider) getName() string {
	return "fnac"
}

var (
	providers = []provider{
		&tmProvider{},
		&stProvider{},
		&fnacProvider{},
	}
	store = &storage{
		Mutex:  &sync.Mutex{},
		values: make(map[string]*exportedProvider),
	}
)

func main() {
	var wg sync.WaitGroup

	for _, pr := range providers {
		wg.Add(1)
		go func(wgrp *sync.WaitGroup, p provider) {
			defer wg.Done()
			p.scrap()
		}(&wg, pr)
	}

	wg.Wait()

	vs := make([]exportedProvider, 0)
	for _, v := range store.values {
		vs = append(vs, *v)
	}
	b, _ := yaml.Marshal(vs)
	os.WriteFile("../../data/"+artistWithoutSpaces+"/dates.yaml", b, 0755)
}
