package tracker

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	tendermintv1beta1 "github.com/dedicatedDev/txproxy/pkg/cosmos/base/tendermint/v1beta1"
)

type TestResult struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
}
type ByHeight []TestResult

func (a ByHeight) Len() int           { return len(a) }
func (a ByHeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHeight) Less(i, j int) bool { return a[i].Height < a[j].Height }

func TrackState(client tendermintv1beta1.ServiceClient, duration int) {
	results := make(map[int64]string)

	for {
		// Check if the duration has been reached
		if len(results) == 5 {
			saveDataToLocal(results)
			results = make(map[int64]string)
		}

		resp, err := client.GetLatestBlock(context.Background(), &tendermintv1beta1.GetLatestBlockRequest{})
		if err != nil {
			log.Printf("failed to get the latest block: %v", err)
			continue
		}

		height := resp.GetBlock().Header.Height
		hash := hex.EncodeToString(resp.GetBlockId().GetHash())

		// result := TestResult{
		// 	Height: height,
		// 	Hash:   hash,
		// }

		results[height] = hash
		//results = append(results, result)
		log.Printf("Added block data: Height = %d, Hash = %s", height, hash)

		// Sleep for the block time (assuming 5 seconds)
		time.Sleep(5 * time.Second)
	}

}

func saveDataToLocal(results map[int64]string) {
	// Load existing test results from the JSON file, if any
	resultList := []TestResult{}
	for key, value := range results {
		resultList = append(resultList, TestResult{key, value})
	}
	existingResults, err := loadTestResults("test_result.json")
	if err != nil {
		existingResults = []TestResult{}
	}

	// Append the new results to the existing ones
	resultList = append(existingResults, resultList...)
	sort.Sort(ByHeight(resultList))

	// Save the combined results to the JSON file
	data, err := json.Marshal(map[string][]TestResult{
		"test_result": resultList,
	})
	if err != nil {
		log.Fatalf("failed to marshal the test results: %v", err)
	}

	if err := ioutil.WriteFile("test_result.json", data, 0644); err != nil {
		log.Fatalf("failed to write the test results to a JSON file: %v", err)
	}

	fmt.Println("Saved test results to test_result.json")
}

func loadTestResults(filename string) ([]TestResult, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var result map[string][]TestResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result["test_result"], nil
}
