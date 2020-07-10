package domain

import (
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/seashell/drago/pkg/logger"
)

const DefaultReconcileInterval = 1 // In minutes

var (
	service networkIPAddressLeaseService
	once    sync.Once
)

type networkIPAddressPool struct {
	network  *net.IPNet
	assigned map[uint32]*net.IP
}

type networkPool struct {
	pool map[string]*networkIPAddressPool
	sync.Mutex
}

type NetworkIPAddressLeaseService interface {
	Lease(i *Interface) (*Interface, error)
	PutIPAddress(i *Interface) error
	PopIPAddress(i *Interface) error
	PutNetwork(n *Network) error
	PopNetwork(n *Network) error
}

type networkIPAddressLeaseService struct {
	networks *networkPool
	ir       InterfaceRepository
	nr       NetworkRepository
	log      logger.Logger
}

func NewIPAddressLeaseService(
	ir InterfaceRepository,
	nr NetworkRepository,
	logger logger.Logger,
) (NetworkIPAddressLeaseService, error) {
	once.Do(func() {
		networks := &networkPool{
			pool: make(map[string]*networkIPAddressPool),
		}

		service = networkIPAddressLeaseService{
			networks: networks,
			ir:       ir,
			nr:       nr,
			log:      logger,
		}

		go func() {
			for {
				time.Sleep(DefaultReconcileInterval * time.Minute)

				now := time.Now()

				service.log.Debugf("Reconciling at %v\n", now.Round(0))

				err := service.reconcile()
				if err != nil {
					now := time.Now()

					service.log.Errorf("Reconciliation error at %v: %v\n", now.Round(0), err)
				}
			}
		}()

		now := time.Now()

		service.log.Debugf("Reconciling at %v\n", now.Round(0))

		err := service.reconcile()
		if err != nil {
			now := time.Now()

			service.log.Errorf("Reconciliation error at %v: %v\n", now.Round(0), err)
		}
	})

	return &service, nil
}

func (s *networkIPAddressLeaseService) Lease(i *Interface) (*Interface, error) {
	s.networks.Lock()

	p, ok := s.networks.pool[*i.NetworkID]

	if !ok {
		s.networks.Unlock()

		return nil, errors.New("Network not found.")
	}

	start, finish := getUint32NetworkMarginalIPAddresses(*p.network)

	var address string

	for i := start; i <= finish; i++ {
		if _, ok := p.assigned[i]; ok {
			if i == finish {
				s.networks.Unlock()

				return nil, errors.New("No address available.")
			}

			continue
		}

		ip := make(net.IP, 4)

		binary.BigEndian.PutUint32(ip, i)

		p.assigned[i] = &ip
		address = p.assigned[i].String()

		break
	}

	i.IPAddress = &address

	s.networks.Unlock()

	return i, nil

}

func (s *networkIPAddressLeaseService) PutIPAddress(i *Interface) error {
	s.networks.Lock()

	p, ok := s.networks.pool[*i.NetworkID]

	if !ok {
		s.networks.Unlock()

		return errors.New("Network not found.")
	}

	address := net.ParseIP(*i.IPAddress)
	if address == nil || address.To4() == nil {
		s.networks.Unlock()

		return errors.New("Invalid IPv4 address.")
	}

	ai := binary.BigEndian.Uint32(address)

	if _, ok := p.assigned[ai]; ok {
		s.networks.Unlock()

		return errors.New("Address already exists.")
	}

	p.assigned[ai] = &address

	s.networks.Unlock()

	return nil
}

func (s *networkIPAddressLeaseService) PopIPAddress(i *Interface) error {
	s.networks.Lock()

	p, ok := s.networks.pool[*i.NetworkID]

	if !ok {
		s.networks.Unlock()

		return errors.New("Network not found.")
	}

	address := net.ParseIP(*i.IPAddress)
	if address == nil || address.To4() == nil {
		s.networks.Unlock()

		return errors.New("Invalid IPv4 address.")
	}

	ai := binary.BigEndian.Uint32(address)

	if _, ok := p.assigned[ai]; !ok {
		s.networks.Unlock()

		return errors.New("Address not found.")
	}

	delete(p.assigned, ai)

	s.networks.Unlock()

	return nil
}

func (s *networkIPAddressLeaseService) PutNetwork(n *Network) error {
	s.networks.Lock()

	if _, ok := s.networks.pool[*n.ID]; ok {
		s.networks.Unlock()

		return errors.New("Network already exists.")
	}

	_, net, err := net.ParseCIDR(*n.IPAddressRange)
	if err != nil {
		s.networks.Unlock()

		return err
	}

	s.networks.pool[*n.ID] = &networkIPAddressPool{
		network: net,
	}

	s.networks.Unlock()

	return nil
}

func (s *networkIPAddressLeaseService) PopNetwork(n *Network) error {
	s.networks.Lock()

	if _, ok := s.networks.pool[*n.ID]; !ok {
		if err := s.reconcile(); err != nil {
			s.networks.Unlock()

			return err
		}

		if _, ok := s.networks.pool[*n.ID]; !ok {
			s.networks.Unlock()

			return errors.New("Network not found.")
		}
	}

	delete(s.networks.pool, *n.ID)

	s.networks.Unlock()

	return nil
}

func (s *networkIPAddressLeaseService) reconcile() error {
	s.networks.Lock()

	pi := PageInfo{
		Page:    DefaultPage,
		PerPage: DefaultPerPage,
	}

	ns, _, err := s.nr.FindAll(pi)
	if err != nil {
		s.networks.Unlock()

		return err
	}

	s.networks.pool = make(map[string]*networkIPAddressPool)

	for i := range ns {
		_, n, err := net.ParseCIDR(*ns[i].IPAddressRange)
		if err != nil {
			s.networks.Unlock()

			return err
		}

		s.networks.pool[*ns[i].ID] = &networkIPAddressPool{
			network:  n,
			assigned: make(map[uint32]*net.IP),
		}

		pi.PerPage = MaxPerPage

		is, p, err := s.ir.FindAllByNetworkID(*ns[i].ID, pi)
		if err != nil {
			s.networks.Unlock()

			return err
		}

		for j := range is {
			ip := net.ParseIP(*is[j].IPAddress)
			ipUint32 := binary.BigEndian.Uint32(ip)

			s.networks.pool[*ns[i].ID].assigned[ipUint32] = &ip
		}

		for k := p.Page + 1; k <= p.PageCount; k++ {
			pi.Page = k

			is, p, err = s.ir.FindAllByNetworkID(*ns[i].ID, pi)
			if err != nil {
				s.networks.Unlock()

				return err
			}

			for j := range is {
				ip := net.ParseIP(*is[j].IPAddress)
				ipUint32 := binary.BigEndian.Uint32(ip)

				s.networks.pool[*ns[i].ID].assigned[ipUint32] = &ip
			}
		}
	}

	s.networks.Unlock()

	return nil
}

func getUint32NetworkMarginalIPAddresses(nip net.IPNet) (uint32, uint32) {
	mask := binary.BigEndian.Uint32(nip.Mask)
	start := binary.BigEndian.Uint32(nip.IP)
	finish := (start & mask) | (mask ^ 0xffffffff)

	return start, finish
}
