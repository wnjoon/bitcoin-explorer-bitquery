package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BitQueryConfig struct {
	URL       string
	API_KEY   string
	Variables Variables
}

type Variables struct {
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	Network    string `json:"network"`
	From       string `json:"from"`
	Till       string `json:"till"`
	Height     int    `json:"height"`
	DateFormat string `json:"dateFormat"`
}

type Timestamp struct {
	Time string `json:"time"`
}

// Common structure to store data obtained by BitQuery
type BitQueryResult struct {
	Data struct {
		Bitcoin struct {
			Blocks []struct {
				Timestamp struct {
					Time string `json:"time"`
				} `json:"timestamp"`
				Height     int     `json:"height"`
				Difficulty float64 `json:"difficulty"`
			} `json:"blocks"`
		} `json:"bitcoin"`
	} `json:"data"`
}

// Getting wanted number of block information from now
func UpdateBitcoinInfo(c BitQueryConfig) (BitQueryResult, error) {

	highestBlock, err := getHigestBlock(c)
	if err != nil {
		return BitQueryResult{}, err
	}

	highestBlockNumber := highestBlock.Data.Bitcoin.Blocks[0].Height
	c.Variables.Height = highestBlockNumber

	// Parameter for querying data from BitQuery
	jsonTypeVariable, err := structToJson(c.Variables)
	if err != nil {
		return BitQueryResult{}, err
	}

	// Query statement
	jsonData := map[string]string{
		"query": `
				query ($network: BitcoinNetwork!, $limit: Int!, $offset: Int!, $height: Int!) {
					bitcoin(network: $network) {
						blocks(options: {desc : "height", limit: $limit, offset: $offset}, height: {lteq: $height}) {
							timestamp {
								time(format: "%Y-%m-%d %H:%M:%S")
							}
						height
						difficulty
					}
				}
	  		}
		`,
		"variables": string(jsonTypeVariable),
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return BitQueryResult{}, err
	}

	bitQueryResult, err := requestBitQuery(c, jsonValue)
	if err != nil {
		return BitQueryResult{}, err
	}
	return bitQueryResult, nil

}

// Getting highest block number
func getHigestBlock(c BitQueryConfig) (BitQueryResult, error) {

	variables := Variables{
		Limit:      1, // Since function's role is just for getting most newest one, limit = 1
		Offset:     0,
		Network:    "bitcoin",
		From:       time.Now().Format("2006-01-02"),
		Till:       fmt.Sprintf("%sT23:59:59", time.Now().Format("2006-01-02")),
		DateFormat: "%Y-%m-%d",
	}

	jsonTypeVariable, err := structToJson(variables)
	if err != nil {
		return BitQueryResult{}, err
	}

	jsonData := map[string]string{
		"query": `
			query ($network: BitcoinNetwork!, $limit: Int!, $offset: Int!, $from: ISO8601DateTime, $till: ISO8601DateTime) {
				bitcoin(network: $network) {
			  		blocks(
						options: {desc: "height", limit: $limit, offset: $offset}
						date: {since: $from, till: $till}
			  		) 	{
							timestamp {
				  			time(format: "%Y-%m-%d %H:%M:%S")
						}
						height
			  		}	
				}
		  	}
        `,
		"variables": string(jsonTypeVariable),
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return BitQueryResult{}, err
	}

	bitQueryResult, err := requestBitQuery(c, jsonValue)
	if err != nil {
		return BitQueryResult{}, err
	}
	return bitQueryResult, nil
}

// Common function for calling BitQuery
func requestBitQuery(c BitQueryConfig, jsonValue []byte) (BitQueryResult, error) {
	request, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonValue))
	if err != nil {
		return BitQueryResult{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", c.API_KEY)
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)
	if err != nil {
		return BitQueryResult{}, err
	}

	defer response.Body.Close()

	readResult, err := io.ReadAll(response.Body)
	if err != nil {
		return BitQueryResult{}, err
	}

	bitQueryResult, err := jsonToStruct(readResult)
	if err != nil {
		return BitQueryResult{}, err
	}
	return bitQueryResult, nil
}

func structToJson(variables Variables) ([]byte, error) {
	jsonTypeVariable, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}
	return jsonTypeVariable, nil
}

func jsonToStruct(data []byte) (BitQueryResult, error) {
	result := BitQueryResult{}

	err := json.Unmarshal(data, &result)
	if err != nil {
		return BitQueryResult{}, err
	}
	return result, nil
}
