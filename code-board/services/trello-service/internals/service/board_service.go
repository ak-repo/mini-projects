package service

import (
	"errors"

	"github.com/ak-repo/code-board/pkg/dto"
	"github.com/ak-repo/code-board/services/trello-service/internals/models"
	"github.com/ak-repo/code-board/services/trello-service/internals/repo"
)

// BoardService defines all operations related to Boards
type BoardService interface {
	// Create a new board
	CreateBoard(req *dto.CreateBoardRequest) (*dto.BoardResponse, error)

	// Get all boards for a user (owner)
	GetBoardsByUser(userID uint) ([]dto.BoardResponse, error)

	// Get single board by ID
	GetBoardByID(id uint) (*dto.BoardResponse, error)

	// Update board details (e.g. name)
	UpdateBoard(req *dto.UpdateBoardRequest) (*dto.BoardResponse, error)

	// Delete a board and all related lists/cards
	DeleteBoard(id uint) error
}

type boardService struct {
	repo repo.BoardRepo
}

func NewBoardService(repo repo.BoardRepo) BoardService {

	return &boardService{repo: repo}
}

func (s *boardService) CreateBoard(req *dto.CreateBoardRequest) (*dto.BoardResponse, error) {

	board := &models.Board{
		Name:    req.Name,
		OwnerID: req.OwnerID,
	}
	if err := s.repo.CreateBoard(board); err != nil {
		return nil, err
	}

	response := &dto.BoardResponse{
		ID:        board.ID,
		Name:      board.Name,
		OwnerID:   board.OwnerID,
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
	}

	return response, nil

}

func (s *boardService) GetBoardsByUser(userID uint) ([]dto.BoardResponse, error) {

	data, err := s.repo.GetBoardsByUser(userID)
	if err != nil {
		return nil, err
	}

	boards := []dto.BoardResponse{}
	for _, v := range data {

		board := dto.BoardResponse{
			ID:        v.ID,
			Name:      v.Name,
			OwnerID:   v.OwnerID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		lists := []dto.ListResponse{}
		for _, l := range board.Lists {
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
			lists = append(lists, list)

		}
		boards = append(boards, board)

	}

	return boards, nil

}

func (s *boardService) GetBoardByID(id uint) (*dto.BoardResponse, error) {

	data, err := s.repo.GetBoardByID(id)
	if err != nil || data == nil {
		return nil, errors.New("no board found " + err.Error())
	}
	board := &dto.BoardResponse{
		ID:      data.ID,
		Name:    data.Name,
		OwnerID: data.OwnerID,

		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	lists := []dto.ListResponse{}
	for _, l := range board.Lists {
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
		lists = append(lists, list)

	}

	return board, nil

}

func (s *boardService) UpdateBoard(req *dto.UpdateBoardRequest) (*dto.BoardResponse, error) {
	data, err := s.repo.GetBoardByID(req.OwnerID)
	if err != nil || data == nil {
		return nil, errors.New("no board found " + err.Error())
	}

	if req.Name != "" {
		data.Name = req.Name
	}

	if err := s.repo.UpdateBoard(data); err != nil {
		return nil, errors.New("board update failed:  " + err.Error())

	}

	board := &dto.BoardResponse{
		ID:        data.ID,
		Name:      data.Name,
		OwnerID:   data.OwnerID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return board, nil

}

func (s *boardService) DeleteBoard(id uint) error {
	return s.repo.DeleteBoard(id)
}
