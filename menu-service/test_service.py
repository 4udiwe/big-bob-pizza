"""
Скрипт для проверки работоспособности menu-service
Запуск: python test_service.py
"""
import asyncio
import httpx
from uuid import uuid4
from decimal import Decimal


BASE_URL = "http://localhost:8000"


async def test_health_check():
    """Проверка health check endpoint"""
    print("1. Проверка health check...")
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{BASE_URL}/health")
        assert response.status_code == 200
        assert response.json() == {"status": "ok"}
        print("   ✓ Health check работает")


async def test_create_dish():
    """Тест создания блюда"""
    print("\n2. Создание блюда...")
    async with httpx.AsyncClient() as client:
        dish_data = {
            "name": "Пицца Маргарита",
            "description": "Классическая пицца с томатами и моцареллой",
            "price": 599.00,
            "category": "Пицца"
        }
        response = await client.post(f"{BASE_URL}/dishes", json=dish_data)
        if response.status_code != 201:
            print(f"   ✗ Ошибка: статус {response.status_code}")
            print(f"   Ответ: {response.text}")
            raise AssertionError(f"Expected 201, got {response.status_code}: {response.text}")
        dish = response.json()
        assert dish["name"] == dish_data["name"]
        # Проверяем price как Decimal (может быть строкой или числом)
        price = float(dish["price"]) if isinstance(dish["price"], str) else dish["price"]
        assert abs(price - dish_data["price"]) < 0.01
        print(f"   ✓ Блюдо создано: {dish['id']}")
        return dish["id"]


async def test_get_dish(dish_id: str):
    """Тест получения блюда"""
    print(f"\n3. Получение блюда {dish_id}...")
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{BASE_URL}/dishes/{dish_id}")
        assert response.status_code == 200
        dish = response.json()
        assert dish["id"] == dish_id
        print(f"   ✓ Блюдо получено: {dish['name']}")


async def test_get_all_dishes():
    """Тест получения всех блюд"""
    print("\n4. Получение всех блюд...")
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{BASE_URL}/dishes")
        assert response.status_code == 200
        dishes = response.json()
        assert isinstance(dishes, list)
        print(f"   ✓ Получено блюд: {len(dishes)}")


async def test_update_dish(dish_id: str):
    """Тест обновления блюда"""
    print(f"\n5. Обновление блюда {dish_id}...")
    async with httpx.AsyncClient() as client:
        update_data = {
            "price": 649.00,
            "is_available": False
        }
        response = await client.put(f"{BASE_URL}/dishes/{dish_id}", json=update_data)
        if response.status_code != 200:
            print(f"   ✗ Ошибка: статус {response.status_code}")
            print(f"   Ответ: {response.text}")
            raise AssertionError(f"Expected 200, got {response.status_code}: {response.text}")
        dish = response.json()
        # Проверяем price как Decimal (может быть строкой или числом)
        price = float(dish["price"]) if isinstance(dish["price"], str) else dish["price"]
        assert abs(price - update_data["price"]) < 0.01
        assert dish["is_available"] == update_data["is_available"]
        print(f"   ✓ Блюдо обновлено: цена={dish['price']}, доступно={dish['is_available']}")


async def test_delete_dish(dish_id: str):
    """Тест удаления блюда"""
    print(f"\n6. Удаление блюда {dish_id}...")
    async with httpx.AsyncClient() as client:
        response = await client.delete(f"{BASE_URL}/dishes/{dish_id}")
        assert response.status_code == 204
        print("   ✓ Блюдо удалено")
        
        # Проверяем, что блюдо действительно удалено
        response = await client.get(f"{BASE_URL}/dishes/{dish_id}")
        assert response.status_code == 404
        print("   ✓ Подтверждено удаление (404 при попытке получить)")


async def test_not_found():
    """Тест обработки несуществующего ресурса"""
    print("\n7. Проверка обработки несуществующего блюда...")
    async with httpx.AsyncClient() as client:
        fake_id = str(uuid4())
        response = await client.get(f"{BASE_URL}/dishes/{fake_id}")
        assert response.status_code == 404
        print("   ✓ Корректная обработка 404")


async def run_all_tests():
    """Запуск всех тестов"""
    print("=" * 50)
    print("Тестирование menu-service")
    print("=" * 50)
    
    try:
        await test_health_check()
        dish_id = await test_create_dish()
        await test_get_dish(dish_id)
        await test_get_all_dishes()
        await test_update_dish(dish_id)
        await test_delete_dish(dish_id)
        await test_not_found()
        
        print("\n" + "=" * 50)
        print("✓ Все тесты пройдены успешно!")
        print("=" * 50)
        
    except AssertionError as e:
        print(f"\n✗ Тест провален: {e}")
        import traceback
        traceback.print_exc()
        return 1
    except httpx.ConnectError:
        print("\n✗ Ошибка: не удалось подключиться к сервису")
        print("   Убедитесь, что сервис запущен на http://localhost:8000")
        return 1
    except httpx.HTTPStatusError as e:
        print(f"\n✗ HTTP ошибка: {e.response.status_code}")
        print(f"   Ответ: {e.response.text}")
        import traceback
        traceback.print_exc()
        return 1
    except Exception as e:
        print(f"\n✗ Неожиданная ошибка: {e}")
        import traceback
        traceback.print_exc()
        return 1
    
    return 0


if __name__ == "__main__":
    exit_code = asyncio.run(run_all_tests())
    exit(exit_code)

