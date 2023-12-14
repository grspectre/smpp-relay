package main

import (
	"strings"
)

type MessageList struct {
	messages map[string]MessageBag
}

func (ml *MessageList) get(id string) MessageBag {
	if val, ok := ml.messages[id]; ok {
		return val
	}
	return MessageBag{
		message: "",
		ok:      false,
		parts:   make([]string, 0),
	}
}

func (ml *MessageList) set(id string, message MessageBag) {
	if ml.messages == nil {
		ml.messages = make(map[string]MessageBag)
	}
	ml.messages[id] = message
}

type MessageBag struct {
	ok       bool
	message  string
	parts    []string
	partsIdx []int
}

func (mb *MessageBag) add(count int, part int, message string) {
	if len(mb.parts) == 0 {
		empty := make([]string, count)
		mb.parts = append(mb.parts, empty...)
	}
	mb.partsIdx = append(mb.partsIdx, part)
	mb.parts[part-1] = message
	if len(mb.partsIdx) == count {
		mb.message = strings.Join(mb.parts, "")
	}
}

var globalMessageList MessageList

func messageOrEmpty(idx string, message string) string {
	bag := globalMessageList.get(idx)

	messageBytes := []byte(message)
	if messageBytes[0] < 32 {
		info := messageBytes[0:6]
		text := messageBytes[6:]
		count, part := info[4], info[5]
		bag.add(int(count), int(part), string(text))
	} else {
		bag.add(1, 1, string(messageBytes))
	}
	globalMessageList.set(idx, bag)

	return bag.message
}
