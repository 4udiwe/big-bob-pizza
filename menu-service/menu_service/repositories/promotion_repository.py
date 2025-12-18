from abc import ABC, abstractmethod
from uuid import UUID
from typing import Optional, List
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, update
from datetime import datetime

from menu_service.domain.promotion import Promotion
from menu_service.database.models import PromotionModel
from decimal import Decimal


class PromotionRepository(ABC):
    @abstractmethod
    async def create(self, promotion: Promotion) -> Promotion:
        pass
    
    @abstractmethod
    async def get_by_id(self, promotion_id: UUID) -> Optional[Promotion]:
        pass
    
    @abstractmethod
    async def update(self, promotion: Promotion) -> Promotion:
        pass
    
    @abstractmethod
    async def get_all(self, limit: int = 100, offset: int = 0) -> List[Promotion]:
        pass


class PostgresPromotionRepository(PromotionRepository):
    def __init__(self, session: AsyncSession):
        self.session = session
    
    def _to_domain(self, model: PromotionModel) -> Promotion:
        return Promotion(
            id=model.id,
            name=model.name,
            description=model.description,
            discount_percent=Decimal(str(model.discount_percent)),
            dish_ids=[UUID(dish_id) for dish_id in model.dish_ids],
            is_active=model.is_active,
            starts_at=model.starts_at,
            ends_at=model.ends_at,
            created_at=model.created_at,
            updated_at=model.updated_at,
        )
    
    def _to_model(self, promotion: Promotion) -> PromotionModel:
        return PromotionModel(
            id=promotion.id,
            name=promotion.name,
            description=promotion.description,
            discount_percent=promotion.discount_percent,
            dish_ids=[str(dish_id) for dish_id in promotion.dish_ids],
            is_active=promotion.is_active,
            starts_at=promotion.starts_at,
            ends_at=promotion.ends_at,
            created_at=promotion.created_at,
            updated_at=promotion.updated_at,
        )
    
    async def create(self, promotion: Promotion) -> Promotion:
        model = self._to_model(promotion)
        self.session.add(model)
        await self.session.flush()
        await self.session.refresh(model)
        return self._to_domain(model)
    
    async def get_by_id(self, promotion_id: UUID) -> Optional[Promotion]:
        result = await self.session.execute(
            select(PromotionModel).where(PromotionModel.id == promotion_id)
        )
        model = result.scalar_one_or_none()
        if model is None:
            return None
        return self._to_domain(model)
    
    async def update(self, promotion: Promotion) -> Promotion:
        await self.session.execute(
            update(PromotionModel)
            .where(PromotionModel.id == promotion.id)
            .values(
                name=promotion.name,
                description=promotion.description,
                discount_percent=promotion.discount_percent,
                dish_ids=[str(dish_id) for dish_id in promotion.dish_ids],
                is_active=promotion.is_active,
                starts_at=promotion.starts_at,
                ends_at=promotion.ends_at,
                updated_at=promotion.updated_at,
            )
        )
        await self.session.flush()
        return promotion
    
    async def get_all(self, limit: int = 100, offset: int = 0) -> List[Promotion]:
        result = await self.session.execute(
            select(PromotionModel)
            .limit(limit)
            .offset(offset)
            .order_by(PromotionModel.created_at.desc())
        )
        models = result.scalars().all()
        return [self._to_domain(model) for model in models]


