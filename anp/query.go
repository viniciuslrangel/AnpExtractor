package anp

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const byStateUrl = "http://preco.anp.gov.br/include/Relatorio_Excel_Resumo_Por_Estado_Municipio.asp"
const byCityURL = "http://preco.anp.gov.br/include/Relatorio_Excel_Resumo_Por_Municipio_Posto.asp"

func ReportByState(week int, state string, fuelId int, fuelName string) ([][]string, error) {
	resp, err := http.PostForm(byStateUrl, url.Values{
		"btnSalvar":       {"Exportar"},
		"cod_semana":      {strconv.Itoa(week)},
		"COD_ESTADO":      {state},
		"COD_COMBUSTIVEL": {strconv.Itoa(fuelId)},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	node, err := goquery.NewDocumentFromReader(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	nodes := node.Find("table[border]").Find("tr").Nodes
	if len(nodes) <= 3 {
		return [][]string{}, nil
	}
	nodes = nodes[3:]
	var output [][]string
	for _, row := range nodes {
		col := []string{fuelName, cleanupName(state)}
		goquery.NewDocumentFromNode(row).ChildrenFiltered("td").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			if i > 2 {
				text = strings.Replace(text, ",", ".", 1)
			}
			col = append(col, text)
		})
		output = append(output, col)
	}
	return output, nil
}

func ReportByCity(week int, cityId int, cityName string, fuelId int, fuelName string) ([][]string, error) {
	resp, err := http.PostForm(byCityURL, url.Values{
		"btnSalvar":       {"Exportar"},
		"COD_SEMANA":      {strconv.Itoa(week)},
		"COD_MUNICIPIO":   {strconv.Itoa(cityId)},
		"COD_COMBUSTIVEL": {strconv.Itoa(fuelId)},
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	node, err := goquery.NewDocumentFromReader(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	nodes := node.Find("table[border]").First().Find("tr").Nodes
	if len(nodes) <= 2 {
		return [][]string{}, nil
	}
	nodes = nodes[2:]
	var output [][]string
	for _, row := range nodes {
		col := []string{fuelName, CityByUF[cityName], cleanupName(cityName)}
		goquery.NewDocumentFromNode(row).ChildrenFiltered("td").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			if i == 4 || i == 5 {
				text = strings.Replace(text, ",", ".", 1)
			}
			col = append(col, text)
		})
		output = append(output, col)
	}
	return output, nil
}
