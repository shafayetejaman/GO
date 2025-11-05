package main

type contact struct {
	userID       string // 16 byte
	sendingLimit int32  // 4 byte
	age          int32  // 4 byte

}

type perms struct {
	permissionLevel int  // 8 byte
	canSend         bool // 1 byte
	canReceive      bool // 1 byte
	canManage       bool // 1 byte
}
