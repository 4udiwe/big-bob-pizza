from uuid import UUID
from typing import Optional

from .repository import OrderRepository, Transactor


ORDER_STATUS_PREPARING = "preparing"
ORDER_STATUS_DELIVERING = "delivering"


class KitchenError(Exception):
    pass


class CannotStartPreparing(KitchenError):
    pass


class CannotCompletePreparing(KitchenError):
    pass


class KitchenService:
    def __init__(self, repo: OrderRepository, tx: Transactor) -> None:
        self.repo = repo
        self.tx = tx

    def start_preparing(self, order_id: UUID) -> None:
        """Move order into PREPARING"""
        raise NotImplementedError

    def complete_preparing(self, order_id: UUID) -> None:
        """Move order into DELIVERING"""
        raise NotImplementedError
