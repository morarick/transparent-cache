package sample1

import (
	"fmt"
	"time"
)

// price is a structure that represents the price as a whole price context.
// It can have several useful properties.
type price struct {
	value    float64
	cachedAt time.Time
}

// priceError is an abstraction for the price and error
type priceError struct {
	price price
	err   error
}

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]price
}

// NewTransparentCache is the implementation for PriceService interface
// It creates a new Transparent Cache based on the arguments
func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]price{},
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	if price, ok := c.prices[itemCode]; ok && time.Since(price.cachedAt) < c.maxAge {
		return price.value, nil
	}
	value, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}
	c.prices[itemCode] = price{value: value, cachedAt: time.Now()}
	return value, nil
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	ch := make(chan priceError, len(itemCodes))
	for _, itemCode := range itemCodes {
		go publishPrice(c.GetPriceFor, itemCode, ch)
	}
	return consumePrices(ch)
}

// publishPrice publish the retrieved price to the queue (channel)
func publishPrice(GetPriceFor func(string) (float64, error), itemCode string, ch chan priceError) {
	value, err := GetPriceFor(itemCode)
	ch <- priceError{price: price{value: value}, err: err}
}

// consumePrices consumes the queued prices from the given channel and returns them into an float64 slice
func consumePrices(ch chan priceError) ([]float64, error) {
	var results []float64
	for i := cap(ch); i > 0; i-- {
		priceError := <-ch
		if priceError.err != nil {
			return results, priceError.err
		}
		results = append(results, priceError.price.value)
	}
	return results, nil
}
