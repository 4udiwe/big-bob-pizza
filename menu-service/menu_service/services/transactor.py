from abc import ABC, abstractmethod
from typing import Callable, TypeVar, Awaitable
from sqlalchemy.ext.asyncio import AsyncSession

T = TypeVar("T")


class Transactor(ABC):
    @abstractmethod
    async def within_transaction(self, fn: Callable[[AsyncSession], Awaitable[T]]) -> T:
        """Выполняет функцию в рамках транзакции"""
        pass


class AsyncSessionTransactor(Transactor):
    def __init__(self, session: AsyncSession):
        self.session = session
    
    async def within_transaction(self, fn: Callable[[AsyncSession], Awaitable[T]]) -> T:
        try:
            result = await fn(self.session)
            await self.session.commit()
            return result
        except Exception:
            await self.session.rollback()
            raise


