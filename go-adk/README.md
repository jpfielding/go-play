# Google Agent Development Kit (ADK) for Go

This directory contains examples demonstrating how to use Google's Agent Development Kit (ADK) for Go to build AI agents powered by Gemini models.

## Overview

The ADK for Go provides a comprehensive framework for building AI agents that can:
- Interact with Gemini models
- Use tools and search capabilities
- Support multiple communication protocols (CLI, Web UI, A2A)
- Implement workflow patterns (sequential, parallel)
- Handle sessions and artifacts
- Connect with remote agents

## Prerequisites

- Go 1.21 or later
- Google API Key with access to Gemini models
- Set your API key as an environment variable:
  ```bash
  export GOOGLE_API_KEY="your-api-key-here"
  ```

## Installation

The ADK is already vendored in this workspace. To use it in your own projects:

```bash
go get google.golang.org/adk
go get google.golang.org/genai
```

## Examples

### Quick Start (`quick_start/`)

A basic example showing how to create a simple LLM agent with Google Search capability.

```bash
cd quick_start
go run main.go cli
```

**Key concepts:**
- Creating a Gemini model client
- Configuring an LLM agent with tools
- Using the full launcher for CLI interaction

### Agent-to-Agent (A2A) Communication (`a2a/`)

Demonstrates how to expose an agent via A2A protocol and connect to it remotely.

```bash
cd a2a
go run main.go cli
```

**Key concepts:**
- Starting an A2A server
- Creating agent cards
- Remote agent connections
- HTTP/JSON-RPC transport

### Web Interface (`web/`)

Shows how to build agents with web UI support, multiple agents, and artifact management.

```bash
cd web
go run main.go webui
```

**Key concepts:**
- Multi-agent setup
- Web UI launcher
- Artifact storage and retrieval
- After-model callbacks
- Authentication interceptors

### Workflow Agents

#### Sequential Agent (`workflowagents/sequential/`)

Runs multiple agents in sequence, one after another.

```bash
cd workflowagents/sequential
go run main.go cli
```

#### Parallel Agent (`workflowagents/parallel/`)

Runs multiple agents concurrently.

```bash
cd workflowagents/parallel
go run main.go cli
```

**Key concepts:**
- Custom agent implementation
- Sub-agent composition
- Workflow orchestration
- Event streaming

## Core Components

### Agent

The base unit of work in ADK. Create agents using:

```go
agent, err := llmagent.New(llmagent.Config{
    Name:        "my_agent",
    Model:       model,
    Description: "What this agent does",
    Instruction: "How the agent should behave",
    Tools:       []tool.Tool{},
})
```

### Model

Connect to Gemini models:

```go
model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})
```

### Tools

Built-in tools available:
- `geminitool.GoogleSearch{}` - Web search capability
- Custom tools can be implemented

### Launcher

The launcher provides multiple interfaces:

```go
config := &launcher.Config{
    AgentLoader:     agentLoader,
    SessionService:  session.InMemoryService(),
    ArtifactService: artifact.InMemoryService(),
}

l := full.NewLauncher()
l.Execute(ctx, config, os.Args[1:])
```

Launch modes:
- `cli` - Command-line interface
- `webui` - Web user interface
- `a2a` - Agent-to-Agent server

## Project Structure

```
go-adk/
├── quick_start/       # Basic agent example
├── a2a/              # Agent-to-Agent protocol example
├── mcp/              # Model Context Protocol example
├── web/              # Web UI with multiple agents
│   ├── main.go
│   ├── llmauditor.go
│   └── image_generator.go
└── workflowagents/
    ├── sequential/   # Sequential workflow
    ├── parallel/     # Parallel workflow
    ├── loop/         # Loop workflow
    └── sequentialCode/ # Sequential with code generation
```

## Common Patterns

### Creating a Multi-Agent System

```go
agentLoader, err := agent.NewMultiLoader(
    agent1,
    agent2,
    agent3,
)
```

### Adding Callbacks

```go
llmagent.New(llmagent.Config{
    // ... other config
    AfterModelCallbacks: []llmagent.AfterModelCallback{
        func(ctx agent.CallbackContext, resp *model.LLMResponse, err error) (*model.LLMResponse, error) {
            // Process response
            return resp, err
        },
    },
})
```

### Custom Agent Implementation

```go
type MyAgent struct{}

func (a MyAgent) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
    return func(yield func(*session.Event, error) bool) {
        yield(&session.Event{
            LLMResponse: model.LLMResponse{
                Content: &genai.Content{
                    Parts: []*genai.Part{{Text: "Response"}},
                },
            },
        }, nil)
    }
}
```

## Documentation

- [ADK Documentation](https://pkg.go.dev/google.golang.org/adk)
- [Gemini API](https://pkg.go.dev/google.golang.org/genai)
- [A2A Protocol](https://github.com/a2aproject/a2a-go)

## Troubleshooting

**Issue**: "Failed to create model" error
- **Solution**: Ensure `GOOGLE_API_KEY` environment variable is set

**Issue**: Model not found
- **Solution**: Verify you have access to the specified Gemini model version

**Issue**: Tools not working
- **Solution**: Some tools may require additional API access or configuration

## License

See the Google ADK license terms for usage information.
