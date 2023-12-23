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

package echo

import (
	"context"
	"log"

	"connectrpc.com/connect"
	echov1 "github.com/nxtaai/echo-grpc/api/echo/v1"
	"github.com/nxtaai/echo-grpc/api/echo/v1/echov1connect"
)

// EchoAPIHandler implements the Echo API service.
type EchoAPIHandler struct {
	echov1connect.UnimplementedEchoAPIHandler
}

// Echo echos the message sent by the client.
func (s *EchoAPIHandler) Echo(ctx context.Context, req *connect.Request[echov1.EchoRequest]) (*connect.Response[echov1.EchoResponse], error) {
	msg := req.Msg.Message
	log.Printf("Request headers: %v", req.Header())
	log.Printf("Received message from client: %s", msg)

	res := connect.NewResponse(&echov1.EchoResponse{
		Reply: msg,
	})
	res.Header().Set("EchoAPI-Version", "v1")
	return res, nil
}
