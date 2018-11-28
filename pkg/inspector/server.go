// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package inspector

import (
	"context"
	"fmt"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/pkg/dht"
	"storj.io/storj/pkg/node"
	"storj.io/storj/pkg/overlay"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/provider"
	"storj.io/storj/pkg/statdb"
	statsproto "storj.io/storj/pkg/statdb/proto"
)

var (
	// ServerError is a gRPC server error for Inspector
	ServerError = errs.Class("inspector server error:")
)

// Server holds references to cache and kad
type Server struct {
	dht      dht.DHT
	cache    *overlay.Cache
	statdb   *statdb.Server
	logger   *zap.Logger
	metrics  *monkit.Registry
	identity *provider.FullIdentity
}

// ---------------------
// Kad/Overlay commands:
// ---------------------

// CountNodes returns the number of nodes in the cache and in kademlia
func (srv *Server) CountNodes(ctx context.Context, req *pb.CountNodesRequest) (*pb.CountNodesResponse, error) {
	return &pb.CountNodesResponse{
		Kademlia: 0,
		Overlay:  0,
	}, nil
}

// GetBuckets returns all kademlia buckets for current kademlia instance
func (srv *Server) GetBuckets(ctx context.Context, req *pb.GetBucketsRequest) (*pb.GetBucketsResponse, error) {
	rt, err := srv.dht.GetRoutingTable(ctx)
	if err != nil {
		return &pb.GetBucketsResponse{}, ServerError.Wrap(err)
	}
	b, err := rt.GetBucketIds()
	if err != nil {
		return nil, err
	}
	bytes := b.ByteSlices()
	return &pb.GetBucketsResponse{
		Total: int64(len(b)),
		Ids:   bytes,
	}, nil
}

// GetBucket retrieves all of a given K buckets contents
func (srv *Server) GetBucket(ctx context.Context, req *pb.GetBucketRequest) (*pb.GetBucketResponse, error) {
	rt, err := srv.dht.GetRoutingTable(ctx)
	if err != nil {
		return nil, err
	}
	bucket, ok := rt.GetBucket(req.Id)
	if !ok {
		return &pb.GetBucketResponse{}, ServerError.New("GetBuckets returned non-OK response")
	}

	return &pb.GetBucketResponse{
		Id:    req.Id,
		Nodes: bucket.Nodes(),
	}, nil
}

// PingNode sends a PING RPC to the provided node ID in the Kad network.
func (srv *Server) PingNode(ctx context.Context, req *pb.PingNodeRequest) (*pb.PingNodeResponse, error) {
	rt, err := srv.dht.GetRoutingTable(ctx)
	if err != nil {
		return &pb.PingNodeResponse{}, ServerError.Wrap(err)
	}

	self := rt.Local()

	nc, err := node.NewNodeClient(srv.identity, self, srv.dht)
	if err != nil {
		return &pb.PingNodeResponse{}, ServerError.Wrap(err)
	}

	p, err := nc.Ping(ctx, pb.Node{
		Id: req.Id,
		Address: &pb.NodeAddress{
			Address: req.Address,
		},
	})

	if err != nil {
		return &pb.PingNodeResponse{}, ServerError.Wrap(err)
	}

	fmt.Printf("---- Pinged Node: %+v\n", p)

	return &pb.PingNodeResponse{}, nil
}

// ---------------------
// StatDB commands:
// ---------------------

// GetStats returns the stats for a particular node ID
func (srv *Server) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	nodeID := node.IDFromString(req.NodeId)
	getReq := &statsproto.GetRequest{
		NodeId: nodeID.Bytes(),
	}
	res, err := srv.statdb.Get(ctx, getReq)
	if err != nil {
		return nil, err
	}

	return &pb.GetStatsResponse{
		AuditRatio:  res.Stats.AuditSuccessRatio,
		UptimeRatio: res.Stats.UptimeRatio,
		AuditCount:  res.Stats.AuditCount,
	}, nil
}

// CreateStats creates a node with specified stats
func (srv *Server) CreateStats(ctx context.Context, req *pb.CreateStatsRequest) (*pb.CreateStatsResponse, error) {
	nodeID := node.IDFromString(req.NodeId)
	node := &statsproto.Node{
		NodeId: nodeID.Bytes(),
	}
	stats := &statsproto.NodeStats{
		AuditCount:         req.AuditCount,
		AuditSuccessCount:  req.AuditSuccessCount,
		UptimeCount:        req.UptimeCount,
		UptimeSuccessCount: req.UptimeSuccessCount,
	}
	createReq := &statsproto.CreateRequest{
		Node:  node,
		Stats: stats,
	}
	_, err := srv.statdb.Create(ctx, createReq)

	return &pb.CreateStatsResponse{}, err
}
