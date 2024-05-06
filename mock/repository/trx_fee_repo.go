// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/trx_fee_repo.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	repository "github.com/jaime1129/fedex/internal/repository"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// BatchInsertUniTrxFee mocks base method.
func (m *MockRepository) BatchInsertUniTrxFee(fees []repository.UniTrxFee) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchInsertUniTrxFee", fees)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchInsertUniTrxFee indicates an expected call of BatchInsertUniTrxFee.
func (mr *MockRepositoryMockRecorder) BatchInsertUniTrxFee(fees interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchInsertUniTrxFee", reflect.TypeOf((*MockRepository)(nil).BatchInsertUniTrxFee), fees)
}

// BatchRecordHistoricalTrx mocks base method.
func (m *MockRepository) BatchRecordHistoricalTrx(fees []repository.UniTrxFee, symbol string, maxBlock uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchRecordHistoricalTrx", fees, symbol, maxBlock)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchRecordHistoricalTrx indicates an expected call of BatchRecordHistoricalTrx.
func (mr *MockRepositoryMockRecorder) BatchRecordHistoricalTrx(fees, symbol, maxBlock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchRecordHistoricalTrx", reflect.TypeOf((*MockRepository)(nil).BatchRecordHistoricalTrx), fees, symbol, maxBlock)
}

// Close mocks base method.
func (m *MockRepository) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}

// GetMaxBlockNum mocks base method.
func (m *MockRepository) GetMaxBlockNum(symbol string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaxBlockNum", symbol)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMaxBlockNum indicates an expected call of GetMaxBlockNum.
func (mr *MockRepositoryMockRecorder) GetMaxBlockNum(symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaxBlockNum", reflect.TypeOf((*MockRepository)(nil).GetMaxBlockNum), symbol)
}

// GetTrxFee mocks base method.
func (m *MockRepository) GetTrxFee(txHash string) (*repository.UniTrxFee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrxFee", txHash)
	ret0, _ := ret[0].(*repository.UniTrxFee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrxFee indicates an expected call of GetTrxFee.
func (mr *MockRepositoryMockRecorder) GetTrxFee(txHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrxFee", reflect.TypeOf((*MockRepository)(nil).GetTrxFee), txHash)
}

// ListTrxFee mocks base method.
func (m *MockRepository) ListTrxFee(symbol string, startTime, endTime int64, page, limit int) ([]repository.UniTrxFee, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTrxFee", symbol, startTime, endTime, page, limit)
	ret0, _ := ret[0].([]repository.UniTrxFee)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTrxFee indicates an expected call of ListTrxFee.
func (mr *MockRepositoryMockRecorder) ListTrxFee(symbol, startTime, endTime, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTrxFee", reflect.TypeOf((*MockRepository)(nil).ListTrxFee), symbol, startTime, endTime, page, limit)
}