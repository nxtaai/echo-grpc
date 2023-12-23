// Copyright 2023 The Echo gRPC Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"connectrpc.com/connect"
	echov1 "github.com/nxtaai/echo-grpc/api/echo/v1"
	"github.com/nxtaai/echo-grpc/api/echo/v1/echov1connect"
)

func main() {
	serverAddr := flag.String("addr", "http://localhost:8080", "The server address in the format of scheme://host:port")
	flag.Parse()

	client := echov1connect.NewEchoAPIClient(http.DefaultClient, *serverAddr, connect.WithGRPC())
	res, err := client.Echo(context.Background(), connect.NewRequest(&echov1.EchoRequest{Message: "Hello gRPC!"}))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", res.Msg.Reply)
}
