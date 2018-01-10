package x

import (
	"context"
	"fmt"
	"time"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-x/rpc"
	"google.golang.org/grpc"
)

// Dial the target identity.MultiAddress using a background context.Context.
// Returns a grpc.ClientConn, or an error. The grpc.ClientConn must be closed
// before it exists scope.
func Dial(target identity.MultiAddress) (*grpc.ClientConn, error) {
	host, err := target.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := target.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Send an order fragment to the target identity.MultiAddress. Returns nil, or an
// error.
func SendOrderFragment(target identity.MultiAddress, fragment *rpc.OrderFragment) (*identity.MultiAddress, error) {
	// Create the client.
	client, conn, err := NewNodeClient(target)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send the order fragment.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	response, err := client.SendOrderFragment(ctx, fragment, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}
	multi, err := identity.NewMultiAddress(response.Multi)
	if err != nil {
		return nil, err
	}
	return &multi, nil
}

// RPCPing sends an Ping RPC request to the target using a new grpc.ClientConn
// and a new rpc.NodeClient. It returns the result of the RPC call, or an
// error.
func (node *Node) RPCPing(target identity.MultiAddress) (identity.MultiAddress, error) {

	// Connect to the target.
	conn, err := Dial(target)
	if err != nil {
		return identity.MultiAddress{}, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	// Create a timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Call the Ping RPC on the target.
	multi, err := client.Ping(ctx, &rpc.MultiAddress{Multi: node.MultiAddress.String()}, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddress{}, err
	}

	return identity.NewMultiAddressFromString(multi.Multi)
}

// RPCPeers sends an Peers RPC request to the target using a new
// grpc.ClientConn and a new rpc.NodeClient. It returns the result of the RPC
// call, or an error.
func RPCPeers(target identity.MultiAddress) (identity.MultiAddresses, error) {

	// Connect to the target.
	conn, err := Dial(target)
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	defer conn.Close()
	client := rpc.NewNodeClient(conn)

	// Create a timeout context.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Call the Peers RPC on the target.
	multis, err := client.Peers(ctx, &rpc.Nothing{}, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}

	peers := make(identity.MultiAddresses, 0, len(multis.Multis))
	for i, multi := range multis.Multis {
		peer, err := identity.NewMultiAddressFromString(multi.Multi)
		if err != nil {
			return peers, err
		}
		peers = append(peers, peer)
	}
	return peers, nil
}
