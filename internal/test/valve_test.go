package test

import (
	"errors"
	"testing"

	"orbital/internal/valve"

	"github.com/stretchr/testify/assert"
)

const vrisingMockUrl string = "127.0.0.1"

func TestQueryVrisingPlayersNumSuccess(t *testing.T) {
	original := valve.QueryVrisingPlayersNum
	defer func() { valve.QueryVrisingPlayersNum = original }()

	valve.QueryVrisingPlayersNum = func(ip, port string) (uint8, error) {
		assert.Equal(t, vrisingMockUrl, ip)
		assert.Equal(t, "9876", port)
		return 5, nil
	}

	players, err := valve.QueryVrisingPlayersNum(vrisingMockUrl, "9876")
	assert.NoError(t, err)
	assert.Equal(t, uint8(5), players)
}

func TestQueryVrisingPlayersNumError(t *testing.T) {
	original := valve.QueryVrisingPlayersNum
	defer func() { valve.QueryVrisingPlayersNum = original }()

	valve.QueryVrisingPlayersNum = func(ip, port string) (uint8, error) {
		return 0, errors.New("mocked error")
	}

	players, err := valve.QueryVrisingPlayersNum(vrisingMockUrl, "9876")
	assert.Error(t, err)
	assert.Equal(t, uint8(0), players)
	assert.Contains(t, err.Error(), "mocked error")
}
