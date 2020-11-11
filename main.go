package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/machinebox/graphql"
	"os"
)

type DataStruct struct {
	TokenDayDatas []struct {
		DailyVolumeETH      string `json:"dailyVolumeETH"`
		DailyVolumeToken    string `json:"dailyVolumeToken"`
		DailyVolumeUSD      string `json:"dailyVolumeUSD"`
		Date                int    `json:"date"`
		ID                  string `json:"id"`
		PriceUSD            string `json:"priceUSD"`
		TotalLiquidityETH   string `json:"totalLiquidityETH"`
		TotalLiquidityToken string `json:"totalLiquidityToken"`
		TotalLiquidityUSD   string `json:"totalLiquidityUSD"`
	} `json:"tokenDayDatas"`
}

func main() {
	graphqlClient := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2")
	// queries 	https://uniswap.org/docs/v2/API/queries/
	graphqlRequest := graphql.NewRequest(`
		{
		 tokenDayDatas(orderBy: date, orderDirection: asc,
		  where: {
			token: "0x6b175474e89094c44da98b954eedeac495271d0f"
		  }
		 ) {
			id
			date
			priceUSD
			totalLiquidityToken
			totalLiquidityUSD
			totalLiquidityETH
			dailyVolumeETH
			dailyVolumeToken
			dailyVolumeUSD
			 }
			}
    `)
	res := new(DataStruct)
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &res); err != nil {
		panic(err.Error())
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil{
		panic(err.Error())
	}
	fmt.Println(res.TokenDayDatas[0])
	output(string(jsonBytes))
}

func output(json string) {
	file, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%v\n", json)
	w.Flush()
}
