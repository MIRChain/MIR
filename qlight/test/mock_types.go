// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package qlighttest is a generated GoMock package.
package test

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/pavelkrolevets/MIR-pro/common"
	types "github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	qlightptm "github.com/pavelkrolevets/MIR-pro/private/engine/qlightptm"
	qlight "github.com/pavelkrolevets/MIR-pro/qlight"
)

// MockPrivateStateRootHashValidator is a mock of PrivateStateRootHashValidator interface.
type MockPrivateStateRootHashValidator struct {
	ctrl     *gomock.Controller
	recorder *MockPrivateStateRootHashValidatorMockRecorder
}

// MockPrivateStateRootHashValidatorMockRecorder is the mock recorder for MockPrivateStateRootHashValidator.
type MockPrivateStateRootHashValidatorMockRecorder struct {
	mock *MockPrivateStateRootHashValidator
}

// NewMockPrivateStateRootHashValidator creates a new mock instance.
func NewMockPrivateStateRootHashValidator(ctrl *gomock.Controller) *MockPrivateStateRootHashValidator {
	mock := &MockPrivateStateRootHashValidator{ctrl: ctrl}
	mock.recorder = &MockPrivateStateRootHashValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrivateStateRootHashValidator) EXPECT() *MockPrivateStateRootHashValidatorMockRecorder {
	return m.recorder
}

// ValidatePrivateStateRoot mocks base method.
func (m *MockPrivateStateRootHashValidator) ValidatePrivateStateRoot(blockHash, blockPublicStateRoot common.Hash) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatePrivateStateRoot", blockHash, blockPublicStateRoot)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidatePrivateStateRoot indicates an expected call of ValidatePrivateStateRoot.
func (mr *MockPrivateStateRootHashValidatorMockRecorder) ValidatePrivateStateRoot(blockHash, blockPublicStateRoot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatePrivateStateRoot", reflect.TypeOf((*MockPrivateStateRootHashValidator)(nil).ValidatePrivateStateRoot), blockHash, blockPublicStateRoot)
}

// MockPrivateClientCache is a mock of PrivateClientCache interface.
type MockPrivateClientCache struct {
	ctrl     *gomock.Controller
	recorder *MockPrivateClientCacheMockRecorder
}

// MockPrivateClientCacheMockRecorder is the mock recorder for MockPrivateClientCache.
type MockPrivateClientCacheMockRecorder struct {
	mock *MockPrivateClientCache
}

// NewMockPrivateClientCache creates a new mock instance.
func NewMockPrivateClientCache(ctrl *gomock.Controller) *MockPrivateClientCache {
	mock := &MockPrivateClientCache{ctrl: ctrl}
	mock.recorder = &MockPrivateClientCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrivateClientCache) EXPECT() *MockPrivateClientCacheMockRecorder {
	return m.recorder
}

// AddPrivateBlock mocks base method.
func (m *MockPrivateClientCache) AddPrivateBlock(blockPrivateData qlight.BlockPrivateData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPrivateBlock", blockPrivateData)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPrivateBlock indicates an expected call of AddPrivateBlock.
func (mr *MockPrivateClientCacheMockRecorder) AddPrivateBlock(blockPrivateData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPrivateBlock", reflect.TypeOf((*MockPrivateClientCache)(nil).AddPrivateBlock), blockPrivateData)
}

// CheckAndAddEmptyEntry mocks base method.
func (m *MockPrivateClientCache) CheckAndAddEmptyEntry(hash common.EncryptedPayloadHash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CheckAndAddEmptyEntry", hash)
}

// CheckAndAddEmptyEntry indicates an expected call of CheckAndAddEmptyEntry.
func (mr *MockPrivateClientCacheMockRecorder) CheckAndAddEmptyEntry(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAndAddEmptyEntry", reflect.TypeOf((*MockPrivateClientCache)(nil).CheckAndAddEmptyEntry), hash)
}

// ValidatePrivateStateRoot mocks base method.
func (m *MockPrivateClientCache) ValidatePrivateStateRoot(blockHash, blockPublicStateRoot common.Hash) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidatePrivateStateRoot", blockHash, blockPublicStateRoot)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidatePrivateStateRoot indicates an expected call of ValidatePrivateStateRoot.
func (mr *MockPrivateClientCacheMockRecorder) ValidatePrivateStateRoot(blockHash, blockPublicStateRoot interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidatePrivateStateRoot", reflect.TypeOf((*MockPrivateClientCache)(nil).ValidatePrivateStateRoot), blockHash, blockPublicStateRoot)
}

// MockPrivateBlockDataResolver is a mock of PrivateBlockDataResolver interface.
type MockPrivateBlockDataResolver struct {
	ctrl     *gomock.Controller
	recorder *MockPrivateBlockDataResolverMockRecorder
}

// MockPrivateBlockDataResolverMockRecorder is the mock recorder for MockPrivateBlockDataResolver.
type MockPrivateBlockDataResolverMockRecorder struct {
	mock *MockPrivateBlockDataResolver
}

// NewMockPrivateBlockDataResolver creates a new mock instance.
func NewMockPrivateBlockDataResolver(ctrl *gomock.Controller) *MockPrivateBlockDataResolver {
	mock := &MockPrivateBlockDataResolver{ctrl: ctrl}
	mock.recorder = &MockPrivateBlockDataResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPrivateBlockDataResolver) EXPECT() *MockPrivateBlockDataResolverMockRecorder {
	return m.recorder
}

// PrepareBlockPrivateData mocks base method.
func (m *MockPrivateBlockDataResolver) PrepareBlockPrivateData(block *types.Block[nist.PublicKey], psi string) (*qlight.BlockPrivateData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareBlockPrivateData", block, psi)
	ret0, _ := ret[0].(*qlight.BlockPrivateData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareBlockPrivateData indicates an expected call of PrepareBlockPrivateData.
func (mr *MockPrivateBlockDataResolverMockRecorder) PrepareBlockPrivateData(block, psi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareBlockPrivateData", reflect.TypeOf((*MockPrivateBlockDataResolver)(nil).PrepareBlockPrivateData), block, psi)
}

// MockAuthProvider is a mock of AuthProvider interface.
type MockAuthProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAuthProviderMockRecorder
}

// MockAuthProviderMockRecorder is the mock recorder for MockAuthProvider.
type MockAuthProviderMockRecorder struct {
	mock *MockAuthProvider
}

// NewMockAuthProvider creates a new mock instance.
func NewMockAuthProvider(ctrl *gomock.Controller) *MockAuthProvider {
	mock := &MockAuthProvider{ctrl: ctrl}
	mock.recorder = &MockAuthProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthProvider) EXPECT() *MockAuthProviderMockRecorder {
	return m.recorder
}

// Authorize mocks base method.
func (m *MockAuthProvider) Authorize(token, psi string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", token, psi)
	ret0, _ := ret[0].(error)
	return ret0
}

// Authorize indicates an expected call of Authorize.
func (mr *MockAuthProviderMockRecorder) Authorize(token, psi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockAuthProvider)(nil).Authorize), token, psi)
}

// Initialize mocks base method.
func (m *MockAuthProvider) Initialize() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockAuthProviderMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockAuthProvider)(nil).Initialize))
}

// MockCacheWithEmpty is a mock of CacheWithEmpty interface.
type MockCacheWithEmpty struct {
	ctrl     *gomock.Controller
	recorder *MockCacheWithEmptyMockRecorder
}

// MockCacheWithEmptyMockRecorder is the mock recorder for MockCacheWithEmpty.
type MockCacheWithEmptyMockRecorder struct {
	mock *MockCacheWithEmpty
}

// NewMockCacheWithEmpty creates a new mock instance.
func NewMockCacheWithEmpty(ctrl *gomock.Controller) *MockCacheWithEmpty {
	mock := &MockCacheWithEmpty{ctrl: ctrl}
	mock.recorder = &MockCacheWithEmptyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheWithEmpty) EXPECT() *MockCacheWithEmptyMockRecorder {
	return m.recorder
}

// Cache mocks base method.
func (m *MockCacheWithEmpty) Cache(privateTxData *qlightptm.CachablePrivateTransactionData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cache", privateTxData)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cache indicates an expected call of Cache.
func (mr *MockCacheWithEmptyMockRecorder) Cache(privateTxData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cache", reflect.TypeOf((*MockCacheWithEmpty)(nil).Cache), privateTxData)
}

// CheckAndAddEmptyToCache mocks base method.
func (m *MockCacheWithEmpty) CheckAndAddEmptyToCache(hash common.EncryptedPayloadHash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CheckAndAddEmptyToCache", hash)
}

// CheckAndAddEmptyToCache indicates an expected call of CheckAndAddEmptyToCache.
func (mr *MockCacheWithEmptyMockRecorder) CheckAndAddEmptyToCache(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAndAddEmptyToCache", reflect.TypeOf((*MockCacheWithEmpty)(nil).CheckAndAddEmptyToCache), hash)
}
