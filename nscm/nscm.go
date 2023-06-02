package nscm

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/KirillMerz/NSCMTelegramBot/models"
)

const NSCM_URL = "http://nscm.ru/giaresult/tablresult.php"

func GetResults(user models.User) (string, error) {
	page, err := fetchResultsFromNSCM(user)
	if err != nil {
		return "", err
	}

	results, err := parseResults(page)
	if err != nil {
		return "", err
	}

	return genMessageWithResults(results), nil
}

func fetchResultsFromNSCM(user models.User) (*html.Node, error) {
	values := url.Values{
		"Lastname":   {user.Lastname},
		"Name":       {user.Name},
		"SecondName": {user.SecondName},
		"DocNumber":  {user.DocNumber},
	}

	resp, err := http.PostForm(NSCM_URL, values)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching results error: response status: %s", resp.Status)
	}

	return html.Parse(resp.Body)
}

func parseResults(page *html.Node) (models.Results, error) {
	doc := goquery.NewDocumentFromNode(page)

	var results models.Results
	results.LastUpdated = time.Now().Unix()

	doc.Find("tbody").Find("tr").Each(func(_ int, tr *goquery.Selection) {
		var buf []string

		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			buf = append(buf, td.Text())
		})

		points, _ := strconv.Atoi(buf[2])

		result := models.Result{
			Subject: buf[0],
			Points:  points,
			Mark:    buf[3],
		}

		results.List = append(results.List, result)
	})

	return results, nil
}

func genMessageWithResults(res models.Results) string {
	var message string

	for _, r := range res.List {
		message += fmt.Sprintf("%s: %s (%s)\n", r.Subject, r.Mark, genPointsDescription(r.Points))
	}

	return message
}

func genPointsDescription(points int) string {
	// example: genPointsDescription(20) -> "20 баллов"
	//          genPointsDescription(21) -> "21 балл"
	//          genPointsDescription(22) -> "22 балла"

	lastDigit := points % 10

	var pointsWord string

	switch {
	case (points > 4 && points < 21) || (lastDigit == 0) || (lastDigit > 4):
		pointsWord = "баллов"
	case lastDigit == 1:
		pointsWord = "балл"
	default:
		pointsWord = "балла"
	}

	return fmt.Sprintf("%d %s", points, pointsWord)
}
