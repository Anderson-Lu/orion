grpcurl -plaintext 127.0.0.1:8001 list
grpcurl -plaintext 127.0.0.1:8081 list todo.UitTodo
grpcurl -plaintext -d '{}' 127.0.0.1:8081 todo.UitTodo.Add
