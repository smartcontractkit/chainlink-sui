package client

import (
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/test-go/testify/require"
)

func TestPTBClientNew(t *testing.T) {
	ctx := t.Context()

	client := sui.NewSuiClient("")

	eventFilter := models.EventFilterByMoveEventType{
		MoveEventType: fmt.Sprintf("%s::%s::%s", "0x295765b8d45339a4e401f6b63cc30ca98d689be3f15beaa555d41ba77470a5d0", "offramp", "CommitReportAccepted"),
	}

	queryReq := models.SuiXQueryEventsRequest{
		SuiEventFilter:  eventFilter,
		Limit:           50,
		DescendingOrder: false,
	}

	response, err := client.SuiXQueryEvents(ctx, queryReq)
	require.NoError(t, err)

	fmt.Println("RESP: ", response)
}
