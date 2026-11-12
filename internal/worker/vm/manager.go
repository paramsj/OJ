package vm

import (
	"net"
	"sync"

	"github.com/firecracker-microvm/firecracker-go-sdk"
)

type VM struct {
	UUID    string
	IP      net.IP
	Machine *firecracker.Machine
	Cfg     firecracker.Config
}

type Manager struct {
	mu     sync.Mutex
	VMPool map[string]*VM
}

func NewManager() *Manager {
	return &Manager{
		VMPool: make(map[string]*VM),
	}
}

func (m *Manager) Add(vm *VM) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.VMPool[vm.UUID] = vm
}

func (m *Manager) Delete(UUID string) (*VM, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if key, ok := m.VMPool[UUID]; ok {
		delete(m.VMPool, UUID)
		return key, ok
	}

	return nil, false

}
