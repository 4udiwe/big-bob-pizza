import logging
from contextlib import asynccontextmanager
from fastapi import FastAPI

from config import settings
from menu_service.database.connection import engine, Base
from menu_service.handlers.dish_handler import router as dish_router
from menu_service.handlers.promotion_handler import router as promotion_router

# Настройка логирования
logging.basicConfig(
    level=getattr(logging, settings.log_level.upper()),
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifecycle events для FastAPI"""
    # Startup
    logger.info("Starting menu-service...")
    
    # Проверяем подключение к БД и создаем таблицы
    try:
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        logger.info("Database connection successful, tables created/verified")
    except Exception as e:
        logger.error(f"Failed to connect to database: {e}")
        logger.error("Please check:")
        logger.error(f"  1. PostgreSQL is running on {settings.db_host}:{settings.db_port}")
        logger.error(f"  2. Database '{settings.db_name}' exists")
        logger.error(f"  3. User '{settings.db_user}' has access")
        logger.error(f"  4. Settings in .env file are correct")
        raise
    
    logger.info("menu-service started")
    
    yield
    
    # Shutdown
    logger.info("Shutting down menu-service...")
    await engine.dispose()
    logger.info("menu-service stopped")


app = FastAPI(
    title="Menu Service",
    description="Сервис управления меню и акциями",
    version="1.0.0",
    lifespan=lifespan,
)

# Подключаем роутеры
app.include_router(dish_router)
app.include_router(promotion_router)


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "main:app",
        host=settings.server_host,
        port=settings.server_port,
        reload=True,
    )

