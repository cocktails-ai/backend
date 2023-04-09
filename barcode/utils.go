package barcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProductInfo struct {
	Status        int    `json:"status"`
	StatusVerbose string `json:"status_verbose"`
	Product       struct {
		ProductName string `json:"product_name"`
	} `json:"product"`
}

func FindBarcode(barcode string) (string, error) {
	url := fmt.Sprintf("https://world.openfoodfacts.org/api/v0/product/%s.json", barcode)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var productInfo ProductInfo
	if err := json.Unmarshal(body, &productInfo); err != nil {
		return "", err
	}

	if productInfo.Status != 1 {
		return "", errors.New(productInfo.StatusVerbose)
	}

	return productInfo.Product.ProductName, nil
}
