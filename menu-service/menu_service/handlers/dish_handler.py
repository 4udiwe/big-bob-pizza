from uuid import UUID
from decimal import Decimal
from typing import Optional
from fastapi import APIRouter, Depends, HTTPException, status
from pydantic import BaseModel, Field
from sqlalchemy.ext.asyncio import AsyncSession

from menu_service.services.menu_service import MenuService, DishNotFoundError
from menu_service.database.connection import get_db
from menu_service.repositories.dish_repository import PostgresDishRepository
from menu_service.repositories.promotion_repository import PostgresPromotionRepository
from menu_service.services.transactor import AsyncSessionTransactor

router = APIRouter(prefix="/dishes", tags=["dishes"])


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


# Request/Response модели
class CreateDishRequest(BaseModel):
    name: str = Field(..., min_length=1, max_length=255)
    description: Optional[str] = Field(None, max_length=1000)
    price: Decimal = Field(..., gt=0)
    category: str = Field(..., min_length=1, max_length=100)


class UpdateDishRequest(BaseModel):
    name: Optional[str] = Field(None, min_length=1, max_length=255)
    description: Optional[str] = Field(None, max_length=1000)
    price: Optional[Decimal] = Field(None, gt=0)
    category: Optional[str] = Field(None, min_length=1, max_length=100)
    is_available: Optional[bool] = None


class DishResponse(BaseModel):
    id: UUID
    name: str
    description: Optional[str]
    price: Decimal
    category: str
    is_available: bool
    created_at: str
    updated_at: str
    
    class Config:
        from_attributes = True


@router.post("/", response_model=DishResponse, status_code=status.HTTP_201_CREATED)
async def create_dish(
    request: CreateDishRequest,
    service: MenuService = Depends(get_menu_service),
):
    """Команда: добавить блюдо"""
    try:
        dish = await service.add_dish(
            name=request.name,
            description=request.description,
            price=request.price,
            category=request.category,
        )
        return DishResponse(
            id=dish.id,
            name=dish.name,
            description=dish.description,
            price=dish.price,
            category=dish.category,
            is_available=dish.is_available,
            created_at=dish.created_at.isoformat(),
            updated_at=dish.updated_at.isoformat(),
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@router.put("/{dish_id}", response_model=DishResponse)
async def update_dish(
    dish_id: UUID,
    request: UpdateDishRequest,
    service: MenuService = Depends(get_menu_service),
):
    """Команда: изменить блюдо"""
    try:
        dish = await service.update_dish(
            dish_id=dish_id,
            name=request.name,
            description=request.description,
            price=request.price,
            category=request.category,
            is_available=request.is_available,
        )
        return DishResponse(
            id=dish.id,
            name=dish.name,
            description=dish.description,
            price=dish.price,
            category=dish.category,
            is_available=dish.is_available,
            created_at=dish.created_at.isoformat(),
            updated_at=dish.updated_at.isoformat(),
        )
    except DishNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e),
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@router.delete("/{dish_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_dish(
    dish_id: UUID,
    service: MenuService = Depends(get_menu_service),
):
    """Команда: удалить блюдо"""
    try:
        await service.delete_dish(dish_id)
    except DishNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e),
        )
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=str(e),
        )


@router.get("/{dish_id}", response_model=DishResponse)
async def get_dish(
    dish_id: UUID,
    service: MenuService = Depends(get_menu_service),
):
    """Получить блюдо по ID"""
    try:
        dish = await service.get_dish(dish_id)
        return DishResponse(
            id=dish.id,
            name=dish.name,
            description=dish.description,
            price=dish.price,
            category=dish.category,
            is_available=dish.is_available,
            created_at=dish.created_at.isoformat(),
            updated_at=dish.updated_at.isoformat(),
        )
    except DishNotFoundError as e:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(e),
        )


@router.get("/", response_model=list[DishResponse])
async def get_all_dishes(
    limit: int = 100,
    offset: int = 0,
    service: MenuService = Depends(get_menu_service),
):
    """Получить все блюда"""
    dishes = await service.get_all_dishes(limit=limit, offset=offset)
    return [
        DishResponse(
            id=dish.id,
            name=dish.name,
            description=dish.description,
            price=dish.price,
            category=dish.category,
            is_available=dish.is_available,
            created_at=dish.created_at.isoformat(),
            updated_at=dish.updated_at.isoformat(),
        )
        for dish in dishes
    ]

