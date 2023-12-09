package errors

//
//import (
//	"fmt"
//	"net/http"
//)
//
//type HTTPError struct {
//	StatusCode int
//	URL        string
//}
//
//func (he *HTTPError) Error() string {
//	return fmt.Sprintf("http status code = %d, url = %s", he.StatusCode, he.URL)
//}
//
//func ReadContents(url string) ([]byte, error) {
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//
//	if err.StatusCode != http.StatusOK {
//		return nil, &HTTPError{
//			StatusCode: resp.StatusCode,
//			URL:        url,
//		}
//	}
//}
