import pytest
from unittest.mock import Mock
from uuid import uuid4

from kitchen.service import (
    KitchenService,
    CannotStartPreparing,
    CannotCompletePreparing,
    ORDER_STATUS_PREPARING,
    ORDER_STATUS_DELIVERING,
)


@pytest.fixture
def repo():
    return Mock()


@pytest.fixture
def tx():
    return Mock()


@pytest.fixture
def service(repo, tx):
    return KitchenService(repo, tx)


def test_start_preparing_success(service, repo, tx):
    order_id = uuid4()

    def tx_fn(fn):
        return fn()

    tx.do.side_effect = tx_fn
    repo.update_order_status.return_value = None

    service.start_preparing(order_id)

    repo.update_order_status.assert_called_once_with(order_id, ORDER_STATUS_PREPARING)


def test_start_preparing_repo_error(service, repo, tx):
    order_id = uuid4()

    def tx_fn(fn):
        return fn()

    tx.do.side_effect = tx_fn
    repo.update_order_status.side_effect = Exception("repo failure")

    with pytest.raises(CannotStartPreparing):
        service.start_preparing(order_id)


def test_start_preparing_tx_error(service, repo, tx):
    order_id = uuid4()

    tx.do.side_effect = Exception("tx failed")

    with pytest.raises(CannotStartPreparing):
        service.start_preparing(order_id)

    repo.update_order_status.assert_not_called()


def test_complete_preparing_success(service, repo, tx):
    order_id = uuid4()

    def tx_fn(fn):
        return fn()

    tx.do.side_effect = tx_fn
    repo.update_order_status.return_value = None

    service.complete_preparing(order_id)

    repo.update_order_status.assert_called_once_with(order_id, ORDER_STATUS_DELIVERING)


def test_complete_preparing_repo_error(service, repo, tx):
    order_id = uuid4()

    def tx_fn(fn):
        return fn()

    tx.do.side_effect = tx_fn
    repo.update_order_status.side_effect = Exception("repo failure")

    with pytest.raises(CannotCompletePreparing):
        service.complete_preparing(order_id)


def test_complete_preparing_tx_error(service, repo, tx):
    order_id = uuid4()

    tx.do.side_effect = Exception("tx failed")

    with pytest.raises(CannotCompletePreparing):
        service.complete_preparing(order_id)

    repo.update_order_status.assert_not_called()
