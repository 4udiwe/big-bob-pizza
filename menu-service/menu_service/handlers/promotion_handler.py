from uuid import UUID
from decimal import Decimal
from typing import Optional
from datetime import datetime
from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, Field
from sqlalchemy.ext.asyncio import AsyncSession

from menu_service.services.menu_service import MenuService, PromotionNotFoundError
from menu_service.database.connection import get_db
from menu_service.repositories.dish_repository import PostgresDishRepository
from menu_service.repositories.promotion_repository import PostgresPromotionRepository
from menu_service.services.transactor import AsyncSessionTransactor

router = APIRouter(prefix="/promotions", tags=["promotions"])


def get_menu_service(db: AsyncSession = Depends(get_db)) -> MenuService:
    """Dependency для получения MenuService"""
    dish_repo = PostgresDishRepository(db)
    promotion_repo = PostgresPromotionRepository(db)
    transactor = AsyncSessionTransactor(db)
    
    return MenuService(
        dish_repo=dish_repo,
        promotion_repo=promotion_repo,
        transactor=transactor,
    )


class PromotionResponse(BaseModel):
    id: UUID
    name: str
    description: Optional[str]
    discount_percent: Decimal
    dish_ids: list[UUID]
    is_active: bool
    starts_at: str
    ends_at: Optional[str]
    created_at: str
    updated_at: str
    
    class Config:
        from_attributes = True


@router.post("/{promotion_id}/activate", response_model=PromotionResponse)
async def activate_promotion(
    promotion_id: UUID,
    service: MenuService = Depends(get_menu_service),
):
    """Команда: активировать акцию (автор - админ)"""
    try:
        promotion = await service.activate_promotion(promotion_id)
        return PromotionResponse(
            id=promotion.id,
            name=promotion.name,
            description=promotion.description,
            discount_percent=promotion.discount_percent,
            dish_ids=promotion.dish_ids,
            is_active=promotion.is_active,
            starts_at=promotion.starts_at.isoformat(),
            ends_at=promotion.ends_at.isoformat() if promotion.ends_at else None,
            created_at=promotion.created_at.isoformat(),
            updated_at=promotion.updated_at.isoformat(),
        )
    except PromotionNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e),
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@router.get("/{promotion_id}", response_model=PromotionResponse)
async def get_promotion(
    promotion_id: UUID,
    service: MenuService = Depends(get_menu_service),
):
    """Получить акцию по ID"""
    try:
        promotion = await service.get_promotion(promotion_id)
        return PromotionResponse(
            id=promotion.id,
            name=promotion.name,
            description=promotion.description,
            discount_percent=promotion.discount_percent,
            dish_ids=promotion.dish_ids,
            is_active=promotion.is_active,
            starts_at=promotion.starts_at.isoformat(),
            ends_at=promotion.ends_at.isoformat() if promotion.ends_at else None,
            created_at=promotion.created_at.isoformat(),
            updated_at=promotion.updated_at.isoformat(),
        )
    except PromotionNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e),
        )

