from uuid import UUID
from decimal import Decimal
from datetime import datetime
from typing import Optional
from sqlalchemy.ext.asyncio import AsyncSession

from menu_service.domain.dish import Dish
from menu_service.domain.promotion import Promotion
from menu_service.repositories.dish_repository import DishRepository, PostgresDishRepository
from menu_service.repositories.promotion_repository import PromotionRepository, PostgresPromotionRepository
from menu_service.services.transactor import Transactor, AsyncSessionTransactor


class MenuServiceError(Exception):
    pass


class DishNotFoundError(MenuServiceError):
    pass


class PromotionNotFoundError(MenuServiceError):
    pass


class MenuService:
    def __init__(
        self,
        dish_repo: DishRepository,
        promotion_repo: PromotionRepository,
        transactor: Transactor,
    ):
        self.dish_repo = dish_repo
        self.promotion_repo = promotion_repo
        self.transactor = transactor
    
    async def add_dish(
        self,
        name: str,
        description: Optional[str],
        price: Decimal,
        category: str,
    ) -> Dish:
        """Команда: добавить блюдо"""
        dish = Dish.create(
            name=name,
            description=description,
            price=price,
            category=category,
        )
        
        async def _do(session: AsyncSession):
            dish_repo = PostgresDishRepository(session)
            created_dish = await dish_repo.create(dish)
            return created_dish
        
        return await self.transactor.within_transaction(_do)
    
    async def update_dish(
        self,
        dish_id: UUID,
        name: Optional[str] = None,
        description: Optional[str] = None,
        price: Optional[Decimal] = None,
        category: Optional[str] = None,
        is_available: Optional[bool] = None,
    ) -> Dish:
        """Команда: изменить блюдо"""
        async def _do(session: AsyncSession):
            dish_repo = PostgresDishRepository(session)
            
            dish = await dish_repo.get_by_id(dish_id)
            if dish is None:
                raise DishNotFoundError(f"Dish {dish_id} not found")
            
            # Обновляем поля
            if name is not None:
                dish.name = name
            if description is not None:
                dish.description = description
            if price is not None:
                dish.price = price
            if category is not None:
                dish.category = category
            if is_available is not None:
                dish.is_available = is_available
            
            dish.updated_at = datetime.utcnow()
            
            updated_dish = await dish_repo.update(dish)
            return updated_dish
        
        return await self.transactor.within_transaction(_do)
    
    async def delete_dish(self, dish_id: UUID) -> None:
        """Команда: удалить блюдо"""
        async def _do(session: AsyncSession):
            dish_repo = PostgresDishRepository(session)
            
            dish = await dish_repo.get_by_id(dish_id)
            if dish is None:
                raise DishNotFoundError(f"Dish {dish_id} not found")
            
            await dish_repo.delete(dish_id)
        
        await self.transactor.within_transaction(_do)
    
    async def activate_promotion(self, promotion_id: UUID) -> Promotion:
        """Команда: активировать акцию (автор - админ)"""
        async def _do(session: AsyncSession):
            promotion_repo = PostgresPromotionRepository(session)
            
            promotion = await promotion_repo.get_by_id(promotion_id)
            if promotion is None:
                raise PromotionNotFoundError(f"Promotion {promotion_id} not found")
            
            promotion.activate()
            
            updated_promotion = await promotion_repo.update(promotion)
            return updated_promotion
        
        return await self.transactor.within_transaction(_do)
    
    async def get_dish(self, dish_id: UUID) -> Dish:
        """Получить блюдо по ID"""
        dish = await self.dish_repo.get_by_id(dish_id)
        if dish is None:
            raise DishNotFoundError(f"Dish {dish_id} not found")
        return dish
    
    async def get_all_dishes(self, limit: int = 100, offset: int = 0) -> list[Dish]:
        """Получить все блюда"""
        return await self.dish_repo.get_all(limit=limit, offset=offset)
    
    async def get_promotion(self, promotion_id: UUID) -> Promotion:
        """Получить акцию по ID"""
        promotion = await self.promotion_repo.get_by_id(promotion_id)
        if promotion is None:
            raise PromotionNotFoundError(f"Promotion {promotion_id} not found")
        return promotion

