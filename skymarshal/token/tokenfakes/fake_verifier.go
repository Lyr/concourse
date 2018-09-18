// Code generated by counterfeiter. DO NOT EDIT.
package tokenfakes

import (
	"context"
	"sync"

	"github.com/concourse/skymarshal/token"
	"golang.org/x/oauth2"
)

type FakeVerifier struct {
	VerifyStub        func(context.Context, *oauth2.Token) (*token.VerifiedClaims, error)
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		arg1 context.Context
		arg2 *oauth2.Token
	}
	verifyReturns struct {
		result1 *token.VerifiedClaims
		result2 error
	}
	verifyReturnsOnCall map[int]struct {
		result1 *token.VerifiedClaims
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVerifier) Verify(arg1 context.Context, arg2 *oauth2.Token) (*token.VerifiedClaims, error) {
	fake.verifyMutex.Lock()
	ret, specificReturn := fake.verifyReturnsOnCall[len(fake.verifyArgsForCall)]
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		arg1 context.Context
		arg2 *oauth2.Token
	}{arg1, arg2})
	fake.recordInvocation("Verify", []interface{}{arg1, arg2})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.verifyReturns.result1, fake.verifyReturns.result2
}

func (fake *FakeVerifier) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *FakeVerifier) VerifyArgsForCall(i int) (context.Context, *oauth2.Token) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].arg1, fake.verifyArgsForCall[i].arg2
}

func (fake *FakeVerifier) VerifyReturns(result1 *token.VerifiedClaims, result2 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 *token.VerifiedClaims
		result2 error
	}{result1, result2}
}

func (fake *FakeVerifier) VerifyReturnsOnCall(i int, result1 *token.VerifiedClaims, result2 error) {
	fake.VerifyStub = nil
	if fake.verifyReturnsOnCall == nil {
		fake.verifyReturnsOnCall = make(map[int]struct {
			result1 *token.VerifiedClaims
			result2 error
		})
	}
	fake.verifyReturnsOnCall[i] = struct {
		result1 *token.VerifiedClaims
		result2 error
	}{result1, result2}
}

func (fake *FakeVerifier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVerifier) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ token.Verifier = new(FakeVerifier)
