package handler

import "net/http"
import "io"
import "fmt"
import "strconv"
import "time"

func (h Handler) delay(w http.ResponseWriter, r *http.Request) {
	strDelay := r.URL.Query().Get("delay")

	intDelay, err := strconv.Atoi(strDelay)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("Invalid delay: %v", strDelay))
		return
	}

	duration := time.Duration(intDelay) * time.Millisecond

	time.Sleep(duration)

	io.WriteString(w, fmt.Sprintf("Delay of %s", duration))
}
