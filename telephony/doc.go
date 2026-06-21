// Package telephony provides phone integration services for Twilio and SIP.
//
// This package enables conversational AI voice agents over phone calls
// using Twilio or SIP trunk integration.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	    "github.com/plexusone/elevenlabs-go/telephony"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	twiml, _ := client.Telephony().RegisterCall(ctx, &telephony.RegisterCallRequest{...})
//	numbers, _ := client.Telephony().ListPhoneNumbers(ctx)
package telephony
