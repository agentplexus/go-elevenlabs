// Package account provides user account, models, history, and pronunciation dictionary services.
//
// This package combines user/subscription info, available models, generation history,
// and pronunciation dictionary management.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	user, _ := client.Account().GetUser(ctx)
//	models, _ := client.Account().ListModels(ctx)
//	history, _ := client.Account().ListHistory(ctx, nil)
package account
