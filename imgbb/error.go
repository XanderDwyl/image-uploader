package imgbb

import "fmt"

// ErrResp ...
type ErrResp struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_txt"`
	ErrMsg     error  `json:"err_msg"`
}

// Error ...
func (e ErrResp) Error() string {
	return fmt.Sprintf("%d %s: %v", e.StatusCode, e.StatusText, e.ErrMsg)
}

// Is ...
func (e ErrResp) Is(target error) bool {
	if err, ok := target.(ErrResp); !ok {
		return false
	} else {
		return err.StatusCode == e.StatusCode && err.StatusText == e.StatusText
	}
}
