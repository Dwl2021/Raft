package raft

import (
	"sync"
	"time"

	"6.824/labrpc"
)

// LogEntry log条目
type LogEntry struct {
	Term    int         // leader收到日志条目时的任期
	Command interface{} // 状态机的command
}

// Raft
// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	dead      int32               // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// 所有的servers拥有的变量:
	currentTerm int        // 记录当前的任期
	votedFor    int        // 记录当前的任期把票投给了谁
	logs        []LogEntry // 日志条目数组，包含了状态机要执行的指令集，以及收到领导时的任期号

	// 所有的servers经常修改的:
	// 正常情况下commitIndex与lastApplied应该是一样的，但是如果有一个新的提交，并且还未应用的话last应该要更小些
	commitIndex int // 状态机中已知的被提交的日志条目的索引值(初始化为0，持续递增）
	lastApplied int // 最后一个被追加到状态机日志的索引值

	// leader拥有的可见变量，用来管理他的follower(leader经常修改的）
	// nextIndex与matchIndex初始化长度应该为len(peers)，Leader对于每个Follower都记录他的nextIndex和matchIndex
	// nextIndex指的是下一个的appendEntries要从哪里开始
	// matchIndex指的是已知的某follower的log与leader的log最大匹配到第几个Index,已经apply
	nextIndex  []int // 对于每一个server，需要发送给他下一个日志条目的索引值（初始化为leader日志index+1,那么范围就对标len）
	matchIndex []int // 对于每一个server，已经复制给该server的最后日志条目下标

	// 由自己追加的:
	status   Status        // 该节点是什么角色（状态）
	overtime time.Duration // 设置超时时间，200-400ms
	timer    *time.Ticker  // 每个节点中的计时器

	applyChan chan ApplyMsg // 日志都是存在这里client取（2B）
}
