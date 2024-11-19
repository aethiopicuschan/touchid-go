//go:build darwin

package touchid

/*
#cgo CFLAGS: -x objective-c -fmodules -fblocks
#cgo LDFLAGS: -framework CoreFoundation -framework LocalAuthentication -framework Foundation
#import <LocalAuthentication/LocalAuthentication.h>

typedef struct {
  LAContext *context;
  dispatch_semaphore_t sema;
  __block int result;
} AuthContext;

AuthContext* InitAuthContext() {
  AuthContext *authCtx = (AuthContext *)malloc(sizeof(AuthContext));
  authCtx->context = [[LAContext alloc] init];
  authCtx->sema = dispatch_semaphore_create(0);
  authCtx->result = 0;
  return authCtx;
}

void ReleaseAuthContext(AuthContext *authCtx) {
  if (authCtx) {
    [authCtx->context release];
    dispatch_release(authCtx->sema);
    free(authCtx);
  }
}

void CancelAuthentication(AuthContext *authCtx) {
  if (authCtx && authCtx->context) {
    [authCtx->context invalidate];
  }
}

int AuthenticateWithContext(AuthContext *authCtx, char const* reason) {
  NSError *authError = nil;
  NSString *nsReason = [NSString stringWithUTF8String:reason];

  // Use LAPolicyDeviceOwnerAuthentication to allow both biometrics and password
  if ([authCtx->context canEvaluatePolicy:LAPolicyDeviceOwnerAuthentication error:&authError]) {
    [authCtx->context evaluatePolicy:LAPolicyDeviceOwnerAuthentication
      localizedReason:nsReason
      reply:^(BOOL success, NSError *error) {
        if (success) {
          authCtx->result = 1;
        } else {
          // Handle failure due to user canceling or other errors
          if (error.code == LAErrorUserFallback) {
            authCtx->result = 2;  // Password used
          } else {
            authCtx->result = 3;  // Other error
          }
        }
        dispatch_semaphore_signal(authCtx->sema);
      }];
  }

  dispatch_semaphore_wait(authCtx->sema, DISPATCH_TIME_FOREVER);
  return authCtx->result;
}
*/
import (
	"C"
)
import (
	"context"
	"errors"
	"unsafe"
)

func Authenticate(ctx context.Context, reason string) (bool, error) {
	reasonStr := C.CString(reason)
	defer C.free(unsafe.Pointer(reasonStr))

	authCtx := C.InitAuthContext()
	defer C.ReleaseAuthContext(authCtx)

	resultChan := make(chan int, 1)
	go func() {
		result := C.AuthenticateWithContext(authCtx, reasonStr)
		resultChan <- int(result)
	}()

	select {
	case <-ctx.Done():
		C.CancelAuthentication(authCtx)
		return false, ctx.Err()
	case result := <-resultChan:
		switch result {
		case 1:
			return true, nil
		case 2:
			return false, nil
		}
	}

	return false, errors.New("unexpected error")
}
