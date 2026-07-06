curl -i -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc":"2.0",
    "id":4,
    "method":"tools/call",
    "params":{
      "name":"add",
      "arguments":{
        "a":10,
        "b":20
      }
    }
  }'