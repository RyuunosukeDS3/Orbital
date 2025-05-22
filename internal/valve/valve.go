package valve

import (
	"fmt"
	"log"
	"time"

	"github.com/rumblefrog/go-a2s"
)

var QueryVrisingPlayersNum = func(ip string, port string) (uint8, error) {
	address := fmt.Sprintf("%s:%s", ip, port)

	client, err := a2s.NewClient(address, a2s.TimeoutOption(20*time.Second))
	if err != nil {
		return 0, fmt.Errorf("failed to create A2S client: %w", err)
	}
	
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("failed to close client: %v", err)
		}
	}()
	

	info, err := client.QueryInfo()
	if err != nil {
		return 0, fmt.Errorf("failed to query server info: %w", err)
	}

	return info.Players, nil
}
