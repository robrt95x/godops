package usecase_test

import (
	"testing"
	"time"

	pkgLogger "github.com/robrt95x/godops/pkg/logger"
	"github.com/robrt95x/godops/services/order/internal/entity"
	"github.com/robrt95x/godops/services/order/internal/errors"
	"github.com/robrt95x/godops/services/order/internal/infra/memory"
	"github.com/robrt95x/godops/services/order/internal/usecase"
)

func TestGetOrderByIDCase_Execute(t *testing.T) {
	// Setup
	repo := memory.NewOrderMemoryRepository()
	testLogger := pkgLogger.Setup(pkgLogger.NewDefaultConfig())
	uc := usecase.NewGetOrderByIDCase(repo, testLogger)

	// Create a test order
	testOrder := &entity.Order{
		ID:     "test-order-123",
		UserID: "user-456",
		Items: []entity.OrderItem{
			{
				ProductID: "product-1",
				Quantity:  2,
				Price:     29.99,
			},
		},
		Status:          entity.Pending,
		CouponCode:      "DISCOUNT10",
		Total:           59.98,
		ShippingAddress: "123 Test St",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save the test order
	err := repo.Save(testOrder)
	if err != nil {
		t.Fatalf("Failed to save test order: %v", err)
	}

	t.Run("should return order when found", func(t *testing.T) {
		// Execute
		result, err := uc.Execute("test-order-123")

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("Expected order, got nil")
		}
		if result.ID != "test-order-123" {
			t.Errorf("Expected order ID 'test-order-123', got '%s'", result.ID)
		}
		if result.UserID != "user-456" {
			t.Errorf("Expected user ID 'user-456', got '%s'", result.UserID)
		}
		if len(result.Items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(result.Items))
		}
		if result.Total != 59.98 {
			t.Errorf("Expected total 59.98, got %f", result.Total)
		}
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		// Execute
		result, err := uc.Execute("non-existent-order")

		// Assert
		if err != errors.ErrOrderNotFound {
			t.Errorf("Expected ErrOrderNotFound, got %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("should return error for invalid order ID", func(t *testing.T) {
		// Execute
		result, err := uc.Execute("")

		// Assert
		if err != errors.ErrOrderInvalidID {
			t.Errorf("Expected ErrOrderInvalidID, got %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}

func TestMemoryRepository_Isolation(t *testing.T) {
	// This test demonstrates that the memory repository provides proper isolation
	repo := memory.NewOrderMemoryRepository()

	// Create test orders
	order1 := &entity.Order{
		ID:     "order-1",
		UserID: "user-1",
		Items:  []entity.OrderItem{{ProductID: "product-1", Quantity: 1, Price: 10.0}},
		Status: entity.Pending,
		Total:  10.0,
	}

	order2 := &entity.Order{
		ID:     "order-2",
		UserID: "user-2",
		Items:  []entity.OrderItem{{ProductID: "product-2", Quantity: 2, Price: 15.0}},
		Status: entity.Completed,
		Total:  30.0,
	}

	// Save orders
	repo.Save(order1)
	repo.Save(order2)

	// Verify count
	if repo.Count() != 2 {
		t.Errorf("Expected 2 orders, got %d", repo.Count())
	}

	// Clear repository
	repo.Clear()

	// Verify repository is empty
	if repo.Count() != 0 {
		t.Errorf("Expected 0 orders after clear, got %d", repo.Count())
	}

	// Verify orders are not found
	_, err := repo.FindByID("order-1")
	if err == nil {
		t.Error("Expected error when finding order after clear")
	}
}
