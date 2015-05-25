/******************************************************************************
 * reply_codes.go
 * Description: reply_codes contains constants that represent status messages
 * to send back to the client so that the client can act based on those codes
 ******************************************************************************/

package main

const (
	// Name          // Code
	ChannelJoinSuccess = 050
	ChannelJoinFailure = 051
)

var FormattedStatusCodeMessages = make(map[int]string)

func InitStatusCodeMessages() {
	StatusCodeMessages[ChannelJoinSuccess] = "%s join successfully"
	StatusCodeMessages[ChannelJoinFailure] = "%s join failed"
}
