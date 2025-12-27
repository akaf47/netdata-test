package discover

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func TestNewVsphereBuilder(t *testing.T) {
	t.Run("should create a new vsphere builder", func(t *testing.T) {
		ctx := context.Background()
		builder := NewVsphereBuilder(ctx)
		
		assert.NotNil(t, builder)
		assert.NotNil(t, builder.Context)
		assert.Equal(t, ctx, builder.Context)
	})

	t.Run("should create builder with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		
		builder := NewVsphereBuilder(ctx)
		assert.NotNil(t, builder)
		assert.Equal(t, ctx, builder.Context)
	})

	t.Run("should create builder with nil context", func(t *testing.T) {
		builder := NewVsphereBuilder(nil)
		assert.NotNil(t, builder)
	})
}

func TestBuildVMs(t *testing.T) {
	t.Run("should return empty slice when vms is nil", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildVMs(nil)
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return empty slice when vms is empty", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildVMs([]*object.VirtualMachine{})
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should convert single virtual machine to target", func(t *testing.T) {
		builder := &VsphereBuilder{}
		vm := createMockVirtualMachine("vm-1", "TestVM")
		
		result := builder.BuildVMs([]*object.VirtualMachine{vm})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "vm-1", result[0].ID)
		assert.Equal(t, "TestVM", result[0].Name)
	})

	t.Run("should convert multiple virtual machines to targets", func(t *testing.T) {
		builder := &VsphereBuilder{}
		vms := []*object.VirtualMachine{
			createMockVirtualMachine("vm-1", "TestVM1"),
			createMockVirtualMachine("vm-2", "TestVM2"),
			createMockVirtualMachine("vm-3", "TestVM3"),
		}
		
		result := builder.BuildVMs(vms)
		
		assert.Len(t, result, 3)
		assert.Equal(t, "vm-1", result[0].ID)
		assert.Equal(t, "vm-2", result[1].ID)
		assert.Equal(t, "vm-3", result[2].ID)
	})

	t.Run("should handle vm with special characters in name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		vm := createMockVirtualMachine("vm-special", "Test-VM_001")
		
		result := builder.BuildVMs([]*object.VirtualMachine{vm})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "Test-VM_001", result[0].Name)
	})

	t.Run("should handle vm with unicode characters in name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		vm := createMockVirtualMachine("vm-unicode", "Test🚀VM")
		
		result := builder.BuildVMs([]*object.VirtualMachine{vm})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "Test🚀VM", result[0].Name)
	})
}

func TestBuildHosts(t *testing.T) {
	t.Run("should return empty slice when hosts is nil", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildHosts(nil)
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return empty slice when hosts is empty", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildHosts([]*object.HostSystem{})
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should convert single host to target", func(t *testing.T) {
		builder := &VsphereBuilder{}
		host := createMockHostSystem("host-1", "TestHost")
		
		result := builder.BuildHosts([]*object.HostSystem{host})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "host-1", result[0].ID)
		assert.Equal(t, "TestHost", result[0].Name)
	})

	t.Run("should convert multiple hosts to targets", func(t *testing.T) {
		builder := &VsphereBuilder{}
		hosts := []*object.HostSystem{
			createMockHostSystem("host-1", "TestHost1"),
			createMockHostSystem("host-2", "TestHost2"),
			createMockHostSystem("host-3", "TestHost3"),
		}
		
		result := builder.BuildHosts(hosts)
		
		assert.Len(t, result, 3)
		assert.Equal(t, "host-1", result[0].ID)
		assert.Equal(t, "host-2", result[1].ID)
		assert.Equal(t, "host-3", result[2].ID)
	})

	t.Run("should handle host with special characters in name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		host := createMockHostSystem("host-special", "Test-Host.001")
		
		result := builder.BuildHosts([]*object.HostSystem{host})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "Test-Host.001", result[0].Name)
	})

	t.Run("should handle host with long name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		longName := "very-long-hostname-that-exceeds-normal-limits-1234567890-abcdefghijklmnop"
		host := createMockHostSystem("host-long", longName)
		
		result := builder.BuildHosts([]*object.HostSystem{host})
		
		assert.Len(t, result, 1)
		assert.Equal(t, longName, result[0].Name)
	})
}

func TestBuildDatacenters(t *testing.T) {
	t.Run("should return empty slice when datacenters is nil", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildDatacenters(nil)
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return empty slice when datacenters is empty", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildDatacenters([]*object.Datacenter{})
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should convert single datacenter to target", func(t *testing.T) {
		builder := &VsphereBuilder{}
		dc := createMockDatacenter("dc-1", "TestDC")
		
		result := builder.BuildDatacenters([]*object.Datacenter{dc})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "dc-1", result[0].ID)
		assert.Equal(t, "TestDC", result[0].Name)
	})

	t.Run("should convert multiple datacenters to targets", func(t *testing.T) {
		builder := &VsphereBuilder{}
		dcs := []*object.Datacenter{
			createMockDatacenter("dc-1", "TestDC1"),
			createMockDatacenter("dc-2", "TestDC2"),
			createMockDatacenter("dc-3", "TestDC3"),
		}
		
		result := builder.BuildDatacenters(dcs)
		
		assert.Len(t, result, 3)
		assert.Equal(t, "dc-1", result[0].ID)
		assert.Equal(t, "dc-2", result[1].ID)
		assert.Equal(t, "dc-3", result[2].ID)
	})

	t.Run("should handle datacenter with special characters in name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		dc := createMockDatacenter("dc-special", "Test_DC-Prod")
		
		result := builder.BuildDatacenters([]*object.Datacenter{dc})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "Test_DC-Prod", result[0].Name)
	})

	t.Run("should handle datacenter with empty name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		dc := createMockDatacenter("dc-empty", "")
		
		result := builder.BuildDatacenters([]*object.Datacenter{dc})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "", result[0].Name)
	})
}

func TestBuildClusters(t *testing.T) {
	t.Run("should return empty slice when clusters is nil", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildClusters(nil)
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return empty slice when clusters is empty", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildClusters([]*object.ClusterComputeResource{})
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should convert single cluster to target", func(t *testing.T) {
		builder := &VsphereBuilder{}
		cluster := createMockCluster("cluster-1", "TestCluster")
		
		result := builder.BuildClusters([]*object.ClusterComputeResource{cluster})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "cluster-1", result[0].ID)
		assert.Equal(t, "TestCluster", result[0].Name)
	})

	t.Run("should convert multiple clusters to targets", func(t *testing.T) {
		builder := &VsphereBuilder{}
		clusters := []*object.ClusterComputeResource{
			createMockCluster("cluster-1", "TestCluster1"),
			createMockCluster("cluster-2", "TestCluster2"),
			createMockCluster("cluster-3", "TestCluster3"),
		}
		
		result := builder.BuildClusters(clusters)
		
		assert.Len(t, result, 3)
		assert.Equal(t, "cluster-1", result[0].ID)
		assert.Equal(t, "cluster-2", result[1].ID)
		assert.Equal(t, "cluster-3", result[2].ID)
	})

	t.Run("should handle cluster with dash in name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		cluster := createMockCluster("cluster-dash", "prod-cluster-01")
		
		result := builder.BuildClusters([]*object.ClusterComputeResource{cluster})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "prod-cluster-01", result[0].Name)
	})

	t.Run("should handle cluster with numeric name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		cluster := createMockCluster("cluster-123", "123")
		
		result := builder.BuildClusters([]*object.ClusterComputeResource{cluster})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "123", result[0].Name)
	})
}

func TestBuildResourcePools(t *testing.T) {
	t.Run("should return empty slice when pools is nil", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildResourcePools(nil)
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return empty slice when pools is empty", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildResourcePools([]*object.ResourcePool{})
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should convert single resource pool to target", func(t *testing.T) {
		builder := &VsphereBuilder{}
		pool := createMockResourcePool("pool-1", "TestPool")
		
		result := builder.BuildResourcePools([]*object.ResourcePool{pool})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "pool-1", result[0].ID)
		assert.Equal(t, "TestPool", result[0].Name)
	})

	t.Run("should convert multiple resource pools to targets", func(t *testing.T) {
		builder := &VsphereBuilder{}
		pools := []*object.ResourcePool{
			createMockResourcePool("pool-1", "TestPool1"),
			createMockResourcePool("pool-2", "TestPool2"),
			createMockResourcePool("pool-3", "TestPool3"),
		}
		
		result := builder.BuildResourcePools(pools)
		
		assert.Len(t, result, 3)
		assert.Equal(t, "pool-1", result[0].ID)
		assert.Equal(t, "pool-2", result[1].ID)
		assert.Equal(t, "pool-3", result[2].ID)
	})

	t.Run("should handle resource pool with slash in name", func(t *testing.T) {
		builder := &VsphereBuilder{}
		pool := createMockResourcePool("pool-slash", "parent/child")
		
		result := builder.BuildResourcePools([]*object.ResourcePool{pool})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "parent/child", result[0].Name)
	})

	t.Run("should handle resource pool with deep hierarchy", func(t *testing.T) {
		builder := &VsphereBuilder{}
		pool := createMockResourcePool("pool-deep", "a/b/c/d/e")
		
		result := builder.BuildResourcePools([]*object.ResourcePool{pool})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "a/b/c/d/e", result[0].Name)
	})
}

func TestBuildFolders(t *testing.T) {
	t.Run("should return empty slice when folders is nil", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildFolders(nil)
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return empty slice when folders is empty", func(t *testing.T) {
		builder := &VsphereBuilder{}
		result := builder.BuildFolders([]*object.Folder{})
		
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("should convert single folder to target", func(t *testing.T) {
		builder := &VsphereBuilder{}
		folder := createMockFolder("folder-1", "TestFolder")
		
		result := builder.BuildFolders([]*object.Folder{folder})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "folder-1", result[0].ID)
		assert.Equal(t, "TestFolder", result[0].Name)
	})

	t.Run("should convert multiple folders to targets", func(t *testing.T) {
		builder := &VsphereBuilder{}
		folders := []*object.Folder{
			createMockFolder("folder-1", "TestFolder1"),
			createMockFolder("folder-2", "TestFolder2"),
			createMockFolder("folder-3", "TestFolder3"),
		}
		
		result := builder.BuildFolders(folders)
		
		assert.Len(t, result, 3)
		assert.Equal(t, "folder-1", result[0].ID)
		assert.Equal(t, "folder-2", result[1].ID)
		assert.Equal(t, "folder-3", result[2].ID)
	})

	t.Run("should handle folder with numeric suffix", func(t *testing.T) {
		builder := &VsphereBuilder{}
		folder := createMockFolder("folder-num", "VMs_01")
		
		result := builder.BuildFolders([]*object.Folder{folder})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "VMs_01", result[0].Name)
	})

	t.Run("should handle folder with case variations", func(t *testing.T) {
		builder := &VsphereBuilder{}
		folder := createMockFolder("folder-case", "FolderName")
		
		result := builder.BuildFolders([]*object.Folder{folder})
		
		assert.Len(t, result, 1)
		assert.Equal(t, "FolderName", result[0].Name)
	})
}

// Helper functions to create mock objects

func createMockVirtualMachine(id, name string) *object.VirtualMachine {
	ref := types.ManagedObjectReference{
		Type:  "VirtualMachine",
		Value: id,
	}
	return &object.VirtualMachine{
		Common: object.Common{
			InventoryPath: name,
			Reference:     ref,
		},
	}
}

func createMockHostSystem(id, name string) *object.HostSystem {
	ref := types.ManagedObjectReference{
		Type:  "HostSystem",
		Value: id,
	}
	return &object.HostSystem{
		Common: object.Common{
			InventoryPath: name,
			Reference:     ref,
		},
	}
}

func createMockDatacenter(id, name string) *object.Datacenter {
	ref := types.ManagedObjectReference{
		Type:  "Datacenter",
		Value: id,
	}
	return &object.Datacenter{
		Common: object.Common{
			InventoryPath: name,
			Reference:     ref,
		},
	}
}

func createMockCluster(id, name string) *object.ClusterComputeResource {
	ref := types.ManagedObjectReference{
		Type:  "ClusterComputeResource",
		Value: id,
	}
	return &object.ClusterComputeResource{
		Common: object.Common{
			InventoryPath: name,
			Reference:     ref,
		},
	}
}

func createMockResourcePool(id, name string) *object.ResourcePool {
	ref := types.ManagedObjectReference{
		Type:  "ResourcePool",
		Value: id,
	}
	return &object.ResourcePool{
		Common: object.Common{
			InventoryPath: name,
			Reference:     ref,
		},
	}
}

func createMockFolder(id, name string) *object.Folder {
	ref := types.ManagedObjectReference{
		Type:  "Folder",
		Value: id,
	}
	return &object.Folder{
		Common: object.Common{
			InventoryPath: name,
			Reference:     ref,
		},
	}
}