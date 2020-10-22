package gossip

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/hashicorp/memberlist"
)

//Broadcast 广播传给gossip集群
type Broadcast struct {
	Msg    []byte
	Notify chan<- struct{}
}

//Invalidates 验证数据
func (b *Broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

//Message 消息
func (b *Broadcast) Message() []byte {
	return b.Msg
}

//Finished 关闭
func (b *Broadcast) Finished() {
	if b.Notify != nil {
		close(b.Notify)
	}
}

//Delegate 客户端要必须实现的接口
//进入会员名单的八卦层。所有方法都必须是线程安全的，
//因为它们可以并且通常会被同时调用。
type Delegate interface {
	//NodeMeta用于检索有关当前节点的 meta-data
	//当广播一条活着的消息时。它的长度限制为给定的字节大小。此 meta-data在节点结构中可用。
	NodeMeta(limit int) []byte

	//NotifyMsg 在收到用户数据消息时调用
	//应注意，这种方法不会阻塞，因为
	//所以会阻塞整个UDP包接收循环。另外，字节[]byte可能在调用返回后被修改，因此如果需要，应该复制它
	NotifyMsg([]byte)

	//GetBroadcasts 当可以广播用户数据消息时，调用GetBroadcasts。
	//它可以返回要发送的缓冲区列表。每个缓冲区应假定在总字节大小上有限制的开销。
	//要发送的结果数据的总字节大小不得超过极限。应注意此方法不会阻塞，因为这样做会阻塞整个UDP包接收循环。
	GetBroadcasts(overhead, limit int) [][]byte

	//LocalState 用于TCP推/拉。发送给远程端除了成员信息。任何数据可以发送到这里。请参见MergeRemoteState。布尔值表示这是用于联接而不是push/pull
	LocalState(join bool) []byte

	//MergeRemoteState 在TCP push/pull之后被调用。是远程端的LocalState被调用后，收到数据的通知。布尔值表示这是用于联接而不是push/pull
	MergeRemoteState(buf []byte, join bool)
}

//Gossip 协议实现
type Gossip struct {
	broadcasts *memberlist.TransmitLimitedQueue
}

//NewGossip 创建gossip协议对象
func NewGossip(name string, port int, delegate Delegate, members []string) *Gossip {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	gossip := new(Gossip)
	c := memberlist.DefaultLocalConfig()
	c.Delegate = delegate
	c.BindPort = port
	c.BindAddr = "0.0.0.0"
	c.Name = name
	c.PushPullInterval = 5 * time.Second
	c.Logger = log.New(bufio.NewWriter(src), "", log.LstdFlags)
	m, err := memberlist.Create(c)
	if err != nil {
		panic(err)
	}
	if len(members) > 0 {
		_, err := m.Join(members)
		if err != nil {
			panic(err)
		}
	}
	gossip.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	//node := m.LocalNode()
	return gossip
}

//QueueBroadcast 广播
func (g *Gossip) QueueBroadcast(b *Broadcast) {
	if g.broadcasts != nil {
		g.broadcasts.QueueBroadcast(b)
	}
}
