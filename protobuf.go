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
		srvDesc := desc.FindService(strs[1])
		if srvDesc == nil {
			return nil, nil, fmt.Errorf("service name not found: %s", strs[1])
		}
		mtdDesc := srvDesc.FindMethodByName(strs[2])
		if mtdDesc == nil {
			return nil, nil, fmt.Errorf("method name not found: %s", strs[2])
		}
		return mtdDesc.GetInputType(), mtdDesc.GetOutputType(), nil
	}
	return nil, nil, fmt.Errorf("message not found: %s", path)
}
