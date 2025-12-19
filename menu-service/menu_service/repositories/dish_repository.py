from abc import ABC, abstractmethod
from uuid import UUID
from typing import Optional, List
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, update, delete
from sqlalchemy.orm import selectinload

from menu_service.domain.dish import Dish
from menu_service.database.models import DishModel
from decimal import Decimal


class DishRepository(ABC):
    @abstractmethod
    async def create(self, dish: Dish) -> Dish:
        pass
    
    @abstractmethod
    async def get_by_id(self, dish_id: UUID) -> Optional[Dish]:
        pass
    
    @abstractmethod
    async def update(self, dish: Dish) -> Dish:
        pass
    
    @abstractmethod
    async def delete(self, dish_id: UUID) -> None:
        pass
    
    @abstractmethod
    async def get_all(self, limit: int = 100, offset: int = 0) -> List[Dish]:
        pass


class PostgresDishRepository(DishRepository):
    def __init__(self, session: AsyncSession):
        self.session = session
    
    def _to_domain(self, model: DishModel) -> Dish:
        return Dish(
            id=model.id,
            name=model.name,
            description=model.description,
            price=Decimal(str(model.price)),
            category=model.category,
            is_available=model.is_available,
            created_at=model.created_at,
            updated_at=model.updated_at,
        )
    
    def _to_model(self, dish: Dish) -> DishModel:
        return DishModel(
            id=dish.id,
            name=dish.name,
            description=dish.description,
            price=dish.price,
            category=dish.category,
            is_available=dish.is_available,
            created_at=dish.created_at,
            updated_at=dish.updated_at,
        )
    
    async def create(self, dish: Dish) -> Dish:
        model = self._to_model(dish)
        self.session.add(model)
        await self.session.flush()
        await self.session.refresh(model)
        return self._to_domain(model)
    
    async def get_by_id(self, dish_id: UUID) -> Optional[Dish]:
        result = await self.session.execute(
            select(DishModel).where(DishModel.id == dish_id)
        )
        model = result.scalar_one_or_none()
        if model is None:
            return None
        return self._to_domain(model)
    
    async def update(self, dish: Dish) -> Dish:
        await self.session.execute(
            update(DishModel)
            .where(DishModel.id == dish.id)
            .values(
                name=dish.name,
                description=dish.description,
                price=dish.price,
                category=dish.category,
                is_available=dish.is_available,
                updated_at=dish.updated_at,
            )
        )
        await self.session.flush()
        return dish
    
    async def delete(self, dish_id: UUID) -> None:
        await self.session.execute(
            delete(DishModel).where(DishModel.id == dish_id)
        )
        await self.session.flush()
    
    async def get_all(self, limit: int = 100, offset: int = 0) -> List[Dish]:
        result = await self.session.execute(
            select(DishModel)
            .limit(limit)
            .offset(offset)
            .order_by(DishModel.created_at.desc())
        )
        models = result.scalars().all()
        return [self._to_domain(model) for model in models]


