# Hypervolt Go SDK

*This code is in no way offical nor affiliated with Hypervolt*

A Go SDK for interacting with a Hypervolt EV charger, it supports retrieval & updating of state. 

## Usage example

```go
package main

import (
	"context"
	"fmt"
	"github.com/anas-domesticus/hypervolt_go_sdk/auth"
	"github.com/anas-domesticus/hypervolt_go_sdk/rest"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc"
	"os"
)

func main() {
	// Authentication, required for both the Websocket based & REST APIs
	token, err := auth.GetToken("<USERNAME>", "<PASSWORD>")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// Getting a REST client
	c, _ := rest.NewClientWithResponses("https://api.hypervolt.co.uk")

	// This section can be omitted if you already have your charger ID
	resp, _ := c.ChargerListByOwnerWithResponse(context.Background(), token.Intercept)
	if len(*resp.JSON200.Chargers) != 1 {
		os.Exit(1)
	}
	chargers := *resp.JSON200.Chargers // Do this nicer in production please
	chargerID := int(*chargers[0].ChargerId)

	// Getting a websocket client
	rpcClient, err := rpc.NewClient(chargerID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Listens for Websocket reponses
	rpcClient.StartResponseLoop()

	// Authenticating to the Websocket client
	_, err = rpcClient.Authenticate(token.AccessToken)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Gets the current state of your charger
	snapshot, err := rpcClient.GetSnapshot()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(*snapshot.Result.LockState)
}

```

## Supported methods

| Method | Parameters | Return Type |
| ------ | ---------- | ----------- |
| Authenticate | `token string` | `*types.LoginResponse, error` |
| SetLedBrightness | `value float64` | `*types.SetLedBrightnessResponse, error` |
| SetLocked | `locked bool` | `*types.SetLockedResponse, error` |
| SetMaxCurrent | `value int` | `*types.SetMaxCurrentResponse, error` |
| GetPlugNCharge | | `*types.GetPlugNChargeResponse, error` |
| GetSchedule | | `*types.GetScheduleResponse, error` |
| SetSchedule | `sessions []types.ScheduleSession, enable bool` | `*types.SetScheduleResponse, error` |
| SetScheduleEnabled | `value bool` | `*types.SetScheduleEnabledResponse, error` |
| GetSnapshot | | `*types.GetSnapshotResponse, error` |

## Known limitations
- Login has to be with username & password