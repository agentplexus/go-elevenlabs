// Package agents provides ElevenAgents (Conversational AI) services.
//
// This package handles all aspects of ElevenLabs' Conversational AI platform:
//
//   - Agent management (create, update, delete, duplicate)
//   - Agent branches and deployments for version control
//   - Conversation history, transcripts, and analysis
//   - Knowledge base document management
//   - Agent testing with test folders and response tests
//   - Conversation simulation (streaming and non-streaming)
//   - Batch calling for outbound campaigns
//   - Analytics and live call counts
//
// # Basic Usage
//
//	client, _ := elevenlabs.NewClient()
//	agentsSvc := client.Agents()
//
//	// List agents
//	agents, _ := agentsSvc.List(ctx, nil)
//
//	// Get conversation history
//	conversations, _ := agentsSvc.ListConversations(ctx, &agents.ListConversationsOptions{
//	    AgentID: "your-agent-id",
//	})
//
// # Knowledge Base
//
// The knowledge base allows you to add documents that agents can reference:
//
//	// Create a text document
//	doc, _ := agentsSvc.CreateTextDocument(ctx, &agents.CreateTextDocumentRequest{
//	    Name:    "FAQ",
//	    Content: "Frequently asked questions...",
//	})
//
//	// Create a URL document (web page crawl)
//	doc, _ := agentsSvc.CreateURLDocument(ctx, &agents.CreateURLDocumentRequest{
//	    Name: "Documentation",
//	    URL:  "https://docs.example.com",
//	})
//
// # Conversation Simulation
//
// Test your agent with simulated conversations:
//
//	// Non-streaming simulation
//	messages, _ := agentsSvc.SimulateConversation(ctx, "agent-id", &agents.SimulateConversationOptions{
//	    SimulationPersona: "A curious customer asking about products",
//	    MaxTurns:          5,
//	})
//
//	// Streaming simulation via WebSocket
//	conn, _ := agentsSvc.SimulateConversationStream(ctx, "agent-id", nil)
//	defer conn.Close()
//	for msg := range conn.Messages() {
//	    fmt.Printf("[%s]: %s\n", msg.Role, msg.Content)
//	}
//
// # Batch Calling
//
// Create outbound calling campaigns:
//
//	batch, _ := agentsSvc.CreateBatchCall(ctx, &agents.CreateBatchCallRequest{
//	    AgentID: "your-agent-id",
//	    Name:    "Sales Campaign",
//	    Recipients: []agents.BatchCallRecipient{
//	        {PhoneNumber: "+1234567890"},
//	        {PhoneNumber: "+0987654321"},
//	    },
//	})
package agents
