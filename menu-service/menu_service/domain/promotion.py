from datetime import datetime
from uuid import UUID, uuid4
from decimal import Decimal
from typing import Optional
from dataclasses import dataclass


@dataclass
class Promotion:
    id: UUID
    name: str
    description: Optional[str]
    discount_percent: Decimal
    dish_ids: list[UUID]
    is_active: bool
    starts_at: datetime
    ends_at: Optional[datetime]
    created_at: datetime
    updated_at: datetime
    
    @classmethod
    def create(
        cls,
        name: str,
        description: Optional[str],
        discount_percent: Decimal,
        dish_ids: list[UUID],
        starts_at: datetime,
        ends_at: Optional[datetime] = None,
    ) -> "Promotion":
        now = datetime.utcnow()
        return cls(
            id=uuid4(),
            name=name,
            description=description,
            discount_percent=discount_percent,
            dish_ids=dish_ids,
            is_active=False,
            starts_at=starts_at,
            ends_at=ends_at,
            created_at=now,
            updated_at=now,
        )
    
    def activate(self) -> None:
        """Активирует акцию"""
        self.is_active = True
        self.updated_at = datetime.utcnow()


