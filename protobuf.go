//    Copyright 2018 <xu shuai <dyxushuai@gmail.com>>
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package grpcdump

import (
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

type protoFileDescs struct {
	descs []*desc.FileDescriptor
}

func protoParse(importPaths []string, protoPaths ...string) (*protoFileDescs, error) {
	p := &protoparse.Parser{
		ImportPaths: importPaths,
	}
	descs, err := p.ParseFiles(protoPaths...)
	if err != nil {
		return nil, err
	}
	return &protoFileDescs{
		descs: descs,
	}, nil
}

func (p *protoFileDescs) pretty(data []byte, msg *desc.MessageDescriptor) (string, error) {
	dmsg := dynamic.NewMessage(msg)
	err := dmsg.Unmarshal(data)
	if err != nil {
		return "", err
	}
	result, err := dmsg.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (p *protoFileDescs) findMehodSignature(path string) (*desc.MessageDescriptor, *desc.MessageDescriptor, error) {
	strs := strings.Split(path, "/")
	if len(strs) != 3 {
		return nil, nil, fmt.Errorf("error path format: %s", path)
	}
	for _, desc := range p.descs {
		// FindService parameter needs package.servcie
		// e.g. helloworld.Greeter
		srvDesc := desc.FindService(strs[1])
		if srvDesc == nil {
			// maybe in other files
			continue
		}
		mtdDesc := srvDesc.FindMethodByName(strs[2])
		if mtdDesc == nil {
			return nil, nil, fmt.Errorf("method: %s not found in file: %s", strs[2], desc.GetName())
		}
		return mtdDesc.GetInputType(), mtdDesc.GetOutputType(), nil
	}
	return nil, nil, fmt.Errorf("grpc path: %s not found", path)
}
