package service

import (
	"errors"

	"github.com/ak-repo/code-board/pkg/dto"
	"github.com/ak-repo/code-board/services/trello-service/internals/models"
	"github.com/ak-repo/code-board/services/trello-service/internals/repo"
)

type ListService interface {
	// Create a new List
	CreateList(req *dto.CreateListRequest) error
	// Get all Lists for a Board
	GetListsByBoard(boardID uint) ([]dto.ListResponse, error)

	// Get single List by ID
	 GetListByID(id uint) (*dto.ListResponse, error) 

	// Update List details (e.g. name)
	UpdateList(req *dto.UpdateListRequest) (*dto.ListResponse, error)

	// Delete a List and all related lists/cards
	DeleteList(id uint) error
}

type listService struct {
	repo repo.ListRepo
}

func NewListService(repo repo.ListRepo) ListService {

	return &listService{repo: repo}
}

func (s *listService) CreateList(req *dto.CreateListRequest) error {

	list := &models.List{
		Name:     req.Name,
		BoardID:  req.BoardID,
		Position: req.Position,
	}
	return s.repo.CreateList(list)

}

func (s *listService) GetListsByBoard(boardID uint) ([]dto.ListResponse, error) {

	data, err := s.repo.GetListsByBoard(boardID)
	if err != nil {
		return nil, err
	}

	response := []dto.ListResponse{}
	for _, l := range data {
		list := dto.ListResponse{
			ID:        l.ID,
			Name:      l.Name,
			BoardID:   l.BoardID,
			Position:  l.Position,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		}
		cards := []dto.CardResponse{}
		for _, c := range l.Cards {
			card := dto.CardResponse{
				ID:          c.ID,
				Title:       c.Title,
				Description: c.Description,
				ListID:      c.ID,
				Position:    c.Position,
				CreatedAt:   c.CreatedAt,
				UpdatedAt:   c.UpdatedAt,
			}
			cards = append(cards, card)
		}
		response = append(response, list)

	}

	return response, nil

}

func (s *listService) GetListByID(id uint) (*dto.ListResponse, error) {

	data, err := s.repo.GetListByID(id)
	if err != nil || data == nil {
		return nil, errors.New("no board found " + err.Error())
	}

	response := &dto.ListResponse{
		ID:        data.ID,
		Name:      data.Name,
		BoardID:   data.BoardID,
		Position:  data.Position,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	cards := []dto.CardResponse{}
	for _, c := range data.Cards {
		card := dto.CardResponse{
			ID:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			ListID:      c.ID,
			Position:    c.Position,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
		cards = append(cards, card)
	}

	return response, nil

}

func (s *listService) UpdateList(req *dto.UpdateListRequest) (*dto.ListResponse, error) {
	data, err := s.repo.GetListByID(req.ID)

	if err != nil || data == nil {
		return nil, errors.New("no list found " + err.Error())
	}

	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Position != 0 {
		data.Position = req.Position
	}

	if err := s.repo.UpdateList(data); err != nil {
		return nil, err
	}

	response := &dto.ListResponse{
		ID:        data.ID,
		Name:      data.Name,
		BoardID:   data.BoardID,
		Position:  data.Position,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	cards := []dto.CardResponse{}
	for _, c := range data.Cards {
		card := dto.CardResponse{
			ID:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			ListID:      c.ID,
			Position:    c.Position,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
		cards = append(cards, card)
	}

	return response, nil

}

func (s *listService) DeleteList(id uint) error {
	return s.repo.DeleteList(id)
}
