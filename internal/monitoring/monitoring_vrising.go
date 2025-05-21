package monitoring

import (
	"log"
	"time"

	"orbital/internal/argocd"
	"orbital/internal/valve"
)

var NewTicker = time.NewTicker

func PlayersOffline(vrisingIp string, vrisingPort string, appName string) {
	ticker := NewTicker(1 * time.Minute)
	defer ticker.Stop()

	offlineDuration := 0

	for {
		players, err := valve.QueryVrisingPlayersNum(vrisingIp, vrisingPort)
		if err != nil {
			log.Printf("Error querying V Rising server: %v", err)
			offlineDuration++
		} else {
			log.Printf("Players online: %d", players)
			if players == 0 {
				offlineDuration++
			} else {
				offlineDuration = 0
			}
		}

		if offlineDuration >= 5 {
			err := argocd.SetReplicaCount(appName, "0")
			if err != nil {
				log.Printf("Failed to set replica count to 0 for app %s: %v", appName, err)
			} else {
				log.Printf("Set replica count to 0 for app %s", appName)
				break
			}
		}

		<-ticker.C
	}
}
