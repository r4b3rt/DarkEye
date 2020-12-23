package spider

import (
	"encoding/json"
	"github.com/zsdevX/DarkEye/common"
	"net/http"
)

//SearchInformation add comment
type SearchInformation struct {
	TotalResults int `json:"total_results"`
}

//ResponseError add comment
type ResponseError struct {
	Info string `json:"info"`
}

//OrganicResult add comment
type OrganicResult struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

//Response add comment
type Response struct {
	Error             ResponseError     `json:"error"`
	SearchInformation SearchInformation `json:"search_information"`
	OrganicResults    []OrganicResult   `json:"organic_results"`
}

//Search add comment
func (s *Spider) Search() {
	s.ErrChannel <- common.LogBuild("[+]", "正在生成搜索结果", common.INFO)
	defer func() {
		s.ErrChannel <- common.LogBuild("[-]", "结束", common.INFO)
	}()

	httpClient := http.Client{}

	req, err := http.NewRequest("GET", "http://api.serpstack.com/search", nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("access_key", s.SearchAPIKey)
	q.Add("query", s.Query)
	q.Add("num", "1000")
	req.URL.RawQuery = q.Encode()

	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var apiResponse Response
	_ = json.NewDecoder(res.Body).Decode(&apiResponse)
	if apiResponse.Error.Info != "" {
		s.ErrChannel <- common.LogBuild("[-]", apiResponse.Error.Info, common.INFO)
		return
	}
	for _, result := range apiResponse.OrganicResults {
		s.ErrChannel <- common.LogBuild("[+]", result.Title+" "+result.Url, common.INFO)
	}

}
