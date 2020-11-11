package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"os"
)

type TokenDayData struct {
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

type GlobalStats struct {
	Data struct {
		UniswapFactory struct {
			TotalLiquidityUSD string `json:"totalLiquidityUSD"`
			TotalVolumeUSD    string `json:"totalVolumeUSD"`
			TxCount           string `json:"txCount"`
		} `json:"uniswapFactory"`
	} `json:"data"`
}

type PairData struct {
	Data struct {
		Pair struct {
			Reserve0   string `json:"reserve0"`
			Reserve1   string `json:"reserve1"`
			ReserveUSD string `json:"reserveUSD"`
			Token0     struct {
				DerivedETH string `json:"derivedETH"`
				ID         string `json:"id"`
				Name       string `json:"name"`
				Symbol     string `json:"symbol"`
			} `json:"token0"`
			Token0Price string `json:"token0Price"`
			Token1      struct {
				DerivedETH string `json:"derivedETH"`
				ID         string `json:"id"`
				Name       string `json:"name"`
				Symbol     string `json:"symbol"`
			} `json:"token1"`
			Token1Price       string `json:"token1Price"`
			TrackedReserveETH string `json:"trackedReserveETH"`
			TxCount           string `json:"txCount"`
			VolumeUSD         string `json:"volumeUSD"`
		} `json:"pair"`
	} `json:"data"`
}

func main() {
	tokenDayDataQuery()
	pairDataQuery()
	globalStatQuery()
}

func graphQuery(query string, res interface{}) error {
	graphqlClient := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2")
	graphqlRequest := graphql.NewRequest(query)
	return graphqlClient.Run(context.Background(), graphqlRequest, &res)
}

func tokenDayDataQuery() {
	query := `
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
    `
	res := new(TokenDayData)
	err := graphQuery(query, res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res.TokenDayDatas[0])
}

func globalStatQuery() {
	query := `
{
 uniswapFactory(id: "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"){
   totalVolumeUSD
   totalLiquidityUSD
   txCount
 }
}
    `
	res := new(GlobalStats)
	err := graphQuery(query, res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
}

func pairDataQuery() {
	query := `
{
 pair(id: "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11"){
     token0 {
       id
       symbol
       name
       derivedETH
     }
     token1 {
       id
       symbol
       name
       derivedETH
     }
     reserve0
     reserve1
     reserveUSD
     trackedReserveETH
     token0Price
     token1Price
     volumeUSD
     txCount
 }
}
    `
	res := new(PairData)
	err := graphQuery(query, res)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
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
