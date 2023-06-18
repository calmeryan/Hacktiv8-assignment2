package service

import (
	"assignment-2/dto"
	"assignment-2/entity"
	"assignment-2/repository/order_repository"
	"net/http"
)

type orderService struct {
	orderRepo   order_repository.OrderRepository
	itemService ItemService
}

type OrderService interface {
	CreateOrder(payload dto.NewOrderRequest) (*dto.NewOrderResponse, error)
	UpdateOrder(orderId int, payload dto.NewOrderRequest) (*dto.GetOrderResponse, error)
}

func NewOrderService(orderRepo order_repository.OrderRepository, itemService ItemService) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		itemService: itemService,
	}
}

func (o *orderService) UpdateOrder(orderId int, payload dto.NewOrderRequest) (*dto.GetOrderResponse, error) {
	itemCodes := payload.ItemsToItemCode() //[]string{"123", "456"}

	_, err := o.itemService.FindItemsByItemCodes(itemCodes)

	if err != nil {
		return nil, err
	}

	orderPayload := entity.Order{
		OrderId:      orderId,
		OrderedAt:    payload.OrderedAt,
		CustomerName: payload.CustomerName,
	}

	itemsPayload := []entity.Item{}

	for _, eachItem := range payload.Items {
		item := entity.Item{
			ItemCode:    eachItem.ItemCode,
			Quantity:    eachItem.Quantity,
			Description: eachItem.Description,
		}

		itemsPayload = append(itemsPayload, item)
	}

	orderItem, err := o.orderRepo.UpdateOrder(orderPayload, itemsPayload)

	if err != nil {
		return nil, err
	}

	itemsResponse := []dto.ItemResponse{}

	for _, eachItem := range orderItem.Items {
		itemResponse := eachItem.ItemToItemResponse()

		itemsResponse = append(itemsResponse, itemResponse)
	}

	result := dto.GetOrderResponse{
		Code: http.StatusOK,
		Data: dto.OrderResponse{
			Id:           orderItem.Order.OrderId,
			CreatedAt:    orderItem.Order.CreatedAt,
			UpdatedAt:    orderItem.Order.UpdatedAt,
			CustomerName: orderItem.Order.CustomerName,
			Items:        itemsResponse,
		},
	}

	return &result, nil
}

func (o *orderService) CreateOrder(payload dto.NewOrderRequest) (*dto.NewOrderResponse, error) {
	orderPayload := entity.Order{
		OrderedAt:    payload.OrderedAt,
		CustomerName: payload.CustomerName,
	}

	itemsPayload := []entity.Item{}

	for _, eachItem := range payload.Items {
		item := entity.Item{
			ItemCode:    eachItem.ItemCode,
			Quantity:    eachItem.Quantity,
			Description: eachItem.Description,
		}

		itemsPayload = append(itemsPayload, item)
	}

	newOrder, err := o.orderRepo.CreateOrder(orderPayload, itemsPayload)

	if err != nil {
		return nil, err
	}

	response := &dto.NewOrderResponse{
		Message: "Success",
		Data: dto.NewOrderRequest{
			OrderedAt:    newOrder.OrderedAt,
			CustomerName: newOrder.CustomerName,
		},
		StatusCode: http.StatusCreated,
	}

	return response, nil
}
