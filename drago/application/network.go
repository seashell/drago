package application

import (
	structs "github.com/seashell/drago/drago/application/structs"
	domain "github.com/seashell/drago/drago/domain"
)

// NetworkService
type NetworkService interface {
	Get(*structs.GetNetworkInput) (*structs.GetNetworkOutput, error)
	Create(*structs.CreateNetworkInput) (*structs.CreateNetworkOutput, error)
	Update(*structs.UpdateNetworkInput) (*structs.UpdateNetworkOutput, error)
	Delete(*structs.DeleteNetworkInput) (*structs.DeleteNetworkOutput, error)
	List(*structs.ListNetworksInput) (*structs.ListNetworksOutput, error)
}

type networkService struct {
	repo domain.NetworkRepository
}

// NewNetworkService
func NewNetworkService(nr domain.NetworkRepository) NetworkService {
	return &networkService{
		repo: nr,
	}
}

// GetByID returns a Network entity by ID
func (s *networkService) Get(in *structs.GetNetworkInput) (*structs.GetNetworkOutput, error) {

	entity, err := s.repo.GetByID(*in.ID)
	if err != nil {
		return nil, err
	}

	out := &structs.GetNetworkOutput{
		ID: entity.ID,
	}

	return out, nil
}

// Create creates a new Network entity
func (s *networkService) Create(in *structs.CreateNetworkInput) (*structs.CreateNetworkOutput, error) {

	entity := &domain.Network{
		Name:           in.Name,
		IPAddressRange: in.IPAddressRange,
	}

	id, err := s.repo.Create(entity)
	if err != nil {
		return nil, err
	}

	out := &structs.CreateNetworkOutput{
		ID: id,
	}

	return out, nil
}

// Update updates an already existing Network entity
func (s *networkService) Update(in *structs.UpdateNetworkInput) (*structs.UpdateNetworkOutput, error) {

	update := &domain.Network{
		Name:           in.Name,
		IPAddressRange: in.IPAddressRange,
	}

	entity, err := s.repo.GetByID(*in.ID)
	if err != nil {
		return nil, err
	}

	entity = entity.Merge(update)

	id, err := s.repo.Update(entity)
	if err != nil {
		return nil, err
	}

	out := &structs.UpdateNetworkOutput{
		ID: id,
	}

	return out, nil
}

// Delete :
func (s *networkService) Delete(in *structs.DeleteNetworkInput) (*structs.DeleteNetworkOutput, error) {

	id, err := s.repo.DeleteByID(in.ID)
	if err != nil {
		return nil, err
	}

	out := &structs.DeleteNetworkOutput{
		ID: *id,
	}

	return out, nil
}

// List :
func (s *networkService) List(in *structs.ListNetworksInput) (*structs.ListNetworksOutput, error) {

	pageInfo := domain.PageInfo{
		Page:    in.Page,
		PerPage: in.PerPage,
	}

	entities, page, err := s.repo.FindAll(pageInfo)
	if err != nil {
		return nil, err
	}

	items := []*structs.GetNetworkOutput{}
	for _, entity := range entities {
		items = append(items, &structs.GetNetworkOutput{
			ID:             entity.ID,
			Name:           entity.Name,
			IPAddressRange: entity.IPAddressRange,
		})
	}

	out := &structs.ListNetworksOutput{
		PageOutput: structs.PageOutput{
			Page:       page.Page,
			PerPage:    page.PerPage,
			PageCount:  page.PageCount,
			TotalCount: page.TotalCount,
		},
		Items: items,
	}

	return out, nil
}
