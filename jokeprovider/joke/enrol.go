package joke

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Enrol() error {
	url := "https://icanhazdadjoke.com/"

	c := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/plain")
	req.Header.Set("User-Agent", "github.com/callisto13/mashed-potatoes")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
