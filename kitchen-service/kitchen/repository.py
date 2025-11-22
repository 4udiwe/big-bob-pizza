from abc import ABC, abstractmethod
from uuid import UUID
from typing import Any


class OrderRepository(ABC):
    @abstractmethod
    def update_order_status(self, order_id: UUID, status: str) -> None:
        pass

    @abstractmethod
    def get_order_status(self, order_id: UUID) -> str:
        pass


class Transactor(ABC):
    @abstractmethod
    def do(self, fn: Any) -> Any:
        """Executes function within a transaction"""
        pass
