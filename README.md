Anytype MCP Lite
===

Yet another Anytype MCP to return lightweight contextual note.

## Why this?

The official Anytype MCP consumes 4K tokens for MCP context, it is hard to use it in local LLMs. The Anytype store notes use "object" which cause API response contains a lot of chunks of JSON data.

My friend cannot use the official MCP in his local LLM which limited by graphics card memory. So I create this MCP to provide more control on the context and reduce the token usage.

## Spotlight

- Reduce MCP context size which minimize the support tools.
- Expand object data instead return pointer to reduce steps to get full note data.

> [!NOTE]
> Currently, only `search global` and `get object` are supported which is enough for my friend to use MCP.

## Usage

Download the binary from [release page](https://github.com/elct9620/anytype-mcp-lite/releases) and update your `mcp.json`

### Unix-like

```json
{
  "mcpServers": {
    "anytype-mcp-lite": {
      "command": "/path/to/anytype-mcp-lite",
      "env": {
        "ANYTYPE_API_KEY": "your_anytype_api_key"
      }
    }
  }
}
```

### Windows

```json
{
  "mcpServers": {
    "anytype-mcp-lite": {
      "command": "C:/path/to/anytype-mcp-lite.exe",
      "env": {
        "ANYTYPE_API_KEY": "your_anytype_api_key"
      }
    }
  }
}
```
