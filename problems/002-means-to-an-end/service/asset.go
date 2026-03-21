package service

import (
	"errors"
	"slices"

	"github.com/hellinvoid/network-programming/problems/002-means-to-an-end/dto"
)

type Asset struct {
	data map[int32]int32
	keys []int32
}

func NewAsset() *Asset {
	return &Asset{
		data: make(map[int32]int32),
		keys: make([]int32, 0),
	}
}

// Function to redirect requests
func (a *Asset) HandleRequest(req *dto.Request) (*int32, error) {
	if req.MessageType == "I" {
		return nil, a.insert(req.Timestamp, req.Price)
	}

	queryRes := a.query(req.MinTime, req.MaxTime)
	return &queryRes, nil
}

// Function to insert the data
func (a *Asset) insert(timestamp, price int32) error {

	// Check if value at timestamp already exists and return error if true
	if _, keyExists := a.data[timestamp]; keyExists {
		return errors.New("Price for Timestamp already exists")
	}

	a.data[timestamp] = price
	a.keys = append(a.keys, timestamp)

	return nil
}

// Function to query data
func (a *Asset) query(minTime, maxTime int32) int32 {

	if minTime < maxTime {
		return 0
	}

	// Sort the keys to query between given time
	slices.Sort(a.keys)

	// Binary search in the sorted keys tofind first index
	idx, _ := slices.BinarySearch(a.keys, minTime)

	var sum, count, mean int64 = 0, 0, 0
	// Calculate sum over given time range
	for ; idx < len(a.keys) && a.keys[idx] <= maxTime; idx++ {
		key := a.keys[idx]
		sum += int64(a.data[key])
		count++
	}
	
	if count == 0 {
		return 0
	}

	mean = sum / count

	return int32(mean)
}
