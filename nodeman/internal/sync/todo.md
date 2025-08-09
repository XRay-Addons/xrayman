// questions about Sync architecture. is it ok?

// nodesync
    syncer.go            // poolsync.NodeSyncer impl
    
// poolsync
    node_client.go       // NodeClient iface
    node_storage.go      // NodeStorage iface
    node_syncer.go       // NodeSyncer iface, use NodeClient, NodeStorage

    pool_client.go       // PoolClient iface, produces NodeClients
    pool_storage.go      // PoolStorage iface

    pool_node_storage.go // NodeStorage impl based on PoolStorage
    pool_syncer.go       // implementation of pool syncer, use PoolClient, PoolStorage, NodeSyncer
