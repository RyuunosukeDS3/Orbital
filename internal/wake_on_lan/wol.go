package wake_on_lan

import (
	"fmt"
	"net"

	"github.com/mdlayher/wol"
)

// In internal/wake_on_lan/wake_on_lan.go

var WakeOnLan = func(mac string) error {
    macAddress, err := net.ParseMAC(mac)
    if err != nil {
        return err
    }

    c, err := wol.NewClient()
    if err != nil {
        return fmt.Errorf("failed to create WOL client: %w", err)
    }
    defer c.Close()

    if err := c.Wake("255.255.255.255", macAddress); err != nil {
        return fmt.Errorf("failed to send magic packet: %w", err)
    }

    fmt.Printf("Magic packet sent to %s\n", mac)
    return nil
}

