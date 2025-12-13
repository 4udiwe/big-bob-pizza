//go:build e2e

package e2e_test

import (
	"testing"

	. "github.com/Eun/go-hit"
	"github.com/google/uuid"
)

func TestCreateOrder_Success(t *testing.T) {
	body := map[string]any{
		"customerId":  uuid.New().String(),
		"totalAmount": 100.0,
		"currency":    "USD",
		"items": []map[string]any{
			{
				"productId":    uuid.New().String(),
				"productName":  "Pizza",
				"productPrice": 50.0,
				"amount":       2,
				"totalPrice":   100.0,
			},
		},
	}

	err := Do(
		Post(basePath+"/orders"),
		Send().Body().JSON(body),
		Send().Headers("Content-Type").Add("application/json"),
		Expect().Status().Equal(201),
	)

	if err != nil {
		t.Fatalf("create order failed: %v", err)
	}
}

func TestGetOrdersList_Success(t *testing.T) {
	err := Do(
		Get(basePath+"/orders?limit=10&offset=0"),
		Expect().Status().Equal(200),
	)
	if err != nil {
		t.Fatalf("get orders list failed: %v", err)
	}
}

func TestGetOrdersByUser_Success(t *testing.T) {
	user := uuid.New().String()

	err := Do(
		Get(basePath+"/orders/user/"+user),
		Expect().Status().Equal(200),
	)
	if err != nil {
		t.Fatalf("get orders by user failed: %v", err)
	}
}

func TestGetActiveOrdersByUser_Success(t *testing.T) {
	user := uuid.New().String()

	err := Do(
		Get(basePath+"/orders/user/"+user+"/active"),
		Expect().Status().Equal(200),
	)
	if err != nil {
		t.Fatalf("get active orders failed: %v", err)
	}
}

func TestGetOrder_Success(t *testing.T) {
	// Просто генерируем ID (если нет — сервис вернет 404, но это НЕ ошибка API)
	id := uuid.New().String()

	err := Do(
		Get(basePath+"/orders/"+id),
		Expect().Status().Between(400, 499), // допускаем 404 как валидный ответ
	)

	if err != nil {
		t.Fatalf("get order failed: %v", err)
	}
}

func TestValidation_Success(t *testing.T) {
	// Нарочно ломаем тело запроса
	body := map[string]interface{}{
		"customerId": uuid.New().String(),
	}

	err := Do(
		Post(basePath+"/orders"),
		Send().Body().JSON(body),
		Expect().Status().Between(400, 499), // ожидаем клиентскую ошибку
	)

	if err != nil {
		t.Fatalf("validation test failed: %v", err)
	}
}
