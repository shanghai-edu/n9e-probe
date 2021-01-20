package core

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/rpc"
	"reflect"
	"time"

	"github.com/toolkits/pkg/logger"
	"github.com/ugorji/go/codec"

	"github.com/shanghai-edu/n9e-probe/config/address"
)

func Push(metricItems []*MetricValue) error {
	var err error
	var items []*MetricValue
	now := time.Now().Unix()

	for _, item := range metricItems {
		logger.Debugf("->recv:%+v", item)
		if err = item.CheckValidity(now); err != nil {
			msg := fmt.Errorf("metric:%v err:%v", item, err)
			logger.Warning(msg)
			// 如果数据有问题，直接跳过吧，比如mymon采集的到的数据，其实只有一个有问题，剩下的都没问题
			continue
		}
		logger.Debugf("push item: %+v", item)
		items = append(items, item)
	}

	addrs := address.GetRPCAddresses("transfer")
	count := len(addrs)
	retry := 0
	for {
		for _, i := range rand.Perm(count) {
			addr := addrs[i]
			reply, err := rpcCall(addr, items)
			if err != nil {
				logger.Error(err)
				continue
			} else {
				if reply.Msg != "ok" {
					err = fmt.Errorf("some item push err: %s", reply.Msg)
					logger.Error(err)
				}
				return err
			}
		}

		time.Sleep(time.Millisecond * 500)

		retry += 1
		if retry == 3 {
			break
		}
	}

	return err
}

func rpcCall(addr string, items []*MetricValue) (TransferResp, error) {
	var reply TransferResp
	var err error

	client := rpcClients.Get(addr)
	if client == nil {
		client, err = rpcClient(addr)
		if err != nil {
			return reply, err
		}
		affected := rpcClients.Put(addr, client)
		if !affected {
			defer func() {
				// 我尝试把自己这个client塞进map失败，说明已经有一个client塞进去了，那我自己用完了就关闭
				client.Close()
			}()

		}
	}

	timeout := time.Duration(8) * time.Second
	done := make(chan error, 1)

	go func() {
		err := client.Call("Transfer.Push", items, &reply)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		logger.Warningf("rpc call timeout, transfer addr: %s\n", addr)
		rpcClients.Put(addr, nil)
		client.Close()
		return reply, fmt.Errorf("%s rpc call timeout", addr)
	case err := <-done:
		if err != nil {
			rpcClients.Del(addr)
			client.Close()
			return reply, fmt.Errorf("%s rpc call done, but fail: %v", addr, err)
		}
	}

	return reply, nil
}

func rpcClient(addr string) (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*3)
	if err != nil {
		err = fmt.Errorf("dial transfer %s fail: %v", addr, err)
		logger.Error(err)
		return nil, err
	}

	var bufConn = struct {
		io.Closer
		*bufio.Reader
		*bufio.Writer
	}{conn, bufio.NewReader(conn), bufio.NewWriter(conn)}

	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	rpcCodec := codec.MsgpackSpecRpc.ClientCodec(bufConn, &mh)
	client := rpc.NewClientWithCodec(rpcCodec)
	return client, nil
}
