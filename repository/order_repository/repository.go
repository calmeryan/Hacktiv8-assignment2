package order_repository

import "assignment-2/entity"

type OrderRepository interface {
	CreateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*entity.Order, error)
	UpdateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*OrderItem, error)
}
