from datetime import datetime
from uuid import UUID, uuid4
from decimal import Decimal
from typing import Optional
from dataclasses import dataclass


@dataclass
class Dish:
    id: UUID
    name: str
    description: Optional[str]
    price: Decimal
    category: str
    is_available: bool
    created_at: datetime
    updated_at: datetime
    
    @classmethod
    def create(
        cls,
        name: str,
        description: Optional[str],
        price: Decimal,
        category: str,
    ) -> "Dish":
        now = datetime.utcnow()
        return cls(
            id=uuid4(),
            name=name,
            description=description,
            price=price,
            category=category,
            is_available=True,
            created_at=now,
            updated_at=now,
        )


